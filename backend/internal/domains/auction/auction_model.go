package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"time"
)

type Auction struct {
	OfferID     uint                  `json:"id" gorm:"primaryKey"`
	DateEnd     time.Time             `json:"date_end" gorm:"not null;index"`
	BuyNowPrice uint                  `json:"buy_now_price" gorm:"index"`
	Offer       *sale_offer.SaleOffer `gorm:"foreignKey:OfferID;references:ID;constraint:OnDelete:CASCADE"`
}
