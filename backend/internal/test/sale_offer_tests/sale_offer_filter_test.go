package sale_offer_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

// ----------------
// Pagination tests
// ----------------

var DB, _ = setupDB()

func TestGetFiltered_PaginationNegativePage(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: -1, PageSize: 8}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, pagination.ErrPageOutOfRange)
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationZeroPage(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 0, PageSize: 8}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, pagination.ErrPageOutOfRange)
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationPositivePage(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 1, PageSize: 8}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationPageSizeNegative(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 1, PageSize: -1}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, pagination.ErrNegativePageSize)
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationPageSizeZero(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 1, PageSize: 0}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, pagination.ErrNegativePageSize)
	u.CleanDB(DB)
}

func TestFiltered_PaginationPageSizePositive(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 1, PageSize: 8}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestFiltered_PaginationPageOutOfRange(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 2, PageSize: 8}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, pagination.ErrPageOutOfRange)
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 1, PageSize: 8}
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationMultipleRecordsBounded(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *createOffer(uint(i)))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 1, PageSize: 3}
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 3)
	u.CleanDB(DB)
}

func TestGetFiltered_PaginationMultipleRecordsNotBounded(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *createOffer(uint(i)))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = pagination.PaginationRequest{Page: 2, PageSize: 3}
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 2)
	u.CleanDB(DB)
}

//------------------------
// Invalid arguments tests
// -----------------------

func TestGetFiltered_InvalidManufacturer(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Manufacturers = &[]string{"invalid"}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidManufacturer)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidOfferType(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	offer := sale_offer.OfferType("invalid")
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &offer
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidSaleOfferType)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidColor(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Colors = &[]enums.Color{"invalid"}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidColor)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidDrive(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Drives = &[]enums.Drive{"invalid"}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDrive)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidFuelType(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.FuelTypes = &[]enums.FuelType{"invalid"}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidFuelType)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidTransmission(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Transmissions = &[]enums.Transmission{"invalid"}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidTransmission)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidPriceRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidPriceRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidMileageRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidMileageRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidYearRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidYearRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidEnginePowerRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidEnginePowerRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidEngineCapacityRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidEngineCapacityRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidCarRegistrationDateRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2022-01-01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidCarRegistrationDateRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2023-01-01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidCarRegistrationDateFormat(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2022/01/01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDateFormat)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidOfferCreationDateRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2022-01-01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidOfferCreationDateRangeBothValues(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2023-01-01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
	u.CleanDB(DB)
}

func TestGetFiltered_InvalidOfferCreationDateFormat(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2022/01/01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDateFormat)
	u.CleanDB(DB)
}

// ---------------------
// Valid arguments tests
// ---------------------

func TestGetFiltered_ValidOfferTypeAuction(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	offer := sale_offer.AUCTION
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &offer
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidOfferTypeRegularOffer(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	offer := sale_offer.REGULAR_OFFER
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &offer
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidOfferTypeBoth(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	offer := sale_offer.BOTH
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &offer
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidManufacturer(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(MANUFACTURERS, manufacturer.MapToName)
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidColor(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Colors = &enums.Colors
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidDrive(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Drives = &enums.Drives
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidFuelType(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.FuelTypes = &enums.Types
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidTransmission(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Transmissions = &enums.Transmissions
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidPriceRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidPriceRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidPriceRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidPriceRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidMileageRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidMileageRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidMileageRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidMileageRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidYearRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidYearRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidYearRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidYearRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidEnginePowerRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidEnginePowerRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidEnginePowerRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidEnginePowerRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidEngineCapacityRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidEngineCapacityRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidEngineCapacityRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(1)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}
func TestGetFiltered_ValidEngineCapacityRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidCarRegistrationDateRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2023-01-02"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidCarRegistrationDateRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2023-01-02"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidCarRegistrationDateRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidCarRegistrationDateRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidOfferCreationDateRange(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	max_ := "2023-01-02"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidOfferCreationDateRangeMinNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2023-01-02"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidOfferCreationDateRangeMaxNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2023-01-01"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

func TestGetFiltered_ValidOfferCreationDateRangeBothNil(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	_, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	u.CleanDB(DB)
}

// ----------------------------------
// Retrieving filtering results tests
// ----------------------------------

func TestGetFiltered_NoFilterEmptyDB(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_NoFilter(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferTypeRegularOffer(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	regularOffer := sale_offer.REGULAR_OFFER
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &regularOffer
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferTypeRegularOfferAuctionInDB(t *testing.T) {
	offers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	regularOffer := sale_offer.REGULAR_OFFER
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &regularOffer
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_OfferTypeAuction(t *testing.T) {
	offers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	auction := sale_offer.AUCTION
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &auction
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferTypeAuctionRegularOfferInDB(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	auction := sale_offer.AUCTION
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &auction
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_OfferTypeBoth(t *testing.T) {
	offers := []models.SaleOffer{*createAuctionSaleOffer(1), *createOffer(2)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	both := sale_offer.BOTH
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferType = &both
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_SingleManufacturer(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(MANUFACTURERS, manufacturer.MapToName)
	filterRequest.Filter.Manufacturers = &[]string{"Audi"}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MultipleManufacturers(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ModelID", uint(1)))), // Audi
		*u.Build(createOffer(2), withCarField(u.WithField[models.Car]("ModelID", uint(2)))), // BMW
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(MANUFACTURERS, manufacturer.MapToName)
	filterRequest.Filter.Manufacturers = &[]string{"Audi", "BMW"}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_NoMatchingManufacturer(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ModelID", uint(1))))} // Audi
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(MANUFACTURERS, manufacturer.MapToName)
	filterRequest.Filter.Manufacturers = &[]string{"BMW"}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_SingleColor(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Color", enums.RED))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Colors = &[]enums.Color{enums.RED}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MultipleColors(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Color", enums.RED))),
		*u.Build(createOffer(2), withCarField(u.WithField[models.Car]("Color", enums.BLUE))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Colors = &[]enums.Color{enums.RED, enums.BLUE}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_NoMatchingColor(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Color", enums.RED))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Colors = &[]enums.Color{enums.GREEN}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_SingleDrive(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Drive", enums.FWD))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Drives = &[]enums.Drive{enums.FWD}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}
func TestGetFiltered_MultipleDrives(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Drive", enums.FWD))),
		*u.Build(createOffer(2), withCarField(u.WithField[models.Car]("Drive", enums.RWD))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Drives = &[]enums.Drive{enums.FWD, enums.RWD}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferNoMatchingDrive(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Drive", enums.FWD))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Drives = &[]enums.Drive{enums.AWD}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_SingleFuelType(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("FuelType", enums.PETROL))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.FuelTypes = &[]enums.FuelType{enums.PETROL}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MultipleFuelTypes(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("FuelType", enums.PETROL))),
		*u.Build(createOffer(2), withCarField(u.WithField[models.Car]("FuelType", enums.DIESEL))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.FuelTypes = &[]enums.FuelType{enums.PETROL, enums.DIESEL}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferNoMatchingFuelType(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("FuelType", enums.PETROL))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.FuelTypes = &[]enums.FuelType{enums.ELECTRIC}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_SingleTransmission(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Transmission", enums.AUTOMATIC))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Transmissions = &[]enums.Transmission{enums.AUTOMATIC}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MultipleTransmissions(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Transmission", enums.AUTOMATIC))),
		*u.Build(createOffer(2), withCarField(u.WithField[models.Car]("Transmission", enums.MANUAL))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Transmissions = &[]enums.Transmission{enums.AUTOMATIC, enums.MANUAL}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferNoMatchingTransmission(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Transmission", enums.AUTOMATIC))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.Transmissions = &[]enums.Transmission{enums.MANUAL}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_PriceInRange(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(100))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_PriceInRangeMinProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(100))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_PriceInRangeMaxProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(100))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_PriceGreater(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(250))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(200)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_PriceLower(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(50))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_PriceUpperBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(100))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_PriceLowerBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Price", uint(100))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.PriceRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MileageInRange(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MileageInRangeMinProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MileageInRangeMaxProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MileageGreater(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(250)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(200)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_MileageLower(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(50)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_MileageUpperBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_MileageLowerBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("Mileage", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.MileageRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_YearInRange(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2024)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2023)
	max_ := uint(2025)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_YearInRangeMinProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2024)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2023)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_YearInRangeMaxProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2024)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2025)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_YearGreater(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2025)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2024)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_YearLower(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2023)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2024)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_YearUpperBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2025)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(2025)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_YearLowerBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("ProductionYear", uint(2025)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(2025)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.YearRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerInRange(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerInRangeMinProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerInRangeMaxProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerGreater(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(250)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(200)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerLower(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(50)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerUpperBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EnginePowerLowerBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EnginePower", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EnginePowerRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityInRange(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityInRangeMinProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(50)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityInRangeMaxProvided(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(150)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityGreater(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(250)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(200)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityLower(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(50)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityUpperBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_EngineCapacityLowerBound(t *testing.T) {
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("EngineCapacity", uint(100)))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := uint(100)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.EngineCapacityRange = &sale_offer.MinMax[uint]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_CarRegistrationDateInRange(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-12"
	max_ := "2025-05-14"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_CarRegistrationDateInRangeMinProvided(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-12"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_CarRegistrationDateInRangeMaxProvided(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2025-05-14"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_CarRegistrationDateGreater(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2025-05-12"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_CarRegistrationDateLower(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-14"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_CarRegistrationDateUpperBound(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2025-05-13"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}
func TestGetFiltered_CarRegistrationDateLowerBound(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), withCarField(u.WithField[models.Car]("RegistrationDate", date))),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-13"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.CarRegistrationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateInRange(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-12"
	max_ := "2025-05-14"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateInRangeMinProvided(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-12"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateInRangeMaxProvided(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2025-05-14"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateGreater(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2025-05-12"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateLower(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-14"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateUpperBound(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	max_ := "2025-05-13"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: nil, Max: &max_}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OfferCreationDateLowerBound(t *testing.T) {
	date, _ := time.Parse(formats.DateLayout, "2025-05-13")
	offers := []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("DateOfIssue", date)),
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	min_ := "2025-05-13"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OfferCreationDateRange = &sale_offer.MinMax[string]{Min: &min_, Max: nil}
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

// ----------------
// Order by tests
// ----------------

func TestGetFiltered_DefaultOrderNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_DefaultOrderSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}
func TestGetFiltered_DefaultOrderMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 3; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), u.WithField[models.SaleOffer]("Margin", enums.Margins[i-1])))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	trueStm := true
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Margin, enums.Margins[len(offers)-i-1])
	}
	u.CleanDB(DB)
}

func TestGetFiltered_DefaultOrderMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 3; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), u.WithField[models.SaleOffer]("Margin", enums.Margins[i-1])))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	falseStm := false
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Margin, enums.Margins[i])
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByPriceNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Price"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByPriceSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Price"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByPriceMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), u.WithField[models.SaleOffer]("Price", uint(i))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Price"
	trueStm := true
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Price, uint(len(offers)-i))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByPriceMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), u.WithField[models.SaleOffer]("Price", uint(i))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Price"
	falseStm := false
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Price, uint(i+1))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByMileageNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Mileage"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByMileageSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Mileage"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByMileageMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("Mileage", uint(i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Mileage"
	trueStm := true
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Mileage, uint(len(offers)-i))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByMileageMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("Mileage", uint(i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Mileage"
	falseStm := false
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].Mileage, uint(i+1))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByYearNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Production year"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByYearSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Production year"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByYearMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("ProductionYear", uint(2000+i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Production year"
	trueStm := true
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].ProductionYear, uint(2000+len(offers)-i))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByYearMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("ProductionYear", uint(2000+i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Production year"
	falseStm := false
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].ProductionYear, uint(2000+i+1))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEnginePowerNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine power"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEnginePowerSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine power"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEnginePowerMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("EnginePower", uint(100+i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine power"
	trueStm := true
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].EnginePower, uint(100+len(offers)-i))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEnginePowerMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("EnginePower", uint(100+i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine power"
	falseStm := false
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].EnginePower, uint(100+i+1))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEngineCapacityNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine capacity"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEngineCapacitySingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine capacity"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByEngineCapacityMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("EngineCapacity", uint(1000+i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine capacity"
	trueStm := true
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].EngineCapacity, uint(1000+len(offers)-i))
	}
	u.CleanDB(DB)
}
func TestGetFiltered_OrderByEngineCapacityMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), withCarField(u.WithField[models.Car]("EngineCapacity", uint(1000+i)))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Engine capacity"
	falseStm := false
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].EngineCapacity, uint(1000+i+1))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByDateOfIssueNoRecords(t *testing.T) {
	var offers []models.SaleOffer
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Date of issue"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByDateOfIssueSingleRecord(t *testing.T) {
	offers := []models.SaleOffer{*createOffer(1)}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Date of issue"
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByDateOfIssueMultipleRecordsDesc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), u.WithField[models.SaleOffer]("DateOfIssue", time.Date(2025, 5, 1+i, 0, 0, 0, 0, time.UTC))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Date of issue"
	trueStm := true
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &trueStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].DateOfIssue, time.Date(2025, 5, 1+(len(offers)-i), 0, 0, 0, 0, time.UTC))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_OrderByDateOfIssueMultipleRecordsAsc(t *testing.T) {
	var offers []models.SaleOffer
	for i := 1; i <= 5; i++ {
		offers = append(offers, *u.Build(createOffer(uint(i)), u.WithField[models.SaleOffer]("DateOfIssue", time.Date(2025, 5, 1+i, 0, 0, 0, 0, time.UTC))))
	}
	db := DB
	repo := getRepositoryWithSaleOffers(db, offers)
	key := "Date of issue"
	falseStm := false
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.Filter.OrderKey = &key
	filterRequest.Filter.IsOrderDesc = &falseStm
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	result, _, err := repo.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
	assert.NoError(t, err)
	for i := range offers {
		assert.Equal(t, result[i].DateOfIssue, time.Date(2025, 5, 1+i+1, 0, 0, 0, 0, time.UTC))
	}
	u.CleanDB(DB)
}
