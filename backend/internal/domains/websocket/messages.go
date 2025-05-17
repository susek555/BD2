package websocket

type Message struct {
	Type string `json:"type"`
}

type BidMessage struct {
	Message
	AuctionId string `json:"auction_id"`
	Amount    int    `json:"amount"`
	UserId    string `json:"user_id"`
}

type EndAuctionMessage struct {
	Message
	AuctionId string `json:"auction_id"`
	Winner    string `json:"winner"`
	Amount    int    `json:"amount"`
}
