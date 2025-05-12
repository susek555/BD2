//go:build integration
// +build integration

package sale_offer_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(offers []sale_offer.SaleOffer) (sale_offer.SaleOfferRepositoryInterface, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&car.Car{},
		&manufacturer.Manufacturer{},
		&model.Model{},
		&sale_offer.SaleOffer{},
		&sale_offer.Auction{},
	)
	if err != nil {
		return nil, err
	}
	repo := sale_offer.NewSaleOfferRepository(db)
	for _, offer := range offers {
		repo.Create(&offer)
	}
	return repo, nil
}

//------------------------
// Invalid arguments tests
// -----------------------

// TODO handle invalid manufacturer

func TestGetFiltered_InvalidOfferType(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType("invalid")
	filter := sale_offer.OfferFilter{OfferType: &offer}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidSaleOfferType)
}

func TestGetFiltered_InvalidColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Colors: &[]car_params.Color{"invaid"}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidColor)
}

func TestGetFiltered_InvalidDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Drives: &[]car_params.Drive{"invaid"}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDrive)
}

func TestGetFiltered_InvalidFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{FuelTypes: &[]car_params.FuelType{"invaid"}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidFuelType)
}

func TestGetFiltered_InvalidTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Transmissions: &[]car_params.Transmission{"invaid"}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidTransmission)
}

func TestGetFiltered_InvalidPriceRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidPriceRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidMileageRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidMileageRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidYearRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidYearRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEnginePowerRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEnginePowerRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEngineCapacityRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(2)
	max := uint(1)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidEngineCapacityRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(1)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidCarRegistrationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022-01-01"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidCarRegistrationDateRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-01"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidCarRegistrationDateFormat(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022/01/01"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDateFromat)
}

func TestGetFiltered_InvalidOfferCreationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022-01-01"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidOfferCreationDateRangeBothValues(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-01"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidRange)
}

func TestGetFiltered_InvalidOfferCreationDateFormat(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2022/01/01"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.ErrorIs(t, err, sale_offer.ErrInvalidDateFromat)
}

// ---------------------
// Valid arguments tests
// ---------------------

func TestGetFiltered_ValidOfferTypeAutcion(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType(sale_offer.AUCTION)
	filter := sale_offer.OfferFilter{OfferType: &offer}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferTypeRegularOffer(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType(sale_offer.REGULAR_OFFER)
	filter := sale_offer.OfferFilter{OfferType: &offer}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidOfferTypeBoth(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	offer := sale_offer.OfferType(sale_offer.BOTH)
	filter := sale_offer.OfferFilter{OfferType: &offer}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Colors: &car_params.Colors}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Drives: &car_params.Drives}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{FuelTypes: &car_params.Types}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Transmissions: &car_params.Transmissions}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidPriceRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidMileageRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidYearRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidEnginePowerRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEnginePowerRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidEnginePowerRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEnginePowerRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	max := uint(2)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := uint(2)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := uint(1)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}
func TestGetFiltered_ValidEngineCapacityRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-02"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := "2023-01-02"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidCarRegistrationDateRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	max := "2023-01-02"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRangeMinNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	max := "2023-01-02"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRangeMaxNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	min := "2023-01-01"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

func TestGetFiltered_ValidOfferCreationDateRangeBothNil(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: nil, Max: nil}}
	_, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
}

// ----------------------------------
// Retrieving filtering results tests
// ----------------------------------

func TestGetFiltered_NoFilterEmptyDB(t *testing.T) {
	offers := []sale_offer.SaleOffer{}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_NoFilter(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1)}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferTypeRegularOffer(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1)}
	repo, _ := setupDB(offers)
	regularOffer := sale_offer.REGULAR_OFFER
	filter := sale_offer.OfferFilter{OfferType: &regularOffer}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferTypeRegularOfferAuctionInDB(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withAuction(time.Now(), 0))}
	repo, _ := setupDB(offers)
	regularOffer := sale_offer.REGULAR_OFFER
	filter := sale_offer.OfferFilter{OfferType: &regularOffer}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferTypeAuction(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withAuction(time.Now(), 0))}
	repo, _ := setupDB(offers)
	auction := sale_offer.AUCTION
	filter := sale_offer.OfferFilter{OfferType: &auction}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferTypeAuctionRegularOfferInDB(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1)}
	repo, _ := setupDB(offers)
	auction := sale_offer.AUCTION
	filter := sale_offer.OfferFilter{OfferType: &auction}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferTypeBoth(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withAuction(time.Now(), 0)), *createOffer(2)}
	repo, _ := setupDB(offers)
	both := sale_offer.BOTH
	filter := sale_offer.OfferFilter{OfferType: &both}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_SingleColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Color", car_params.RED))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Colors: &[]car_params.Color{car_params.RED}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleColors(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Color", car_params.RED)), *createOffer(2, withCarField("Color", car_params.BLUE))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Colors: &[]car_params.Color{car_params.RED, car_params.BLUE}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_NoMatchingColor(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Color", car_params.RED))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Colors: &[]car_params.Color{car_params.GREEN}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Drive", car_params.FWD))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Drives: &[]car_params.Drive{car_params.FWD}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}
func TestGetFiltered_MultipleDrives(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Drive", car_params.FWD)), *createOffer(2, withCarField("Drive", car_params.RWD))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Drives: &[]car_params.Drive{car_params.FWD, car_params.RWD}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferNoMatchingDrive(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Drive", car_params.FWD))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Drives: &[]car_params.Drive{car_params.AWD}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("FuelType", car_params.PETROL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{FuelTypes: &[]car_params.FuelType{car_params.PETROL}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleFuelTypes(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("FuelType", car_params.PETROL)), *createOffer(2, withCarField("FuelType", car_params.DIESEL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{FuelTypes: &[]car_params.FuelType{car_params.PETROL, car_params.DIESEL}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferNoMatchingFuelType(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("FuelType", car_params.PETROL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{FuelTypes: &[]car_params.FuelType{car_params.ELECTRIC}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_SingleTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Transmission", car_params.AUTOMATIC))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Transmissions: &[]car_params.Transmission{car_params.AUTOMATIC}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MultipleTransmissions(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Transmission", car_params.AUTOMATIC)), *createOffer(2, withCarField("Transmission", car_params.MANUAL))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Transmissions: &[]car_params.Transmission{car_params.AUTOMATIC, car_params.MANUAL}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferNoMatchingTransmission(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Transmission", car_params.AUTOMATIC))}
	repo, _ := setupDB(offers)
	filter := sale_offer.OfferFilter{Transmissions: &[]car_params.Transmission{car_params.MANUAL}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_PriceInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_PriceLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetfiltered_PriceUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_PriceLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("Price", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{PriceRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_MileageLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_MileageUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_MileageLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("Mileage", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{MileageRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	min := uint(2025)
	max := uint(2026)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2024)))}
	repo, _ := setupDB(offers)
	min := uint(2023)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2024)))}
	repo, _ := setupDB(offers)
	max := uint(2025)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	max := uint(2024)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_YearLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2023)))}
	repo, _ := setupDB(offers)
	min := uint(2024)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_YearUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	max := uint(2025)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_YearLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("ProductionYear", uint(2025)))}
	repo, _ := setupDB(offers)
	min := uint(2025)
	filter := sale_offer.OfferFilter{YearRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EnginePowerLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EnginePowerUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EnginePowerLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EnginePower", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{EnginePowerRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityInRange(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	max := uint(150)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityInRangeMinProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(50)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityInRangeMaxProvided(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(150)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityGreater(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(250)))}
	repo, _ := setupDB(offers)
	max := uint(200)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EngineCapacityLower(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(50)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_EngineCapacityUpperBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	max := uint(100)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_EngineCapacityLowerBound(t *testing.T) {
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("EngineCapacity", uint(100)))}
	repo, _ := setupDB(offers)
	min := uint(100)
	filter := sale_offer.OfferFilter{EngineCapacityRange: &sale_offer.MinMax[uint]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateInRange(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	max := "2025-05-14"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateInRangeMinProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateInRangeMaxProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-14"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_CarRegistrationDateGreater(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-12"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_CarRegistrationDateLower(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-14"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_CarRegistratoinDateUpperBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-13"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}
func TestGetFiltered_CarRegistrationDateLowerBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withCarField("RegistrationDate", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-13"
	filter := sale_offer.OfferFilter{CarRegistrationDateRagne: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateInRange(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	max := "2025-05-14"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateInRangeMinProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-12"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateInRangeMaxProvided(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-14"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateGreater(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-12"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferCreationDateLower(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-14"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), 0)
}

func TestGetFiltered_OfferCreationDateUpperBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	max := "2025-05-13"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: nil, Max: &max}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}

func TestGetFiltered_OfferCreationDateLowerBound(t *testing.T) {
	date, _ := time.Parse(sale_offer.LAYOUT, "2025-05-13")
	offers := []sale_offer.SaleOffer{*createOffer(1, withOfferField("DateOfIssue", date))}
	repo, _ := setupDB(offers)
	min := "2025-05-13"
	filter := sale_offer.OfferFilter{OfferCreationDateRange: &sale_offer.MinMax[string]{Min: &min, Max: nil}}
	result, _, err := repo.GetFiltered(&filter)
	assert.NoError(t, err)
	assert.Equal(t, len(result), len(offers))
}
