package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
)

type Car struct {
	OfferID            uint                    `json:"id" gorm:"primaryKey"`
	Vin                string                  `json:"vin"`
	ProductionYear     uint                    `json:"production_year"`
	Mileage            uint                    `json:"mileage"`
	NumberOfDoors      uint                    `json:"number_of_door" gorm:"check:number_of_doors BETWEEN 1 AND 6"`
	NumberOfSeats      uint                    `json:"number_of_seats" gorm:"check:number_of_seats BETWEEN 2 AND 100"`
	EnginePower        uint                    `json:"engine_power" gorm:"check:engine_power BETWEEN 1 AND 9999"`
	EngineCapacity     uint                    `json:"engine_capacity" gorm:"check:engine_capacity BETWEEN 1 AND 9000"`
	RegistrationNumber string                  `json:"registration_number"`
	RegistrationDate   time.Time               `json:"registration_date"`
	Color              car_params.Color        `json:"color"`
	FuelType           car_params.FuelType     `json:"fuel_type"`
	Transmission       car_params.Transmission `json:"transmission"`
	NumberOfGears      uint                    `json:"number_of_gears" gorm:"check:number_of_gears BETWEEN 1 AND 10"`
	Drive              car_params.Drive        `json:"drive"`
	ModelID            uint                    `json:"model_id"`
	Model              model.Model             `gorm:"foreignKey:ModelID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	SaleOffer          SaleOffer               `gorm:"foreignKey:OfferID;references:ID"`
}
