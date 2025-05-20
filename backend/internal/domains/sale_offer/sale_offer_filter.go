package sale_offer

import (
	"slices"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

// Constants

var OrderKeysMap = map[string]string{
	"Price":           "price",
	"Mileage":         "cars.mileage",
	"Production year": "cars.production_year",
	"Engine power":    "cars.engine_power",
	"Engine capacity": "cars.engine_capacity",
	"Date of issue":   "date_of_issue"}

type MinMax[T uint | string | time.Time] struct {
	Min *T `json:"min"`
	Max *T `json:"max"`
}

type FieldsConstraints struct {
	OfferTypes    []OfferType
	Manufacturers []string
	Colors        []car_params.Color
	Drives        []car_params.Drive
	FuelTypes     []car_params.FuelType
	Transmissions []car_params.Transmission
}

type OfferFilter struct {
	Pagination               pagination.PaginationRequest `json:"pagination"`
	UserID                   *uint                        `json:"user_id"`
	Query                    *string                      `json:"query"`
	OrderKey                 *string                      `json:"order_key"`
	IsOrderDesc              *bool                        `json:"is_order_desc"`
	LikedOnly                *bool                        `json:"liked_only"`
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
	CarRegistrationDateRange *MinMax[string]              `json:"car_registration_date_range"`
	OfferCreationDateRange   *MinMax[string]              `json:"offer_creation_date_range"`
	Constraints              FieldsConstraints            `json:"-"`
}

func NewOfferFilter() *OfferFilter {
	return &OfferFilter{Constraints: FieldsConstraints{
		OfferTypes:    OfferTypes,
		Colors:        car_params.Colors,
		Drives:        car_params.Drives,
		FuelTypes:     car_params.Types,
		Transmissions: car_params.Transmissions,
	}}
}

func (of *OfferFilter) ApplyOfferFilters(query *gorm.DB) (*gorm.DB, error) {
	if err := of.validateParams(); err != nil {
		return nil, err
	}
	query = applyUserFilter(query, of.UserID)
	query = applyOfferTypeFilter(query, of.OfferType)
	query = applyLikedOnlyFilter(query, of.LikedOnly, of.UserID)
	query = applyManufacturersFilter(query, of.Manufacturers)
	query = applyInSliceFilter(query, "cars.color", of.Colors)
	query = applyInSliceFilter(query, "cars.drive", of.Drives)
	query = applyInSliceFilter(query, "cars.fuel_type", of.FuelTypes)
	query = applyInSliceFilter(query, "cars.transmission", of.Transmissions)
	query = applyInRangeFilter(query, "price", of.PriceRange)
	query = applyInRangeFilter(query, "cars.mileage", of.MileageRange)
	query = applyInRangeFilter(query, "cars.production_year", of.YearRange)
	query = applyInRangeFilter(query, "cars.engine_power", of.EnginePowerRange)
	query = applyInRangeFilter(query, "cars.engine_capacity", of.EngineCapacityRange)
	query = applyDateInRangeFilter(query, "cars.registration_date", of.CarRegistrationDateRange)
	query = applyDateInRangeFilter(query, "date_of_issue", of.OfferCreationDateRange)
	query = applyOrderFilter(query, of.OrderKey, of.IsOrderDesc)
	return query, nil
}

func applyUserFilter(query *gorm.DB, userID *uint) *gorm.DB {
	if userID != nil {
		query = query.Where("sale_offers.user_id != ?", *userID)
	}
	return query
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

func applyManufacturersFilter(query *gorm.DB, values *[]string) *gorm.DB {
	if values != nil && len(*values) > 0 {
		query = query.
			Joins("JOIN models ON models.id = cars.model_id").
			Joins("JOIN manufacturers ON manufacturers.id = models.manufacturer_id").
			Where("manufacturers.name IN ?", *values)
	}
	return query
}

func applyLikedOnlyFilter(query *gorm.DB, likedOnly *bool, userID *uint) *gorm.DB {
	if likedOnly != nil && userID != nil {
		query = query.
			Joins("JOIN liked_offers ON liked_offers.offer_id = sale_offers.id").
			Where("liked_offers.user_id = ?", *userID)
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

func (of *OfferFilter) validateParams() error {
	validators := []func() error{of.validateEnums, of.validateRanges, of.validateDates, of.validateOrderKey}
	for _, validate := range validators {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}
func (of *OfferFilter) validateEnums() error {
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

func (of *OfferFilter) validateRanges() error {
	ranges := []*MinMax[uint]{of.PriceRange, of.YearRange, of.MileageRange, of.EnginePowerRange, of.EngineCapacityRange}
	for _, r := range ranges {
		if r != nil && !areMinMaxValidNumbers(*r) {
			return ErrInvalidRange
		}
	}
	return nil
}

func (of *OfferFilter) validateDates() error {
	datesRanges := []*MinMax[string]{of.CarRegistrationDateRange, of.OfferCreationDateRange}
	for _, r := range datesRanges {
		if err := validateDateRange(r); err != nil {
			return err
		}
	}
	return nil
}

func (of *OfferFilter) validateOrderKey() error {
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

func areMinMaxValidNumbers(minmax MinMax[uint]) bool {
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
