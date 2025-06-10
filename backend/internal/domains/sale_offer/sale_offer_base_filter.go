package sale_offer

import (
	"slices"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

var OrderKeysMap = map[string]string{
	"Price":           "price",
	"Mileage":         "mileage",
	"Production year": "production_year",
	"Engine power":    "engine_power",
	"Engine capacity": "engine_capacity",
	"Date of issue":   "date_of_issue"}

type MinMax[T uint | string | time.Time] struct {
	Min *T `json:"min"`
	Max *T `json:"max"`
}

type FieldsConstraints struct {
	OfferTypes    []OfferType
	Manufacturers []string
	Colors        []enums.Color
	Drives        []enums.Drive
	FuelTypes     []enums.FuelType
	Transmissions []enums.Transmission
}

type BaseOfferFilter struct {
	UserID                   *uint                 `json:"user_id"`
	Query                    *string               `json:"query"`
	OrderKey                 *string               `json:"order_key"`
	IsOrderDesc              *bool                 `json:"is_order_desc"`
	OfferType                *OfferType            `json:"offer_type"`
	Manufacturers            *[]string             `json:"manufacturers"`
	Colors                   *[]enums.Color        `json:"colors"`
	Drives                   *[]enums.Drive        `json:"drives"`
	FuelTypes                *[]enums.FuelType     `json:"fuel_types"`
	Transmissions            *[]enums.Transmission `json:"transmissions"`
	PriceRange               *MinMax[uint]         `json:"price_range"`
	MileageRange             *MinMax[uint]         `json:"mileage_range"`
	YearRange                *MinMax[uint]         `json:"year_range"`
	EnginePowerRange         *MinMax[uint]         `json:"engine_power_range"`
	EngineCapacityRange      *MinMax[uint]         `json:"engine_capacity_range"`
	CarRegistrationDateRange *MinMax[string]       `json:"car_registration_date_range"`
	OfferCreationDateRange   *MinMax[string]       `json:"offer_creation_date_range"`
	Constraints              FieldsConstraints     `json:"-"`
}

func NewOfferFilter() *BaseOfferFilter {
	return &BaseOfferFilter{Constraints: FieldsConstraints{
		OfferTypes:    OfferTypes,
		Colors:        enums.Colors,
		Drives:        enums.Drives,
		FuelTypes:     enums.Types,
		Transmissions: enums.Transmissions,
	}}
}

type OfferFilterRequest struct {
	Filter     BaseOfferFilter              `json:"filter"`
	PagRequest pagination.PaginationRequest `json:"pagination"`
}

func NewOfferFilterRequest() *OfferFilterRequest {
	return &OfferFilterRequest{Filter: *NewOfferFilter()}
}

func (of *BaseOfferFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	if err := of.validateParams(); err != nil {
		return nil, err
	}
	query = applyOfferTypeFilter(query, of.OfferType)
	query = applyInSliceFilter(query, "brand", of.Manufacturers)
	query = applyInSliceFilter(query, "color", of.Colors)
	query = applyInSliceFilter(query, "drive", of.Drives)
	query = applyInSliceFilter(query, "fuel_type", of.FuelTypes)
	query = applyInSliceFilter(query, "transmission", of.Transmissions)
	query = applyInRangeFilter(query, "price", of.PriceRange)
	query = applyInRangeFilter(query, "mileage", of.MileageRange)
	query = applyInRangeFilter(query, "production_year", of.YearRange)
	query = applyInRangeFilter(query, "engine_power", of.EnginePowerRange)
	query = applyInRangeFilter(query, "engine_capacity", of.EngineCapacityRange)
	query = applyDateInRangeFilter(query, "registration_date", of.CarRegistrationDateRange)
	query = applyDateInRangeFilter(query, "date_of_issue", of.OfferCreationDateRange)
	query = applyOrderFilter(query, of.OrderKey, of.IsOrderDesc)
	return query, nil
}

func (of *BaseOfferFilter) GetBase() *BaseOfferFilter {
	return of
}

func applyOfferTypeFilter(query *gorm.DB, offerType *OfferType) *gorm.DB {
	if offerType == nil {
		return query
	}
	switch *offerType {
	case REGULAR_OFFER:
		return query.Where("is_auction IS FALSE")
	case AUCTION:
		return query.Where("is_auction IS TRUE")
	default:
		return query
	}
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

func applyOrderFilter(query *gorm.DB, orderKey *string, isOrderDesc *bool) *gorm.DB {
	var orderDirection string
	if isOrderDesc == nil || *isOrderDesc {
		orderDirection = "DESC"
	} else {
		orderDirection = "ASC"
	}
	if orderKey != nil {
		return query.Order(OrderKeysMap[*orderKey] + " " + orderDirection + ", margin " + orderDirection)
	}
	return query.Order("margin " + orderDirection)
}

func (of *BaseOfferFilter) validateParams() error {
	validators := []func() error{of.validateEnums, of.validateRanges, of.validateDates, of.validateOrderKey}
	for _, validate := range validators {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}
func (of *BaseOfferFilter) validateEnums() error {
	if of.OfferType != nil && !IsParamValid(*of.OfferType, OfferTypes) {
		return ErrInvalidSaleOfferType
	}
	if of.Manufacturers != nil && !AreParamsValid(of.Manufacturers, &of.Constraints.Manufacturers) {
		return ErrInvalidManufacturer
	}
	if of.Colors != nil && !AreParamsValid(of.Colors, &of.Constraints.Colors) {
		return ErrInvalidColor
	}
	if of.Drives != nil && !AreParamsValid(of.Drives, &of.Constraints.Drives) {
		return ErrInvalidDrive
	}
	if of.FuelTypes != nil && !AreParamsValid(of.FuelTypes, &of.Constraints.FuelTypes) {
		return ErrInvalidFuelType
	}
	if of.Transmissions != nil && !AreParamsValid(of.Transmissions, &of.Constraints.Transmissions) {
		return ErrInvalidTransmission
	}
	return nil
}

func (of *BaseOfferFilter) validateRanges() error {
	ranges := []*MinMax[uint]{of.PriceRange, of.YearRange, of.MileageRange, of.EnginePowerRange, of.EngineCapacityRange}
	for _, r := range ranges {
		if r != nil && !isMinMaxValidNumbers(*r) {
			return ErrInvalidRange
		}
	}
	return nil
}

func (of *BaseOfferFilter) validateDates() error {
	datesRanges := []*MinMax[string]{of.CarRegistrationDateRange, of.OfferCreationDateRange}
	for _, r := range datesRanges {
		if err := validateDateRange(r); err != nil {
			return err
		}
	}
	return nil
}

func (of *BaseOfferFilter) validateOrderKey() error {
	if of.OrderKey != nil && !slices.Contains(GetKeysFromMap(OrderKeysMap), *of.OrderKey) {
		return ErrInvalidOrderKey
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
	var minValue, maxValue *time.Time
	var err error
	if minmax.Min != nil {
		minValue, err = ParseDate(*minmax.Min)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
	}
	if minmax.Max != nil {
		maxValue, err = ParseDate(*minmax.Max)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
	}
	return &MinMax[time.Time]{Min: minValue, Max: maxValue}, nil
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

func GetKeysFromMap[T comparable](m map[T]T) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
