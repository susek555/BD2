package models

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
)

type SaleOffer struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	UserID      uint              `json:"user_id"`
	Description string            `json:"description"`
	Price       uint              `json:"price"`
	DateOfIssue time.Time         `json:"date_of_issue"`
	Margin      enums.MarginValue `json:"margin" gorm:"type:MARGIN_VALUE"`
	Status      enums.Status      `json:"status" gorm:"type:OFFER_STATUS"`
	IsAuction   bool              `json:"is_auction"`
	User        *User             `gorm:"foreignKey:UserID;references:ID"`
	Car         *Car              `gorm:"foreignKey:OfferID;references:ID"`
	Auction     *Auction          `gorm:"foreignKey:OfferID;references:ID"`
}

func (o *SaleOffer) GetID() uint {
	return o.ID
}

func (o *SaleOffer) IsAuctionOffer() bool {
	return o.IsAuction
}

func (o *SaleOffer) BelongsToUser(userID uint) bool {
	return o.UserID == userID
}

func (o *SaleOffer) GetStatus() enums.Status {
	return o.Status
}
