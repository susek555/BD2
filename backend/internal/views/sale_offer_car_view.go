package views

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
)

type SaleOfferView struct {
	ID                 uint
	UserID             uint
	Username           string
	Description        string
	Price              uint
	DateOfIssue        time.Time
	Margin             enums.MarginValue
	Status             enums.Status
	IsAuction          bool
	Vin                string
	ProductionYear     uint
	Mileage            uint
	NumberOfDoors      uint
	NumberOfSeats      uint
	EnginePower        uint
	EngineCapacity     uint
	RegistrationNumber string
	RegistrationDate   time.Time
	Color              enums.Color
	FuelType           enums.FuelType
	Transmission       enums.Transmission
	NumberOfGears      uint
	Drive              enums.Drive
	Brand              string
	Model              string
	DateEnd            *time.Time
	BuyNowPrice        *uint
}

func (v *SaleOfferView) GetID() uint {
	return v.ID
}

func (v *SaleOfferView) IsAuctionOffer() bool {
	return v.IsAuction
}

func (v *SaleOfferView) BelongsToUser(userID uint) bool {
	return v.UserID == userID
}

func (v *SaleOfferView) GetStatus() enums.Status {
	return v.Status
}
