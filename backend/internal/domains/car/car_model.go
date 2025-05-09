package car

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/model"
)

type Car struct {
	ID                 uint                    `json:"id" gorm:"primaryKey;autoIncrement"`
	Vin                string                  `json:"vin"`
	ProductionYear     uint                    `json:"production_year"`
	Mileage            uint                    `json:"mileage"`
	NumberOfDoors      uint                    `json:"number_of_door" gorm:"check:number_of_doors BETWEEN 1 AND 6"`
	NumberOfSeats      uint                    `json:"number_of_seats" gorm:"check:number_of_seats BETWEEN 2 AND 100"`
	EnginePower        uint                    `json:"engine_power" gorm:"check:engine_power <= 9999"`
	EngineCapacity     uint                    `json:"engine_capacity" gorm:"check engine_capacity <= 9000"`
	RegistrationNumber string                  `json:"registration_number"`
	RegistrationDate   time.Time               `json:"registration_date"`
	Color              car_params.Color        `json:"color"`
	FuelType           car_params.FuelType     `json:"fuel_type"`
	Transmission       car_params.Transmission `json:"transmission"`
	NumberOfGears      uint                    `json:"number_of_gears"`
	Drive              car_params.Drive        `json:"drive"`
	ModelID            uint                    `json:"model_id"`
	Model              *model.Model            `gorm:"foreignKey:ModelID;references:ID;OnDelete:SET NULL"`
}
