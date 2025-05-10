package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
)

type CreateSaleOfferDTO struct {
	Description        string                  `json:"description"`
	Price              uint                    `json:"price"`
	Margin             uint                    `json:"margin"`
	Vin                string                  `json:"vin"`
	ProductionYear     uint                    `json:"production_year"`
	Mileage            uint                    `json:"mileage"`
	NumberOfDoors      uint                    `json:"number_of_doors"`
	NumberOfSeats      uint                    `json:"number_of_seats"`
	EnginePower        uint                    `json:"engine_power"`
	EngineCapacity     uint                    `json:"engine_capacity"`
	RegistrationNumber string                  `json:"registration_number"`
	RegistrationDate   string                  `json:"registration_date"`
	Color              car_params.Color        `json:"color"`
	FuelType           car_params.FuelType     `json:"fuel_type"`
	Transmission       car_params.Transmission `json:"transmission"`
	NumberOfGears      uint                    `json:"number_of_gears"`
	Drive              car_params.Drive        `json:"drive"`
	ModelID            uint                    `json:"model_id"`
}

type RetrieveSaleOfferDTO struct {
	Description        string                  `json:"description"`
	Price              uint                    `json:"price"`
	Margin             uint                    `json:"margin"`
	Vin                string                  `json:"vin"`
	ProductionYear     uint                    `json:"production_year"`
	Mileage            uint                    `json:"mileage"`
	NumberOfDoors      uint                    `json:"number_of_doors"`
	NumberOfSeats      uint                    `json:"number_of_seats"`
	EnginePower        uint                    `json:"engine_power"`
	EngineCapacity     uint                    `json:"engine_capacity"`
	RegistrationNumber string                  `json:"registration_number"`
	RegistrationDate   string                  `json:"registration_date"`
	Color              car_params.Color        `json:"color"`
	FuelType           car_params.FuelType     `json:"fuel_type"`
	Transmission       car_params.Transmission `json:"transmission"`
	NumberOfGears      uint                    `json:"number_of_gears"`
	Drive              car_params.Drive        `json:"drive"`
	Brand              string                  `json:"brand"`
	Model              string                  `json:"model"`
}
