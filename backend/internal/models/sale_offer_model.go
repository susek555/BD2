package models

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
)

type SaleOffer struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	UserID      uint              `json:"user_id"`
	User        *User             `gorm:"foreignKey:UserID;references:ID"`
	Description string            `json:"description"`
	Price       uint              `json:"price"`
	Margin      enums.MarginValue `json:"margin"`
	DateOfIssue time.Time         `json:"date_of_issue"`
	Status      enums.Status      `json:"status"`
	Car         *Car              `gorm:"foreignKey:OfferID;references:ID"`
	Auction     *Auction          `gorm:"foreignKey:OfferID;references:ID"`
}

func (o *SaleOffer) BelongsToUser(userID uint) bool {
	return o.UserID == userID
}
