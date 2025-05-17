package auctionws

import "encoding/json"

type MsgType string

const (
	MsgBid         MsgType = "bid"
	MsgEndAuction  MsgType = "end_auction"
	MsgError       MsgType = "error"
	MsgSubscribe   MsgType = "subscribe"
	MsgUnsubscribe MsgType = "unsubscribe"
)

type Envelope struct {
	MessageType MsgType         `json:"type"`
	Data        json.RawMessage `json:"data,omitempty"`
}
type BidPayload struct {
	AuctionID string `json:"auction_id"`
	Amount    int64  `json:"amount"`
	UserID    string `json:"user_id"`
}

type EndAuctionPayload struct {
	AuctionID string `json:"auction_id"`
	Winner    string `json:"winner"`
	Amount    int64  `json:"amount"`
}

type SubscribePayload struct {
	Auctions []string `json:"auctions"`
}
type UnsubscribePayload struct {
	Auctions []string `json:"auctions"`
}

func NewBidEnvelope(auctionID string, amount int64, userID string) *Envelope {
	payload := BidPayload{
		AuctionID: auctionID,
		Amount:    amount,
		UserID:    userID,
	}
	data, _ := json.Marshal(payload)
	return &Envelope{
		MessageType: MsgBid,
		Data:        data,
	}
}
