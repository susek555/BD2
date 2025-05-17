package auctionws

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	rooms       map[string]map[*Client]struct{}
	clients     map[string]*Client
	register    chan *Client
	unregister  chan *Client
	subscribe   chan subscription
	unsubscribe chan subscription
	broadcast   chan outbound
	mu          sync.RWMutex
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

func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]map[*Client]struct{}),
		clients:     make(map[string]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		subscribe:   make(chan subscription),
		unsubscribe: make(chan subscription),
		broadcast:   make(chan outbound, 1024),
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
