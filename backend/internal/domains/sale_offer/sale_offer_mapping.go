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
		DateOfIssue: time.Now(),
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
	return &RetrieveSaleOfferDTO{
		ID:             offer.ID,
		Username:       offer.User.Username,
		Name:           offer.Car.Model.Manufacturer.Name + " " + offer.Car.Model.Name,
		Price:          offer.Price,
		Mileage:        offer.Car.Mileage,
		ProductionYear: offer.Car.ProductionYear,
		Color:          offer.Car.Color,
		IsAuction:      offer.Auction != nil,
	}
}

func (offer *SaleOffer) MapToDetailedDTO() *RetrieveDetailedSaleOfferDTO {
	buyNow, endDate := offer.prepareAuctionValues()
	return &RetrieveDetailedSaleOfferDTO{
		ID:                 offer.ID,
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
	type EnumValidation struct {
		value          interface{}
		possibleValues interface{}
		err            error
	}
	validations := []EnumValidation{
		{value: dto.Color, possibleValues: car_params.Colors, err: ErrInvalidColor},
		{value: dto.FuelType, possibleValues: car_params.Types, err: ErrInvalidFuelType},
		{value: dto.Transmission, possibleValues: car_params.Transmissions, err: ErrInvalidTransmission},
		{value: dto.Drive, possibleValues: car_params.Drives, err: ErrInvalidDrive},
		{value: dto.Margin, possibleValues: Margins, err: ErrInvalidMargin},
	}
	for _, validation := range validations {
		if !IsParamValid(&validation.value, validation.possibleValues.([]*any)) {
			return validation.err
		}
	}
	return nil
}
