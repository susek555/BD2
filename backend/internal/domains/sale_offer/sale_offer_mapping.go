package sale_offer

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
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
	date, _ := ParseDate(dto.RegistrationDate)
	offer := &models.SaleOffer{Car: &models.Car{}}
	if err := copier.Copy(offer, dto); err != nil {
		return nil, err
	}
	if err := copier.Copy(offer.Car, dto); err != nil {
		return nil, err
	}
	offer.DateOfIssue = time.Now().UTC()
	offer.Car.RegistrationDate = *date
	offer.Status = enums.PENDING
	return offer, nil
}

func (dto *UpdateSaleOfferDTO) UpdateOfferFromDTO(offer *models.SaleOffer) (*models.SaleOffer, error) {
	if err := dto.validateParams(); err != nil {
		return nil, err
	}
	if dto.RegistrationDate != nil {
		date, _ := ParseDate(*dto.RegistrationDate)
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

func MapViewToDTO(offerView *views.SaleOfferView) *RetrieveSaleOfferDTO {
	dto := &RetrieveSaleOfferDTO{}
	_ = copier.Copy(dto, offerView)
	dto.Name = offerView.Brand + " " + offerView.Model
	return dto
}

func MapViewToDetailedDTO(offerView *views.SaleOfferView) *RetrieveDetailedSaleOfferDTO {
	dto := &RetrieveDetailedSaleOfferDTO{}
	_ = copier.Copy(dto, offerView)
	dto.DateOfIssue = offerView.DateOfIssue.Format(formats.DateLayout)
	dto.RegistrationDate = offerView.RegistrationDate.Format(formats.DateLayout)
	if dto.DateEnd != nil {
		date := offerView.DateEnd.Format(formats.DateTimeLayout)
		dto.DateEnd = &date
	}
	return dto
}

func (dto *CreateSaleOfferDTO) validateParams() error {
	if !IsParamValid(dto.Color, enums.Colors) {
		return ErrInvalidColor
	}
	if !IsParamValid(dto.FuelType, enums.Types) {
		return ErrInvalidFuelType
	}
	if !IsParamValid(dto.Transmission, enums.Transmissions) {
		return ErrInvalidTransmission
	}
	if !IsParamValid(dto.Drive, enums.Drives) {
		return ErrInvalidDrive
	}
	if !IsParamValid(dto.Margin, enums.Margins) {
		return ErrInvalidMargin
	}
	if dto.ProductionYear > uint(time.Now().Year()) || dto.ProductionYear < 1886 {
		return ErrInvalidProductionYear
	}
	d, err := ParseDate(dto.RegistrationDate)
	if err != nil {
		return err
	}
	if d.After(time.Now()) {
		return ErrInvalidRegistrationDate
	}
	return nil
}

func (dto *UpdateSaleOfferDTO) validateParams() error {
	if dto.Color != nil && !IsParamValid(*dto.Color, enums.Colors) {
		return ErrInvalidColor
	}
	if dto.FuelType != nil && !IsParamValid(*dto.FuelType, enums.Types) {
		return ErrInvalidFuelType
	}
	if dto.Transmission != nil && !IsParamValid(*dto.Transmission, enums.Transmissions) {
		return ErrInvalidTransmission
	}
	if dto.Drive != nil && !IsParamValid(*dto.Drive, enums.Drives) {
		return ErrInvalidDrive
	}
	if dto.Margin != nil && !IsParamValid(*dto.Margin, enums.Margins) {
		return ErrInvalidMargin
	}
	if dto.ProductionYear != nil && (*dto.ProductionYear > uint(time.Now().Year()) || *dto.ProductionYear < 1886) {
		return ErrInvalidProductionYear
	}
	if dto.RegistrationDate != nil {
		d, err := ParseDate(*dto.RegistrationDate)
		if err != nil {
			return err
		}
		if d.After(time.Now()) {
			return ErrInvalidRegistrationDate
		}
	}
	return nil
}
