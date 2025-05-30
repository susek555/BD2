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
)

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
}
type subscription struct {
	auctionID string
	client    *Client
}
type outbound struct {
	auctionID string
	data      []byte
	excludeID string
}

func NewHub(clientNotificationRepo notification.ClientNotificationRepositoryInterface) *Hub {
	return &Hub{
		rooms:                  make(map[string]map[*Client]struct{}),
		clients:                make(map[string]*Client),
		register:               make(chan *Client),
		unregister:             make(chan *Client),
		subscribe:              make(chan subscription),
		unsubscribe:            make(chan subscription),
		broadcast:              make(chan outbound, 1024),
		clientNotificationRepo: clientNotificationRepo,
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
			h.addToRoom(sub.auctionID, sub.client)
		case sub := <-h.unsubscribe:
			h.removeFromRoom(sub.auctionID, sub.client)
		case msg := <-h.broadcast:
			h.fanOut(msg)
		}
	}
}

func (h *Hub) addToRoom(auctionID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[auctionID] == nil {
		h.rooms[auctionID] = make(map[*Client]struct{})
	}
	h.rooms[auctionID][client] = struct{}{}
	client.rooms[auctionID] = true
}

func (h *Hub) removeFromRoom(auctionID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[auctionID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.rooms, auctionID)
		}
	}
	delete(client.rooms, auctionID)
}

func (h *Hub) removeClient(client *Client) {
	for auctionID := range client.rooms {
		h.removeFromRoom(auctionID, client)
	}
	close(client.send)
}

func (h *Hub) fanOut(msg outbound) {
	h.mu.RLock()
	room, ok := h.rooms[msg.auctionID]
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
			pubsub := rdb.PSubscribe(ctx, "auction.*")
			ch := pubsub.Channel()

			for msg := range ch {
				id := strings.TrimPrefix(msg.Channel, "auction.")
				h.broadcast <- outbound{
					auctionID: id,
					data:      []byte(msg.Payload),
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

func (h *Hub) SubscribeUser(uid, auctionID string) {
	h.mu.RLock()
	cl, ok := h.clients[uid]
	h.mu.RUnlock()
	if !ok {
		return
	}
	h.subscribe <- subscription{auctionID, cl}
}

func (h *Hub) BroadcastLocal(auctionID string, data []byte, excludeID string) {
	h.broadcast <- outbound{
		auctionID: auctionID,
		data:      data,
		excludeID: excludeID,
	}
}

func (h *Hub) SaveNotificationForClients(auctionID string, userID uint, n *models.Notification) error {
	h.mu.RLock()
	clients, ok := h.rooms[auctionID]
	h.mu.RUnlock()
	if !ok {
		return nil
	}
	for client := range clients {
		clientID, err := strconv.ParseUint(client.userID, 10, 64)
		if err != nil {
			log.Printf("Failed to convert userID %s to uint: %v", client.userID, err)
			continue
		}
		if clientID == uint64(userID) {
			log.Printf("Skipping saving notification for client %s as it is the same as the userID %d", client.userID, userID)
			continue
		}
		clientNotification := notification.MapToClientNotification(n, uint(clientID))
		if err := h.clientNotificationRepo.Create(clientNotification); err != nil {
			log.Printf("Failed to save notification for client %s: %v", client.userID, err)
			continue
		}
	}
	return nil
}

func (h *Hub) SendFourLatestNotificationsToClient(auctionID, userID string) {
	h.mu.RLock()
	room, ok := h.rooms[auctionID]
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

		notifications, err := h.clientNotificationRepo.GetLatestByUserId(uint(uid), 4)
		if err != nil {
			log.Printf("hub: cannot fetch latest notification for userID %q: %v", client.userID, err)
			continue
		}
		if len(notifications) == 0 {
			log.Printf("hub: no notifications for userID %q", client.userID)
			continue
		}
		bare := notification.MapToNotifications(notifications)
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
