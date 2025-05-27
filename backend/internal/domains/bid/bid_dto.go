package bid

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
)

type CreateBidDTO struct {
	AuctionID uint `json:"auction_id" binding:"required"`
	Amount    uint `json:"amount" binding:"required"`
}

type ProcessingBidDTO struct {
	AuctionID uint            `json:"auction_id" binding:"required"`
	BidderID  uint            `json:"bidder_id" binding:"required"`
	Amount    uint            `json:"amount" binding:"required"`
	Auction   *models.Auction `json:"auction,omitempty"`
}

type RetrieveBidDTO struct {
	AuctionID uint `json:"auction_id" binding:"required"`
	BidderID  uint `json:"bidder_id" binding:"required"`
	Amount    uint `json:"amount" binding:"required"`
}
