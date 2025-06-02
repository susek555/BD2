package sale_offer

import "time"

type CarSaleOfferView struct {
	ID                 uint
	UserID             uint
	Description        string
	Price              int
	DateOfIssue        time.Time
	Margin             int
	Status             string
	Vin                string
	ProductionYear     int
	Mileage            int
	NumberOfDoors      int
	NumberOfSeats      int
	EnginePower        int
	EngineCapacity     int
	RegistrationNumber string
	RegistrationDate   time.Time
	Color              string
	FuelType           string
	Drive              string
	Transmission       string
	NumberOfGears      int
	Brand              string
	Model              string
}
