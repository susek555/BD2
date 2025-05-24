package models

import (
	"time"
)

type Auction struct {
	OfferID     uint       `json:"id" gorm:"primaryKey"`
	DateEnd     time.Time  `json:"date_end" gorm:"type:timestamptz;not null;index"`
	BuyNowPrice uint       `json:"buy_now_price" gorm:"index"`
	Offer       *SaleOffer `gorm:"foreignKey:OfferID;references:ID"`
}
