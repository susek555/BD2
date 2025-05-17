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
