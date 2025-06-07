package models

import (
	"time"
)

type Auction struct {
	OfferID      uint       `json:"id" gorm:"primaryKey"`
	DateEnd      time.Time  `json:"date_end"`
	BuyNowPrice  *uint      `json:"buy_now_price,omitempty"`
	InitialPrice uint       `json:"initial_price"`
	Offer        *SaleOffer `gorm:"foreignKey:OfferID;references:ID"`
}
