package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"gorm.io/gorm"
)

type MinMax struct {
	Min *uint `json:"min"`
	Max *uint `json:"max"`
}

type OfferFilter struct {
	OrderKeys             *[]string                  `json:"order_keys"`
	IsOrderDesc           *bool                      `json:"is_order_desc"`
	OfferType             *string                    `json:"offer_type"`
	Manufacturers         *[]string                  `json:"manufacturers"`
	Colors                *[]car_params.Color        `json:"colors"`
	Drives                *[]car_params.Drive        `json:"drives"`
	FuelTypes             *[]car_params.FuelType     `json:"fuel_types"`
	Transmissions         *[]car_params.Transmission `json:"transmissions"`
	PriceRange            *MinMax                    `json:"price_range"`
	MileageRange          *MinMax                    `json:"mileage_range"`
	YearRange             *MinMax                    `json:"year_range"`
	EnginePowerRange      *MinMax                    `json:"engine_power_range"`
	EngineCapacityRange   *MinMax                    `json:"engine_capacity_range"`
	RegistrationDateRagne *MinMax                    `json:"registration_date_range"`
}

func (of *OfferFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	if err := of.validateParams(); err != nil {
		return nil, err
	}
	query = applyInSliceFilter(query, "cars.color", of.Colors)
	query = applyInSliceFilter(query, "cars.drive", of.Drives)
	query = applyInSliceFilter(query, "cars.fuel_type", of.FuelTypes)
	query = applyInSliceFilter(query, "cars.transmission", of.Transmissions)
	query = applyInRangeFilter(query, "price", of.PriceRange)
	query = applyInRangeFilter(query, "cars.mileage", of.MileageRange)
	query = applyInRangeFilter(query, "cars.engine_power", of.EnginePowerRange)
	query = applyInRangeFilter(query, "cars.engine_capacity", of.EngineCapacityRange)
	query = applyInRangeFilter(query, "cars.registration_date", of.RegistrationDateRagne)
	return query, nil
}

func applyInSliceFilter[T any](query *gorm.DB, column string, values *[]T) *gorm.DB {

	if values != nil && len(*values) > 0 {
		query = query.Where(column+" IN ?", *values)
	}
	return query
}

func applyInRangeFilter(query *gorm.DB, column string, minmax *MinMax) *gorm.DB {
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
	ranges := []*MinMax{of.PriceRange, of.YearRange, of.EnginePowerRange, of.EngineCapacityRange}
	for _, r := range ranges {
		if r != nil && !r.isMinMaxValid() {
			return ErrInvalidRange
		}
	}
	return nil
}

func areParamsValid[T comparable](params *[]T, validParams *[]T) bool {
	for _, param := range *params {
		if !IsParamValid(param, *validParams) {
			return false
		}
	}
	return true
}

func (mm *MinMax) isMinMaxValid() bool {
	if mm.Min != nil && mm.Max != nil {
		return *mm.Max >= *mm.Min
	}
	return true
}
