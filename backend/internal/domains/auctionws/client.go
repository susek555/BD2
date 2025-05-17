package auctionws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = 30 * time.Second
	writeWait  = 10 * time.Second
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userID int64
	hub    *Hub
	rooms  map[string]bool
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(64 << 10)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var env Envelope
		if err = json.Unmarshal(raw, &env); err != nil {
			continue
		}
		switch env.MessageType {
		case MsgSubscribe:
			var p SubscribePayload
			json.Unmarshal(env.Data, &p)
			for _, id := range p.Auctions {
				c.hub.subscribe <- subscription{id, c}
			}
		case MsgUnsubscribe:
			var p SubscribePayload
			json.Unmarshal(env.Data, &p)
			for _, id := range p.Auctions {
				c.hub.unsubscribe <- subscription{id, c}
			}
		}
	}
}
