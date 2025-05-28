package sale_offer

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"

	"github.com/go-playground/validator/v10"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
)

func (dto *CreateSaleOfferDTO) MapToSaleOffer() (*models.SaleOffer, error) {
	v := validator.New()
	err := v.Struct(dto)
	if err != nil {
		return nil, ErrMissingFields
	}
	if err := dto.validateParams(); err != nil {
		return nil, err
	}
	date, err := ParseDate(dto.RegistrationDate)
	if err != nil {
		return nil, err
	}
	offer := &models.SaleOffer{Car: &models.Car{}}
	if err := copier.Copy(offer, dto); err != nil {
		return nil, err
	}
	if err := copier.Copy(offer.Car, dto); err != nil {
		return nil, err
	}
	offer.DateOfIssue = time.Now().UTC()
	offer.Car.RegistrationDate = *date
	offer.Status = models.PENDING
	return offer, nil
}

func (dto *UpdateSaleOfferDTO) UpdatedOfferFromDTO(offer *models.SaleOffer) (*models.SaleOffer, error) {
	if err := dto.validateParams(); err != nil {
		return nil, err
	}
	if dto.RegistrationDate != nil {
		date, err := ParseDate(*dto.RegistrationDate)
		if err != nil {
			return nil, err
		}
		offer.Car.RegistrationDate = *date
	}
	if err := copier.Copy(offer, dto); err != nil {
		return nil, err
	}
	if err := copier.Copy(offer.Car, dto); err != nil {
		return nil, err
	}
	return offer, nil
}

func MapToDTO(offer *models.SaleOffer) *RetrieveSaleOfferDTO {
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

func MapToDetailedDTO(offer *models.SaleOffer) *RetrieveDetailedSaleOfferDTO {
	buyNow, endDate := prepareAuctionValues(offer)
	return &RetrieveDetailedSaleOfferDTO{
		ID:                 offer.ID,
		UserID:             offer.UserID,
		Username:           offer.User.Username,
		Description:        offer.Description,
		Price:              offer.Price,
		Margin:             offer.Margin,
		DateOfIssue:        offer.DateOfIssue.Format(formats.DateLayout),
		Vin:                offer.Car.Vin,
		ProductionYear:     offer.Car.ProductionYear,
		Mileage:            offer.Car.Mileage,
		NumberOfDoors:      offer.Car.NumberOfDoors,
		NumberOfSeats:      offer.Car.NumberOfSeats,
		EnginePower:        offer.Car.EnginePower,
		EngineCapacity:     offer.Car.EngineCapacity,
		RegistrationNumber: offer.Car.RegistrationNumber,
		RegistrationDate:   offer.Car.RegistrationDate.Format(formats.DateLayout),
		Color:              offer.Car.Color,
		FuelType:           offer.Car.FuelType,
		Transmission:       offer.Car.Transmission,
		NumberOfGears:      offer.Car.NumberOfGears,
		Drive:              offer.Car.Drive,
		Brand:              offer.Car.Model.Manufacturer.Name,
		Model:              offer.Car.Model.Name,
		Status:             offer.Status,
		DateEnd:            endDate,
		BuyNowPrice:        buyNow,
		IsAuction:          offer.Auction != nil,
	}
}

func prepareAuctionValues(offer *models.SaleOffer) (*uint, *time.Time) {
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
	if !IsParamValid(dto.Margin, models.Margins) {
		return ErrInvalidMargin
	}
	return nil
}

func (dto *UpdateSaleOfferDTO) validateParams() error {
	if dto.Color != nil && !IsParamValid(*dto.Color, car_params.Colors) {
		return ErrInvalidColor
	}
	if dto.FuelType != nil && !IsParamValid(*dto.FuelType, car_params.Types) {
		return ErrInvalidFuelType
	}
	if dto.Transmission != nil && !IsParamValid(*dto.Transmission, car_params.Transmissions) {
		return ErrInvalidTransmission
	}
	if dto.Drive != nil && !IsParamValid(*dto.Drive, car_params.Drives) {
		return ErrInvalidDrive
	}
	if dto.Margin != nil && !IsParamValid(*dto.Margin, models.Margins) {
		return ErrInvalidMargin
	}
	return nil
}
