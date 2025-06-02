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
}

func (v *SaleOfferView) BelongsToUser(userID uint) bool {
	return v.UserID == userID
}
