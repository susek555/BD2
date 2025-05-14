package sale_offer

import (
	"os/user"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
)

type SaleOffer struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint       `json:"user_id"`
	User        *user.User `gorm:"foreignKey:UserID;references:ID"`
	Description string     `json:"description"`
	Price       uint       `json:"price"`
	Margin      uint       `json:"margin"`
	DateOfIssue time.Time  `jsont:"date_of_issue"`
	CarID       uint       `json:"car_id"`
	Car         *car.Car   `gorm:"foreignKey:CarID;references:ID"`
	Auction     *Auction   `gorm:"foreignKey:OfferID;references:ID;constraint:OnDelete:CASCADE"`
}
