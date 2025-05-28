package models

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
)

type Car struct {
	OfferID            uint                    `json:"id" gorm:"primaryKey"`
	Vin                string                  `json:"vin"`
	ProductionYear     uint                    `json:"production_year"`
	Mileage            uint                    `json:"mileage"`
	NumberOfDoors      uint                    `json:"number_of_doors"`
	NumberOfSeats      uint                    `json:"number_of_seats"`
	EnginePower        uint                    `json:"engine_power"`
	EngineCapacity     uint                    `json:"engine_capacity"`
	RegistrationNumber string                  `json:"registration_number"`
	RegistrationDate   time.Time               `json:"registration_date"`
	Color              car_params.Color        `json:"color"`
	FuelType           car_params.FuelType     `json:"fuel_type"`
	Transmission       car_params.Transmission `json:"transmission"`
	NumberOfGears      uint                    `json:"number_of_gears"`
	Drive              car_params.Drive        `json:"drive"`
	ModelID            uint                    `json:"model_id"`
	Model              *Model                  `gorm:"foreignKey:ModelID;references:ID"`
}
