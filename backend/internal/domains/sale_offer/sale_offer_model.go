package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

type SaleOffer struct {
	ID          uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint        `json:"user_id"`
	User        *user.User  `gorm:"foreignKey:UserID;references:ID"`
	Description string      `json:"description"`
	Price       uint        `json:"price"`
	Margin      MarginValue `json:"margin"`
	DateOfIssue time.Time   `json:"date_of_issue"`
	Car         *car.Car    `gorm:"foreignKey:OfferID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Auction     *Auction    `gorm:"foreignKey:OfferID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
