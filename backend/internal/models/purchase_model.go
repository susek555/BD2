package models

import "time"

type Purchase struct {
	ID         uint       `json:"id"`
	OfferID    uint       `json:"offer_id"`
	BuyerID    uint       `json:"buyer_id"`
	FinalPrice uint       `json:"final_price"`
	IssueDate  time.Time  `json:"issue_date"`
	Offer      *SaleOffer `gorm:"foreignKey:OfferID;references:ID"`
	Buyer      *User      `gorm:"foreignKey:BuyerID;references:ID"`
}
