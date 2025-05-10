package sale_offer

import (
	"time"
)

type Auction struct {
	OfferID     uint       `json:"id" gorm:"primaryKey"`
	DateEnd     time.Time  `json:"date_end" gorm:"not null;index"`
	BuyNowPrice uint       `json:"buy_now_price" gorm:"index"`
	Offer       *SaleOffer `gorm:"foreignKey:OfferID;references:ID;constraint:OnDelete:CASCADE"`
}
