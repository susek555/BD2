package bid

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

type Bid struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Amount    uint      `json:"amount" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;index"`
	AuctionID uint      `json:"auction_id" gorm:"not null;index"`
	// TODO: fix after refactor models dir `gorm:"foreignKey:AuctionID;references:OfferID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BidderID uint       `json:"bidder_id" gorm:"not null;index"`
	Bidder   *user.User `gorm:"foreignKey:BidderID;references:ID"`
}
