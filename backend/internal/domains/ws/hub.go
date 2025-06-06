package ws

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=HubInterface --output=../../test/mocks --case=snake --with-expecter
type HubInterface interface {
	Run()
	StartRedisFanIn(ctx context.Context, rdb *redis.Client)
	SubscribeUser(uid, offerID string)
	BroadcastLocal(offerID string, data []byte, excludeID string)
	SaveNotificationForClients(offerID string, userID uint, n *models.Notification) error
	SendFourLatestNotificationsToClient(offerID, userID string)
	LoadClientToRooms(userID string)
	UnsubscribeUser(userID, offerID string)
	RemoveRoom(offerID string)
}

type Hub struct {
	rooms                  map[string]map[*Client]struct{}
	clients                map[string]*Client
	register               chan *Client
	unregister             chan *Client
	subscribe              chan subscription
	unsubscribe            chan subscription
	broadcast              chan outbound
	mu                     sync.RWMutex
	clientNotificationRepo notification.ClientNotificationRepositoryInterface
	db                     *gorm.DB
}
type subscription struct {
	offerID string
	client  *Client
}
type outbound struct {
	offerID   string
	data      []byte
	excludeID string
}

func NewHub(clientNotificationRepo notification.ClientNotificationRepositoryInterface, db *gorm.DB) HubInterface {
	return &Hub{
		rooms:                  make(map[string]map[*Client]struct{}),
		clients:                make(map[string]*Client),
		register:               make(chan *Client),
		unregister:             make(chan *Client),
		subscribe:              make(chan subscription),
		unsubscribe:            make(chan subscription),
		broadcast:              make(chan outbound, 1024),
		clientNotificationRepo: clientNotificationRepo,
		db:                     db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Println("Client registered:", client.userID)
		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, client.userID)
			h.mu.Unlock()
			h.removeClient(client)
		case sub := <-h.subscribe:
			h.addToRoom(sub.offerID, sub.client)
		case sub := <-h.unsubscribe:
			h.removeFromRoom(sub.offerID, sub.client)
		case msg := <-h.broadcast:
			h.fanOut(msg)
		}
	}
}

func (h *Hub) addToRoom(offerID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[offerID] == nil {
		h.rooms[offerID] = make(map[*Client]struct{})
	}
	h.rooms[offerID][client] = struct{}{}
	client.rooms[offerID] = true
}

func (h *Hub) removeFromRoom(offerID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[offerID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.rooms, offerID)
		}
	}
	delete(client.rooms, offerID)
}

func (h *Hub) removeClient(client *Client) {
	for offerID := range client.rooms {
		h.removeFromRoom(offerID, client)
	}
	close(client.send)
}

func (h *Hub) fanOut(msg outbound) {
	h.mu.RLock()
	room, ok := h.rooms[msg.offerID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	for client := range room {
		if msg.excludeID != "" && client.userID == msg.excludeID {
			continue
		}
		select {
		case client.send <- msg.data:
		default:
			go client.conn.Close()
		}
	}
}

func (h *Hub) StartRedisFanIn(ctx context.Context, rdb *redis.Client) {
	go func() {
		for {
			pubsub := rdb.PSubscribe(ctx, "offer.*")
			ch := pubsub.Channel()

			for msg := range ch {
				id := strings.TrimPrefix(msg.Channel, "offer.")
				h.broadcast <- outbound{
					offerID: id,
					data:    []byte(msg.Payload),
				}
			}

			_ = pubsub.Close()
			for backoff := time.Second; ; {
				select {
				case <-ctx.Done():
					return
				case <-time.After(backoff):
				}

				if backoff < 30*time.Second {
					backoff *= 2
				}
			}
		}
	}()
}

func (h *Hub) SubscribeUser(uid, offerID string) {
	h.mu.RLock()
	cl, ok := h.clients[uid]
	h.mu.RUnlock()
	if !ok {
		return
	}
	h.subscribe <- subscription{offerID, cl}
}

func (h *Hub) BroadcastLocal(offerID string, data []byte, excludeID string) {
	h.broadcast <- outbound{
		offerID:   offerID,
		data:      data,
		excludeID: excludeID,
	}
}

func (h *Hub) SaveNotificationForClients(offerID string, userID uint, n *models.Notification) error {
	var interactions []UserOfferRecord
	err := h.db.
		Table("user_offer_interactions").
		Select("user_id, offer_id").
		Where("offer_id = ?", offerID).
		Find(&interactions).Error
	if err != nil {
		log.Printf("Failed to fetch user offer interactions for offerID %s: %v", offerID, err)
		return err
	}
	unique := make(map[uint]struct{})
	for _, interaction := range interactions {
		unique[interaction.UserID] = struct{}{}
	}
	h.mu.RLock()
	if room, ok := h.rooms[offerID]; ok {
		for client := range room {
			uid, err := strconv.ParseUint(client.userID, 10, 64)
			if err != nil {
				log.Printf("Failed to convert client.userID %s to uint: %v", client.userID, err)
				continue
			}
			unique[uint(uid)] = struct{}{}
		}
	}
	h.mu.RUnlock()
	for uid := range unique {
		if uid == userID {
			continue
		}

		cn := notification.MapToClientNotification(n, uid)
		if err := h.clientNotificationRepo.Create(cn); err != nil {
			log.Printf("hub: cannot save notif for user %d: %v", uid, err)
		}
		log.Printf("Saved notification for user %d", uid)
	}
	return nil
}

func (h *Hub) SendFourLatestNotificationsToClient(offerID, userID string) {
	h.mu.RLock()
	room, ok := h.rooms[offerID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	for client := range room {
		if client.userID == userID {
			continue
		}

		uid, err := strconv.ParseUint(client.userID, 10, 64)
		if err != nil {
			log.Printf("hub: invalid userID %q: %v", client.userID, err)
			continue
		}

		notifications, err := h.clientNotificationRepo.GetLatestByUserID(uint(uid), 4)
		if err != nil {
			log.Printf("hub: cannot fetch latest notification for userID %q: %v", client.userID, err)
			continue
		}
		if len(notifications) == 0 {
			log.Printf("hub: no notifications for userID %q", client.userID)
			continue
		}
		unseenCount, err := h.clientNotificationRepo.GetUnseenCountByUserId(uint(uid))
		if err != nil {
			log.Printf("hub: cannot fetch unseen count for userID %q: %v", client.userID, err)
			continue
		}
		allCount, err := h.clientNotificationRepo.GetAllCountByUserId(uint(uid))
		if err != nil {
			log.Printf("hub: cannot fetch all count for userID %q: %v", client.userID, err)
			continue
		}
		bare := notification.MapToNotificationsDTO(notifications, unseenCount, allCount)
		payload, err := json.Marshal(bare)
		if err != nil {
			log.Printf("hub: cannot marshal notifications for userID %q: %v", client.userID, err)
			continue
		}
		select {
		case client.send <- payload:
		default:
			log.Printf("hub: dropping notifications for userID %q - channel full", client.userID)
			go client.conn.Close()
		}
	}
}

func (h *Hub) UnsubscribeUser(userID, offerID string) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	h.unsubscribe <- subscription{offerID, client}
	log.Printf("hub: removed client %s from room %s", userID, offerID)
}

func (h *Hub) RemoveRoom(offerID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[offerID]; ok {
		delete(h.rooms, offerID)
		log.Printf("hub: removed room %s", offerID)
	} else {
		log.Printf("hub: room %s does not exist", offerID)
	}
}
