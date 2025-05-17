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