package scheduler

type BuyNowRecord struct {
	AuctionID string `json:"auction_id"`
	Price     uint   `json:"price"`
	BuyerID   uint   `json:"buyer_id"`
}
