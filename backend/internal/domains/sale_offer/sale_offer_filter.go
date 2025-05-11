package sale_offer

import (
	"cmp"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"gorm.io/gorm"
)

type OfferType string

const (
	REGULAR_OFFER OfferType = "Regular offer"
	AUCTION       OfferType = "Auction"
	BOTH          OfferType = "Both"
)

var OfferTypes = []OfferType{REGULAR_OFFER, AUCTION, BOTH}

type MinMax[T cmp.Ordered] struct {
	Min *T `json:"min"`
	Max *T `json:"max"`
}

type OfferFilter struct {
	OrderKeys                *[]string                  `json:"order_keys"`
	IsOrderDesc              *bool                      `json:"is_order_desc"`
	OfferType                *OfferType                 `json:"offer_type"`
	Manufacturers            *[]string                  `json:"manufacturers"`
	Colors                   *[]car_params.Color        `json:"colors"`
	Drives                   *[]car_params.Drive        `json:"drives"`
	FuelTypes                *[]car_params.FuelType     `json:"fuel_types"`
	Transmissions            *[]car_params.Transmission `json:"transmissions"`
	PriceRange               *MinMax[uint]              `json:"price_range"`
	MileageRange             *MinMax[uint]              `json:"mileage_range"`
	YearRange                *MinMax[uint]              `json:"year_range"`
	EnginePowerRange         *MinMax[uint]              `json:"engine_power_range"`
	EngineCapacityRange      *MinMax[uint]              `json:"engine_capacity_range"`
	CarRegistrationDateRagne *MinMax[string]            `json:"car_registration_date_range"`
	OfferCreationDateRange   *MinMax[string]            `json:"offer_creation_date_range"`
}

func (of *OfferFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	if err := of.validateParams(); err != nil {
		return nil, err
	}
	query = applyOfferTypeFilter(query, of.OfferType)
	query = applyManufacturesrsFilter(query, of.Manufacturers)
	query = applyInSliceFilter(query, "cars.color", of.Colors)
	query = applyInSliceFilter(query, "cars.drive", of.Drives)
	query = applyInSliceFilter(query, "cars.fuel_type", of.FuelTypes)
	query = applyInSliceFilter(query, "cars.transmission", of.Transmissions)
	query = applyInRangeFilter(query, "price", of.PriceRange)
	query = applyInRangeFilter(query, "cars.mileage", of.MileageRange)
	query = applyInRangeFilter(query, "cars.engine_power", of.EnginePowerRange)
	query = applyInRangeFilter(query, "cars.engine_capacity", of.EngineCapacityRange)
	query = applyDateInRangeFilter(query, "cars.registration_date", of.CarRegistrationDateRagne)
	query = applyDateInRangeFilter(query, "date_of_issue", of.OfferCreationDateRange)
	return query, nil
}

func applyOfferTypeFilter(query *gorm.DB, offerType *OfferType) *gorm.DB {
	if offerType == nil {
		return query
	}
	switch *offerType {
	case "Auction":
		return query.Where("auctions.offer_id IS NOT NULL")
	case "Regular offer":
		return query.Where("auctions.offer_id IS NULL")
	default:
		return query
	}
}

func applyManufacturesrsFilter(query *gorm.DB, values *[]string) *gorm.DB {
	if values != nil && len(*values) > 0 {
		query = query.
			Joins("JOIN models ON models.id = cars.model_id").
			Joins("JOIN manufacturers ON manufacturers.id = models.manufacturer_id").
			Where("manufacturers.name IN ?", *values)
	}
	return query
}

func applyInSliceFilter[T any](query *gorm.DB, column string, values *[]T) *gorm.DB {
	if values != nil && len(*values) > 0 {
		query = query.Where(column+" IN ?", *values)
	}
	return query
}

func applyInRangeFilter[T cmp.Ordered](query *gorm.DB, column string, minmax *MinMax[T]) *gorm.DB {
	if minmax == nil {
		return query
	}
	if minmax.Max != nil {
		query = query.Where(column+" < ?", *minmax.Max)
	}
	if minmax.Min != nil {
		query = query.Where(column+" > ?", *minmax.Min)
	}
	return query
}

func applyDateInRangeFilter(query *gorm.DB, column string, minmax *MinMax[string]) *gorm.DB {
	dates, _ := parseDateRange(minmax)
	return applyInRangeFilter(query, column, dates)
}

func (of *OfferFilter) validateParams() error {
	if err := of.validateEnums(); err != nil {
		return err
	}
	if err := of.validateRanges(); err != nil {
		return err
	}
	return nil
}

func (of *OfferFilter) validateEnums() error {
	if of.OfferType != nil && !IsParamValid(*of.OfferType, OfferTypes) {
		return ErrInvalidSaleOfferType
	}
	if of.Colors != nil && !areParamsValid(of.Colors, &car_params.Colors) {
		return ErrInvalidColor
	}
	if of.Drives != nil && !areParamsValid(of.Drives, &car_params.Drives) {
		return ErrInvalidDrive
	}
	if of.FuelTypes != nil && !areParamsValid(of.FuelTypes, &car_params.Types) {
		return ErrInvalidFuelType
	}
	if of.Transmissions != nil && !areParamsValid(of.Transmissions, &car_params.Transmissions) {
		return ErrInvalidTransmission
	}
	return nil
}

func (of *OfferFilter) validateRanges() error {
	creationDates, err := parseDateRange(of.OfferCreationDateRange)
	if err != nil {
		return err
	}
	registrationDates, err := parseDateRange(of.CarRegistrationDateRagne)
	if err != nil {
		return err
	}
	ranges := []*MinMax[uint]{of.PriceRange, of.YearRange, of.EnginePowerRange, of.EngineCapacityRange, creationDates, registrationDates}
	for _, r := range ranges {
		if r != nil && !r.isMinMaxValid() {
			return ErrInvalidRange
		}
	}
	return nil
}

func parseDateRange(minmax *MinMax[string]) (*MinMax[uint], error) {
	min, err := parseDate(*minmax.Min)
	if err != nil {
		return nil, err
	}
	max, err := parseDate(*minmax.Max)
	if err != nil {
		return nil, err
	}
	minUint := uint(min.Unix())
	maxUint := uint(max.Unix())
	return &MinMax[uint]{Min: &minUint, Max: &maxUint}, nil
}

func areParamsValid[T comparable](params *[]T, validParams *[]T) bool {
	for _, param := range *params {
		if !IsParamValid(param, *validParams) {
			return false
		}
	}
	return true
}

func (mm *MinMax[T]) isMinMaxValid() bool {
	if mm.Min != nil && mm.Max != nil {
		return *mm.Max > *mm.Min
	}
	return true
}
