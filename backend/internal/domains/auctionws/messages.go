package auctionws

import "encoding/json"

type MsgType string

const (
	MsgBid        MsgType = "bid"
	MsgEndAuction MsgType = "end_auction"
	MsgError      MsgType = "error"
)

type Message struct {
	Type MsgType `json:"type"`
}

type BidMessage struct {
	Message
	AuctionID string `json:"auction_id"`
	Amount    int64    `json:"amount"`
	UserID    string `json:"user_id"`
}

type EndAuctionMessage struct {
	Message
	AuctionID string `json:"auction_id"`
	Winner    string `json:"winner"`
	Amount    int64    `json:"amount"`
}

type Envelope struct {
	MessageType MsgType `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}
