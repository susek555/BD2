package auction_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

func makeFullAuctionEntity(id uint, end time.Time, buyNow uint) *models.Auction {
	const (
		mockUserID = 99
	)
	man := models.Manufacturer{ID: 2, Name: "Tesla"}
	mod := models.Model{ID: 1, Name: "ModelS", Manufacturer: &man}
	car := &models.Car{
		OfferID:        id,
		Mileage:        1000,
		ProductionYear: 2021,
		Color:          enums.RED,
		ModelID:        mod.ID,
		Model:          &mod,
	}
	so := &models.SaleOffer{
		ID:     id,
		UserID: mockUserID,
		User:   &models.User{ID: mockUserID, Username: "alice"},
		Price:  5000,
		Car:    car,
	}
	auc := &models.Auction{
		OfferID:     id,
		DateEnd:     end,
		BuyNowPrice: buyNow,
		Offer:       so,
	}
	so.Auction = auc
	return auc
}

func makeValidCreateDTO() *auction.CreateAuctionDTO {
	future := time.Date(2100, 1, 1, 12, 0, 0, 0, time.UTC)
	return &auction.CreateAuctionDTO{
		CreateSaleOfferDTO: sale_offer.CreateSaleOfferDTO{
			UserID:             42,
			Description:        "desc",
			Price:              1000,
			Margin:             10,
			Vin:                "VIN123456789",
			ProductionYear:     2020,
			Mileage:            12345,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   time.Now().Format(formats.DateLayout),
			Color:              enums.RED,
			FuelType:           enums.DIESEL,
			Transmission:       enums.MANUAL,
			NumberOfGears:      6,
			Drive:              enums.FWD,
			ManufacturerName:   "Tesla",
			ModelName:          "ModelS",
		},
		DateEnd:     future.Format(formats.DateTimeLayout),
		BuyNowPrice: 5000,
	}
}

func makeValidUpdateDTO() *auction.UpdateAuctionDTO {
	future := time.Date(2100, 1, 1, 12, 0, 0, 0, time.UTC)
	description := "updated desc"
	price := uint(1200)
	margin := enums.HIGH_MARGIN
	vin := "VIN987654321"
	var productionYear uint = 2021
	var mileage uint = 15000
	var numberOfDoors uint = 4
	var numberOfSeats uint = 5
	var enginePower uint = 160
	var engineCapacity uint = 2200
	registrationNumber := "XYZ789"
	registrationDate := time.Now().Format(formats.DateLayout)
	color := enums.BLUE
	fuelType := enums.PETROL
	transmission := enums.AUTOMATIC
	var numberOfGears uint = 7
	drive := enums.RWD
	manufacturerName := "Tesla"
	modelName := "ModelX"
	dateEnd := future.Format(formats.DateTimeLayout)
	buyNowPrice := uint(6000)
	return &auction.UpdateAuctionDTO{
		ID: 42,
		UpdateSaleOfferDTO: &sale_offer.UpdateSaleOfferDTO{
			ID:                 42,
			Description:        &description,
			Price:              &price,
			Margin:             &margin,
			Vin:                &vin,
			ProductionYear:     &productionYear,
			Mileage:            &mileage,
			NumberOfDoors:      &numberOfDoors,
			NumberOfSeats:      &numberOfSeats,
			EnginePower:        &enginePower,
			EngineCapacity:     &engineCapacity,
			RegistrationNumber: &registrationNumber,
			RegistrationDate:   &registrationDate,
			Color:              &color,
			FuelType:           &fuelType,
			Transmission:       &transmission,
			NumberOfGears:      &numberOfGears,
			Drive:              &drive,
			Manufacturer:       &manufacturerName,
			Model:              &modelName,
		},
		DateEnd:     &dateEnd,
		BuyNowPrice: &buyNowPrice,
	}
}

func TestAuctionService_Create_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)

	dtoIn := makeValidCreateDTO()

	// Mock the GetModelID call
	saleOfferSvc.On("Create", mock.AnythingOfType("*sale_offer.CreateSaleOfferDTO")).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 7, Price: dtoIn.Price}, nil)
	repo.On("Create", mock.AnythingOfType("*models.Auction")).Run(func(args mock.Arguments) {
		a := args.Get(0).(*models.Auction)
		full := makeFullAuctionEntity(7, a.DateEnd, a.BuyNowPrice)
		*a = *full
	}).Return(nil)

	// Mock the GetByID call that happens after Create
	repo.On("GetByID", uint(7)).Return(makeFullAuctionEntity(7, time.Now().Add(time.Hour), dtoIn.BuyNowPrice), nil)

	out, err := svc.Create(dtoIn)
	assert.NoError(t, err)

	assert.Equal(t, uint(7), out.ID)
	assert.Equal(t, dtoIn.BuyNowPrice, out.BuyNowPrice)

	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}

func TestAuctionService_Create_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)

	dtoIn := makeValidCreateDTO()
	expected := errors.New("db failure")

	// Mock the GetModelID call that happens before Create
	saleOfferSvc.On("Create", mock.AnythingOfType("*sale_offer.CreateSaleOfferDTO")).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 1, Price: dtoIn.Price}, nil)
	repo.On("Create", mock.Anything).Return(expected)

	out, err := svc.Create(dtoIn)
	assert.Nil(t, out)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}
func TestAuctionService_Update_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)
	update := makeValidUpdateDTO()

	// Mock GetByID which is called before Update
	repo.On("GetByID", uint(42)).Return(makeFullAuctionEntity(42, time.Now().Add(time.Hour), 300), nil)
	saleOfferSvc.On("Update", mock.AnythingOfType("*sale_offer.UpdateSaleOfferDTO"), uint(99)).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 99}, nil)
	repo.On("Update", mock.AnythingOfType("*models.Auction")).Run(func(args mock.Arguments) {
		a := args.Get(0).(*models.Auction)
		full := makeFullAuctionEntity(5, a.DateEnd, a.BuyNowPrice)
		*a = *full
	}).Return(nil)

	// Mock the GetByID call that happens after Update
	repo.On("GetByID", uint(5)).Return(makeFullAuctionEntity(5, time.Now().Add(time.Hour), 6000), nil)

	out, svcErr := svc.Update(update, uint(99))
	assert.NoError(t, svcErr)

	assert.Equal(t, uint(5), out.ID)
	assert.Equal(t, "alice", out.Username)
	repo.AssertExpectations(t)
}
func TestAuctionService_Update_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)

	update := makeValidUpdateDTO()
	expected := errors.New("update failed")

	// Mock GetByID which is called before Update
	repo.On("GetByID", uint(42)).Return(makeFullAuctionEntity(42, time.Now().Add(time.Hour), 300), nil)

	saleOfferSvc.On("Update", mock.AnythingOfType("*sale_offer.UpdateSaleOfferDTO"), uint(99)).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 99}, nil)

	repo.On("Update", mock.Anything).Return(expected)

	_, err := svc.Update(update, uint(99))
	assert.ErrorIs(t, err, expected)
	repo.AssertExpectations(t)
}
func TestAuctionService_Delete_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	repo.On("GetByID", uint(8)).Return(full, nil)
	repo.On("Delete", uint(8)).Return(nil)
	err := svc.Delete(8, full.Offer.UserID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_Unauthorized(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	repo.On("GetByID", uint(8)).Return(full, nil)
	svc := auction.NewAuctionService(repo, saleOfferSvc)

	err := svc.Delete(8, full.Offer.UserID+1)
	assert.EqualError(t, err, "you are not the owner of this auction")

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_GetByID_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)
	expected := errors.New("auction not found")
	repo.On("GetByID", uint(8)).Return(nil, expected)

	err := svc.Delete(8, 1)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	repo.On("GetByID", uint(8)).Return(full, nil)
	repo.On("Delete", uint(8)).Return(errors.New("delete failed"))

	err := svc.Delete(8, full.Offer.UserID)
	assert.Error(t, err)

	repo.AssertExpectations(t)
}
