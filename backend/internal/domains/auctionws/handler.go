package auctionws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}
		uid, err := auth.GetUserId(c)
		if err != nil {
			log.Println("Error getting user ID:", err)
			conn.Close()
			return
		}
		uidStr := strconv.Itoa(int(uid))
		client := &Client{
			conn:   conn,
			send:   make(chan []byte, 128),
			userID: uidStr,
			hub:    hub,
			rooms:  make(map[string]bool),
		}
		hub.register <- client
		go client.writePump()
		go client.readPump()
	}
}
