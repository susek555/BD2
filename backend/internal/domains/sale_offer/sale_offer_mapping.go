package sale_offer

import (
	"slices"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
)

func (dto *CreateSaleOfferDTO) MapToSaleOffer() (*SaleOffer, error) {
	if err := dto.validateParams(); err != nil {
		return nil, err
	}
	date, err := parseDate(dto.RegistrationDate)
	if err != nil {
		return nil, err
	}
	return &SaleOffer{
		Description: dto.Description,
		Price:       dto.Price,
		Margin:      dto.Margin,
		Car: &car.Car{
			Vin:                dto.Vin,
			ProductionYear:     dto.ProductionYear,
			Mileage:            dto.Mileage,
			NumberOfDoors:      dto.NumberOfDoors,
			NumberOfSeats:      dto.NumberOfSeats,
			EnginePower:        dto.EnginePower,
			EngineCapacity:     dto.EngineCapacity,
			RegistrationNumber: dto.RegistrationNumber,
			RegistrationDate:   date,
			Color:              dto.Color,
			FuelType:           dto.FuelType,
			Transmission:       dto.Transmission,
			NumberOfGears:      dto.NumberOfGears,
			Drive:              dto.Drive,
			ModelID:            dto.ModelID,
		},
	}, nil
}

func parseDate(date string) (time.Time, error) {
	layout := "02/01/2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func IsParamValid[T comparable](param T, validParams []T) bool {
	return slices.Contains(validParams, param)
}

func (dto *CreateSaleOfferDTO) validateParams() error {
	if !IsParamValid(dto.Color, car_params.Colors) {
		return ErrInvalidColor
	}
	if !IsParamValid(dto.FuelType, car_params.Types) {
		return ErrInvalidFuelType
	}
	if !IsParamValid(dto.Transmission, car_params.Transmissions) {
		return ErrInvalidTransmission
	}
	if !IsParamValid(dto.Drive, car_params.Drives) {
		return ErrInvalidDrive
	}
	return nil
}
