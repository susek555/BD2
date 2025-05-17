package auctionws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}
		uid := r.Context().Value("userID").(string)
		client := &Client{
			conn:   conn,
			send:   make(chan []byte, 128),
			userID: uid,
			hub:    hub,
			rooms:  make(map[string]bool),
		}
		hub.register <- client
		go client.writePump()
		go client.readPump()
	}
}
