package ws

import (
	"encoding/json"

	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

type MsgType string

const (
	MsgNotification MsgType = "notification"
	MsgError        MsgType = "error"
	MsgSubscribe    MsgType = "subscribe"
	MsgUnsubscribe  MsgType = "unsubscribe"
)

type Envelope struct {
	MessageType MsgType         `json:"type"`
	Data        json.RawMessage `json:"data,omitempty"`
}

type SubscribePayload struct {
	Offers []string `json:"offers"`
}
type UnsubscribePayload struct {
	Offers []string `json:"offers"`
}

func NewNotificationEnvelope(notification *models.Notification) *Envelope {
	data, err := json.Marshal(notification)
	if err != nil {
		return nil
	}
	return &Envelope{
		MessageType: MsgNotification,
		Data:        data,
	}
}
