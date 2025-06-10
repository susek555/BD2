package models

import (
	"time"
)

type Bid struct {
	ID        uint      `json:"id"`
	Amount    uint      `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	AuctionID uint      `json:"auction_id"`
	BidderID  uint      `json:"bidder_id"`
	Bidder    *User     `gorm:"foreignKey:BidderID;references:ID"`
	Auction   *Auction  `gorm:"foreignKey:AuctionID;references:OfferID;"`
}
