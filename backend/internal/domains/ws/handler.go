package ws

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
		uid, err := auth.GetUserID(c)
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
		hub.LoadClientToRooms(uidStr)
		log.Printf("Client %s connected\n", uidStr)
		go client.writePump()
		go client.readPump()
	}
}

func (h *Hub) LoadClientToRooms(userID string) {
	var records []UserOfferRecord
	err := h.db.
		Table("user_offer_interactions").
		Where("user_id = ?", userID).
		Find(&records).Error
	if err != nil {
		log.Println("Error loading client rooms:", err)
		return
	}
	for _, record := range records {
		h.SubscribeUser(userID, strconv.FormatUint(uint64(record.OfferID), 10))
		log.Printf("User %s subscribed to offer %d\n", userID, record.OfferID)
	}
}
