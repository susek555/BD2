package models

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
)

type Car struct {
	OfferID            uint               `json:"id" gorm:"primaryKey"`
	Vin                string             `json:"vin"`
	ProductionYear     uint               `json:"production_year"`
	Mileage            uint               `json:"mileage"`
	NumberOfDoors      uint               `json:"number_of_doors"`
	NumberOfSeats      uint               `json:"number_of_seats"`
	EnginePower        uint               `json:"engine_power"`
	EngineCapacity     uint               `json:"engine_capacity"`
	RegistrationNumber string             `json:"registration_number"`
	RegistrationDate   time.Time          `json:"registration_date"`
	Color              enums.Color        `json:"color" gorm:"type:COLOR"`
	FuelType           enums.FuelType     `json:"fuel_type" gorm:"type:FUEL_TYPE"`
	Drive              enums.Drive        `json:"drive" gorm:"type:DRIVE"`
	Transmission       enums.Transmission `json:"transmission" gorm:"type:TRANSMISSION"`
	NumberOfGears      uint               `json:"number_of_gears" `
	ModelID            uint               `json:"model_id"`
	Model              *Model             `gorm:"foreignKey:ModelID;references:ID"`
}
