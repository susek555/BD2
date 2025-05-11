package bid

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"time"
)

type Bid struct {
	ID        uint                `json:"id" gorm:"primary_key"`
	Amount    uint                `json:"amount" gorm:"not null;index"`
	CreatedAt time.Time           `json:"created_at" gorm:"not null;index"`
	AuctionID uint                `json:"auction_id" gorm:"not null;index"`
	Auction   *sale_offer.Auction `gorm:"foreignKey:AuctionID;references:OfferID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BidderID  uint                `json:"bidder_id" gorm:"not null;index"`
	Bidder    *user.User          `gorm:"foreignKey:BidderID;references:ID"`
}
