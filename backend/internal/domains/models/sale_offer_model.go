package models

import (
	"time"
)

type SaleOffer struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"user_id"`
	User        *User       `gorm:"foreignKey:UserID;references:ID"`
	Description string      `json:"description"`
	Price       uint        `json:"price"`
	Margin      MarginValue `json:"margin"`
	DateOfIssue time.Time   `json:"date_of_issue"`
	Status      Status      `json:"status"`
	Car         *Car        `gorm:"foreignKey:OfferID;references:ID"`
	Auction     *Auction    `gorm:"foreignKey:OfferID;references:ID"`
}
