package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
)

func (dto *CreateSaleOfferDTO) MapToSaleOffer() (*SaleOffer, error) {
	if err := dto.validateParams(); err != nil {
		return nil, err
	}
	date, err := ParseDate(dto.RegistrationDate)
	if err != nil {
		return nil, err
	}
	return &SaleOffer{
		UserID:      dto.UserID,
		Description: dto.Description,
		Price:       dto.Price,
		Margin:      dto.Margin,
		DateOfIssue: dto.DateOfIssue,
		Car: &Car{
			Vin:                dto.Vin,
			ProductionYear:     dto.ProductionYear,
			Mileage:            dto.Mileage,
			NumberOfDoors:      dto.NumberOfDoors,
			NumberOfSeats:      dto.NumberOfSeats,
			EnginePower:        dto.EnginePower,
			EngineCapacity:     dto.EngineCapacity,
			RegistrationNumber: dto.RegistrationNumber,
			RegistrationDate:   *date,
			Color:              dto.Color,
			FuelType:           dto.FuelType,
			Transmission:       dto.Transmission,
			NumberOfGears:      dto.NumberOfGears,
			Drive:              dto.Drive,
			ModelID:            dto.ModelID,
		},
	}, nil
}

func (offer *SaleOffer) MapToDTO() *RetrieveSaleOfferDTO {
	buyNow, endDate := offer.prepareAuctionValues()
	return &RetrieveSaleOfferDTO{
		Username:           offer.User.Username,
		Description:        offer.Description,
		Price:              offer.Price,
		Margin:             offer.Margin,
		DateOfIssue:        offer.DateOfIssue,
		Vin:                offer.Car.Vin,
		ProductionYear:     offer.Car.ProductionYear,
		Mileage:            offer.Car.Mileage,
		NumberOfDoors:      offer.Car.NumberOfDoors,
		NumberOfSeats:      offer.Car.NumberOfSeats,
		EnginePower:        offer.Car.EnginePower,
		EngineCapacity:     offer.Car.EngineCapacity,
		RegistrationNumber: offer.Car.RegistrationNumber,
		RegistrationDate:   offer.Car.RegistrationDate.Format(LAYOUT),
		Color:              offer.Car.Color,
		FuelType:           offer.Car.FuelType,
		Transmission:       offer.Car.Transmission,
		NumberOfGears:      offer.Car.NumberOfGears,
		Drive:              offer.Car.Drive,
		Brand:              offer.Car.Model.Manufacturer.Name,
		Model:              offer.Car.Model.Name,
		DateEnd:            endDate,
		BuyNowPrice:        buyNow,
	}
}

func (offer *SaleOffer) prepareAuctionValues() (*uint, *time.Time) {
	if offer.Auction == nil {
		return nil, nil
	}
	return &offer.Auction.BuyNowPrice, &offer.Auction.DateEnd
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
