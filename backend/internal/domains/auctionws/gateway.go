package auctionws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userID int64
	hub    *Hub
	rooms  map[string]bool
}

type Hub struct {
	rooms       map[string]map[*Client]struct{}
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
}

func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]map[*Client]struct{}),
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
		case _ := <-h.register:
		case client := <-h.unregister:
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
