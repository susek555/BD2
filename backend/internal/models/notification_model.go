package models

import "time"

type Notification struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	OfferID     uint       `json:"offer_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	Offer       *SaleOffer `json:"sale_offer,omitempty" gorm:"foreignKey:OfferID;references:ID"`
}
