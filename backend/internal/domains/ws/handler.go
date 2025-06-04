package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken, ok := c.Get("wsToken")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "missing protocol token"})
			return
		}
		token := rawToken.(string)

		up := websocket.Upgrader{
			CheckOrigin:  func(r *http.Request) bool { return true },
			Subprotocols: []string{token}, // echo
		}

		conn, err := up.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("ws upgrade:", err)
			return
		}

		uidAny, _ := c.Get("userID")
		uid := uidAny.(uint)
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
