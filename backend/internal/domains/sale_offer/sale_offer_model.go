package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
)

type SaleOffer struct {
	ID          uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Description string   `json:"description"`
	Price       uint     `json:"price"`
	Margin      uint     `json:"margin"`
	CarID       uint     `json:"car_id"`
	Car         *car.Car `gorm:"foreignKey:CarID;references:ID"`
	Auction     *Auction
}
