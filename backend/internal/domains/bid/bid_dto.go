package bid

type CreateBidDTO struct {
	AuctionID uint `json:"auction_id" binding:"required"`
	Amount    uint `json:"amount" binding:"required"`
}

type RetrieveBidDTO struct {
	AuctionID uint `json:"auction_id" binding:"required"`
	BidderID  uint `json:"bidder_id" binding:"required"`
	Amount    uint `json:"amount" binding:"required"`
}
