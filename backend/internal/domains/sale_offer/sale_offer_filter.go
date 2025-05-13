package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type OfferType string

const (
	REGULAR_OFFER OfferType = "Regular offer"
	AUCTION       OfferType = "Auction"
	BOTH          OfferType = "Both"
)

var OfferTypes = []OfferType{REGULAR_OFFER, AUCTION, BOTH}

type MinMax[T uint | string | time.Time] struct {
	Min *T `json:"min"`
	Max *T `json:"max"`
}

type OfferFilter struct {
	Pagination               pagination.PaginationRequest `json:"pagination"`
	Query                    *string                      `json:"query"`
	OrderKey                 *string                      `json:"order_key"`
	IsOrderDesc              *bool                        `json:"is_order_desc"`
	OfferType                *OfferType                   `json:"offer_type"`
	Manufacturers            *[]string                    `json:"manufacturers"`
	Colors                   *[]car_params.Color          `json:"colors"`
	Drives                   *[]car_params.Drive          `json:"drives"`
	FuelTypes                *[]car_params.FuelType       `json:"fuel_types"`
	Transmissions            *[]car_params.Transmission   `json:"transmissions"`
	PriceRange               *MinMax[uint]                `json:"price_range"`
	MileageRange             *MinMax[uint]                `json:"mileage_range"`
	YearRange                *MinMax[uint]                `json:"year_range"`
	EnginePowerRange         *MinMax[uint]                `json:"engine_power_range"`
	EngineCapacityRange      *MinMax[uint]                `json:"engine_capacity_range"`
	CarRegistrationDateRagne *MinMax[string]              `json:"car_registration_date_range"`
	OfferCreationDateRange   *MinMax[string]              `json:"offer_creation_date_range"`
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
	query = applyInRangeFilter(query, "cars.production_year", of.YearRange)
	query = applyInRangeFilter(query, "cars.engine_power", of.EnginePowerRange)
	query = applyInRangeFilter(query, "cars.engine_capacity", of.EngineCapacityRange)
	query = applyDateInRangeFilter(query, "cars.registration_date", of.CarRegistrationDateRagne)
	query = applyDateInRangeFilter(query, "date_of_issue", of.OfferCreationDateRange)
	query = applyOrderFilter(query, of.OrderKey)
	return query, nil
}

func applyOfferTypeFilter(query *gorm.DB, offerType *OfferType) *gorm.DB {
	if offerType == nil {
		return query
	}
	switch *offerType {
	case AUCTION:
		return query.Where("auctions.offer_id IS NOT NULL")
	case REGULAR_OFFER:
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

func applyInRangeFilter[T uint | time.Time](query *gorm.DB, column string, minmax *MinMax[T]) *gorm.DB {
	if minmax == nil {
		return query
	}
	if minmax.Max != nil {
		query = query.Where(column+" <= ?", *minmax.Max)
	}
	if minmax.Min != nil {
		query = query.Where(column+" >= ?", *minmax.Min)
	}
	return query
}

func applyDateInRangeFilter(query *gorm.DB, column string, minmax *MinMax[string]) *gorm.DB {
	if minmax != nil {
		dates, _ := parseDateRange(minmax)
		return applyInRangeFilter(query, column, dates)
	}
	return query
}

func applyOrderFilter(query *gorm.DB, orderKey *string) *gorm.DB {
	if orderKey != nil {
		return query.Order(*orderKey)
	}
	return query
}

func (of *OfferFilter) validateParams() error {
	if err := of.validateEnums(); err != nil {
		return err
	}
	if err := of.validateDates(); err != nil {
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
	if of.Colors != nil && !AreParamsValid(of.Colors, &car_params.Colors) {
		return ErrInvalidColor
	}
	if of.Drives != nil && !AreParamsValid(of.Drives, &car_params.Drives) {
		return ErrInvalidDrive
	}
	if of.FuelTypes != nil && !AreParamsValid(of.FuelTypes, &car_params.Types) {
		return ErrInvalidFuelType
	}
	if of.Transmissions != nil && !AreParamsValid(of.Transmissions, &car_params.Transmissions) {
		return ErrInvalidTransmission
	}
	return nil
}

func (of *OfferFilter) validateRanges() error {
	ranges := []*MinMax[uint]{of.PriceRange, of.YearRange, of.MileageRange, of.EnginePowerRange, of.EngineCapacityRange}
	for _, r := range ranges {
		if r != nil && !isMinMaxValidNumbers(*r) {
			return ErrInvalidRange
		}
	}
	return nil
}

func (of *OfferFilter) validateDates() error {
	if err := validateDateRange(of.CarRegistrationDateRagne); err != nil {
		return err
	}
	if err := validateDateRange(of.OfferCreationDateRange); err != nil {
		return err
	}
	return nil
}

func validateDateRange(minmax *MinMax[string]) error {
	if minmax == nil {
		return nil
	}
	dates, err := parseDateRange(minmax)
	if err != nil {
		return err
	}
	if !isMinMaxValidDates(*dates) {
		return ErrInvalidRange
	}
	return nil
}

func parseDateRange(minmax *MinMax[string]) (*MinMax[time.Time], error) {
	var min, max *time.Time
	var err error
	if minmax.Min != nil {
		min, err = ParseDate(*minmax.Min)
		if err != nil {
			return nil, ErrInvalidDateFromat
		}
	}
	if minmax.Max != nil {
		max, err = ParseDate(*minmax.Max)
		if err != nil {
			return nil, ErrInvalidDateFromat
		}
	}
	return &MinMax[time.Time]{Min: min, Max: max}, nil
}

func isMinMaxValidNumbers(minmax MinMax[uint]) bool {
	if minmax.Min != nil && minmax.Max != nil {
		return *minmax.Max > *minmax.Min
	}
	return true
}

func isMinMaxValidDates(minmax MinMax[time.Time]) bool {
	if minmax.Min != nil && minmax.Max != nil {
		return (*minmax.Max).After(*minmax.Min)
	}
	return true
}
