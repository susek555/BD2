package sale_offer_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

// ----------------
// Pagination tests
// ----------------

func TestGetFiltered_PaginationNegaitvePage(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: -1, PageSize: 8}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, pagination.ErrPageOutOfRange)
}

func TestGetFiltered_PaginationZeroPage(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 0, PageSize: 8}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, pagination.ErrPageOutOfRange)
}

func TestGetFiltered_PaginationPositivePage(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 1, PageSize: 8}
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_PaginationPageSizeNegative(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 1, PageSize: -1}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, pagination.ErrNegativePageSize)
}

func TestGetFiltered_PaginationPageSizeZero(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 1, PageSize: 0}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, pagination.ErrNegativePageSize)
}

func TestFiltered_PaginationPageSizePositive(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 1, PageSize: 8}
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestFiltered_PaginationPageOutOfRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 2, PageSize: 8}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, pagination.ErrPageOutOfRange)
}

func TestGetFiltered_PaginationSingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 1, PageSize: 8}
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PaginationMultipleRecordsBounded(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i)))
	}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 1, PageSize: 3}
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 3)
}

func TestGetFiltered_PaginationMultipleRecordsNotBounded(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i)))
	}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = pagination.PaginationRequest{Page: 2, PageSize: 3}
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 2)
}

//------------------------
// Invalid arguments tests
// -----------------------

func TestGetFiltered_InvalidManufacturer(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Manufacturers = &[]string{"invalid"}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidManufacturer)
}

func TestGetFiltered_InvalidOfferType(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType("invalid")
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &offer
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidSaleOfferType)
}

func TestGetFiltered_InvalidColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Colors = &[]car_params.Color{"invaid"}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidColor)
}

func TestGetFiltered_InvalidDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Drives = &[]car_params.Drive{"invaid"}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDrive)
}

func TestGetFiltered_InvalidFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.FuelTypes = &[]car_params.FuelType{"invaid"}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidFuelType)
}

func TestGetFiltered_InvalidTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Transmissions = &[]car_params.Transmission{"invaid"}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidTransmission)
}

func TestGetFiltered_InvalidPriceRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidPriceRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidMileageRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidMileageRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidYearRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidYearRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEnginePowerRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEnginePowerRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEngineCapacityRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEngineCapacityRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidCarRegistrationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022-01-01"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidCarRegistrationDateRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-01"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidCarRegistrationDateFormat(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022/01/01"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDateFromat)
}

func TestGetFiltered_InvalidOfferCreationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022-01-01"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidOfferCreationDateRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-01"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidOfferCreationDateFormat(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022/01/01"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	_, _, err := repo.GetFiltered(filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDateFromat)
}

// ---------------------
// Valid arguments tests
// ---------------------

func TestGetFiltered_ValidOfferTypeAutcion(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType(sale_offer.AUCTION)
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &offer
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferTypeRegularOffer(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType(sale_offer.REGULAR_OFFER)
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &offer
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidOfferTypeBoth(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType(sale_offer.BOTH)
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &offer
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidManufacturer(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Constriants.Manufacturers = manufacturers
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Colors = &car_params.Colors
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Drives = &car_params.Drives
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.FuelTypes = &car_params.Types
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Transmissions = &car_params.Transmissions
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidEnginePowerRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEnginePowerRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidEnginePowerRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEnginePowerRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-02"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := "2023-01-02"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-02"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := "2023-01-02"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
}

// ----------------------------------
// Retrieving filtering results tests
// ----------------------------------

func TestGetFiltered_NoFilterEmptyDB(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_NoFilter(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferTypeRegularOffer(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	regularOffer := sale_offer.REGULAR_OFFER
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &regularOffer
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferTypeRegularOfferAuctionInDB(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithAuction(time.Now(), 0))}
	repo, _ := setupDB(offers)
	regularOffer := sale_offer.REGULAR_OFFER
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &regularOffer
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferTypeAuction(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithAuction(time.Now(), 0))}
	repo, _ := setupDB(offers)
	auction := sale_offer.AUCTION
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &auction
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferTypeAuctionRegularOfferInDB(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	auction := sale_offer.AUCTION
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &auction
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferTypeBoth(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithAuction(time.Now(), 0)), *CreateOffer(2)}
	repo, _ := setupDB(offers)
	both := sale_offer.BOTH
	filter := sale_offer.NewOfferFilter()
	filter.OfferType = &both
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_SingleManufacturer(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Constriants.Manufacturers = manufacturers
	filter.Manufacturers = &[]string{"Audi"}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleManufacturers(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1), *CreateOffer(2)}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Constriants.Manufacturers = manufacturers
	filter.Manufacturers = &[]string{"Audi", "BMW"}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_NoMatchingManufacturer(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Constriants.Manufacturers = manufacturers
	filter.Manufacturers = &[]string{"BMW"}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Color", car_params.RED))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Colors = &[]car_params.Color{car_params.RED}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleColors(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Color", car_params.RED)), *CreateOffer(2, WithCarField("Color", car_params.BLUE))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Colors = &[]car_params.Color{car_params.RED, car_params.BLUE}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_NoMatchingColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Color", car_params.RED))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Colors = &[]car_params.Color{car_params.GREEN}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Drive", car_params.FWD))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Drives = &[]car_params.Drive{car_params.FWD}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}
func TestGetFiltered_MultipleDrives(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Drive", car_params.FWD)), *CreateOffer(2, WithCarField("Drive", car_params.RWD))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Drives = &[]car_params.Drive{car_params.FWD, car_params.RWD}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferNoMatchingDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Drive", car_params.FWD))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Drives = &[]car_params.Drive{car_params.AWD}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("FuelType", car_params.PETROL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.FuelTypes = &[]car_params.FuelType{car_params.PETROL}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleFuelTypes(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("FuelType", car_params.PETROL)), *CreateOffer(2, WithCarField("FuelType", car_params.DIESEL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.FuelTypes = &[]car_params.FuelType{car_params.PETROL, car_params.DIESEL}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferNoMatchingFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("FuelType", car_params.PETROL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.FuelTypes = &[]car_params.FuelType{car_params.ELECTRIC}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Transmission", car_params.AUTOMATIC))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Transmissions = &[]car_params.Transmission{car_params.AUTOMATIC}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleTransmissions(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Transmission", car_params.AUTOMATIC)), *CreateOffer(2, WithCarField("Transmission", car_params.MANUAL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Transmissions = &[]car_params.Transmission{car_params.AUTOMATIC, car_params.MANUAL}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferNoMatchingTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Transmission", car_params.AUTOMATIC))}
	repo, _ := setupDB(offers)
	filter := sale_offer.NewOfferFilter()
	filter.Transmissions = &[]car_params.Transmission{car_params.MANUAL}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_PriceInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_PriceLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetfiltered_PriceUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_MileageLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_MileageUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	min := uint(2025)
	max := uint(2026)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2024)))}
	repo, _ := setupDB(offers)
	min := uint(2023)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2024)))}
	repo, _ := setupDB(offers)
	max := uint(2025)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	max := uint(2024)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_YearLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2023)))}
	repo, _ := setupDB(offers)
	min := uint(2024)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_YearUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	max := uint(2025)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	min := uint(2025)
	filter := sale_offer.NewOfferFilter()
	filter.YearRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EnginePowerLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EnginePowerUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EngineCapacityLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EngineCapacityUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.NewOfferFilter()
	filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateInRange(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	max := "2025-05-14"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateInRangeMinProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateInRangeMaxProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-14"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateGreater(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-12"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_CarRegistrationDateLower(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-14"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_CarRegistratoinDateUpperBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-13"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}
func TestGetFiltered_CarRegistrationDateLowerBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-13"
	filter := sale_offer.NewOfferFilter()
	filter.CarRegistrationDateRagne = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateInRange(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	max := "2025-05-14"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateInRangeMinProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateInRangeMaxProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-14"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateGreater(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-12"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferCreationDateLower(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-14"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferCreationDateUpperBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-13"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateLowerBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*CreateOffer(1, WithOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-13"
	filter := sale_offer.NewOfferFilter()
	filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min, Max: nil}
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

// ----------------
// Order by tests
// ----------------

func TestGetFiltered_OrderByPriceNoRecords(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	key := "Price"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByPriceSingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	key := "Price"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByPriceMultipleRecordsDesc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithOfferField("Price", uint(i))))
	}
	repo, _ := setupDB(offers)
	key := "Price"
	trueStm := true
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &trueStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Price, uint(len(offers)-i))
	}
}

func TestGetFiltered_OrderByPriceMultipleRecordsAsc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithOfferField("Price", uint(i))))
	}
	repo, _ := setupDB(offers)
	key := "Price"
	falseStm := false
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &falseStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Price, uint(i+1))
	}
}

func TestGetFiltered_OrderByMileageNoRecords(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	key := "Mileage"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByMileageSingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	key := "Mileage"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByMileageMultipleRecordsDesc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("Mileage", uint(i))))
	}
	repo, _ := setupDB(offers)
	key := "Mileage"
	trueStm := true
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &trueStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.Mileage, uint(len(offers)-i))
	}
}

func TestGetFiltered_OrderByMileageMultipleRecordsAsc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("Mileage", uint(i))))
	}
	repo, _ := setupDB(offers)
	key := "Mileage"
	falseStm := false
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &falseStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.Mileage, uint(i+1))
	}
}

func TestGetFiltered_OrderByYearNoRecords(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	key := "Production year"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByYearSingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	key := "Production year"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByYearMultipleRecordsDesc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("ProductionYear", uint(i+2000))))
	}
	repo, _ := setupDB(offers)
	key := "Production year"
	trueStm := true
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &trueStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.ProductionYear, uint(2000+len(offers)-i))
	}
}

func TestGetFiltered_OrderByYearMultipleRecordsAsc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("ProductionYear", uint(i+2000))))
	}
	repo, _ := setupDB(offers)
	key := "Production year"
	falseStm := false
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &falseStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.ProductionYear, uint(2000+i+1))
	}
}

func TestGetFiltered_OrderByEnginePowerNoRecords(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	key := "Engine power"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByEnginePowerSingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	key := "Engine power"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByEnginePowerMultipleRecordsDesc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("EnginePower", uint(i+100))))
	}
	repo, _ := setupDB(offers)
	key := "Engine power"
	trueStm := true
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &trueStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.EnginePower, uint(100+len(offers)-i))
	}
}

func TestGetFiltered_OrderByEnginePowerMultipleRecordsAsc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("EnginePower", uint(i+100))))
	}
	repo, _ := setupDB(offers)
	key := "Engine power"
	falseStm := false
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &falseStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.EnginePower, uint(100+i+1))
	}
}

func TestGetFiltered_OrderByEngineCapacityNoRecords(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	key := "Engine capacity"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByEngineCapacitySingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	key := "Engine capacity"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByEngineCapacityMultipleRecordsDesc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("EngineCapacity", uint(i+1000))))
	}
	repo, _ := setupDB(offers)
	key := "Engine capacity"
	trueStm := true
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &trueStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.EngineCapacity, uint(1000+len(offers)-i))
	}
}
func TestGetFiltered_OrderByEngineCapacityMultipleRecordsAsc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithCarField("EngineCapacity", uint(1000+i))))
	}
	repo, _ := setupDB(offers)
	key := "Engine capacity"
	falseStm := false
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &falseStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Car.EngineCapacity, uint(1000+i+1))
	}
}

func TestGetFiltered_OrderByDateOfIssueNoRecords(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	key := "Date of issue"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByDateOfIssueSingleRecord(t *testing.T) {
	offers := []sale_offer.SaleOffer{*CreateOffer(1)}
	repo, _ := setupDB(offers)
	key := "Date of issue"
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OrderByDateOfIssueMultipleRecordsDesc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithOfferField("DateOfIssue", time.Date(2025, 5, 1+i, 0, 0, 0, 0, time.UTC))))
	}
	repo, _ := setupDB(offers)
	key := "Date of issue"
	trueStm := true
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &trueStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].DateOfIssue, time.Date(2025, 5, 1+(len(offers)-i), 0, 0, 0, 0, time.UTC))
	}
}

func TestGetFiltered_OrderByDateOfIssueMultipleRecordsAsc(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	for i := 1; i <= 5; i++ {
		offers = append(offers, *CreateOffer(uint(i), WithOfferField("DateOfIssue", time.Date(2025, 5, 1+i, 0, 0, 0, 0, time.UTC))))
	}
	repo, _ := setupDB(offers)
	key := "Date of issue"
	falseStm := false
	filter := sale_offer.NewOfferFilter()
	filter.OrderKey = &key
	filter.IsOrderDesc = &falseStm
	filter.Pagination = *GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(filter)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].DateOfIssue, time.Date(2025, 5, 1+i+1, 0, 0, 0, 0, time.UTC))
	}
}
