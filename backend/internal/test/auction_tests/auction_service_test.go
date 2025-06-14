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
		mockUserID = 5
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
		BuyNowPrice: &buyNow,
		Offer:       so,
	}
	user := &models.User{
		ID:       mockUserID,
		Username: "alice",
	}
	so.Auction = auc
	so.User = user
	return auc
}

func makeValidCreateDTO() *auction.CreateAuctionDTO {
	future := time.Date(2100, 1, 1, 12, 0, 0, 0, time.UTC)
	buyNowPrice := uint(5000)
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
		BuyNowPrice: &buyNowPrice,
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
		UpdateSaleOfferDTO: sale_offer.UpdateSaleOfferDTO{
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
			ManufacturerName:   &manufacturerName,
			ModelName:          &modelName,
		},
		DateEnd:     &dateEnd,
		BuyNowPrice: &buyNowPrice,
	}
}

func TestAuctionService_Create_OK(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)

	dtoIn := makeValidCreateDTO()

	// The mock expects two arguments for PrepareForCreateSaleOffer
	saleOfferSvc.On("PrepareForCreateSaleOffer", mock.AnythingOfType("*sale_offer.CreateSaleOfferDTO")).
		Return(&models.SaleOffer{
			ID:     7,
			UserID: dtoIn.UserID}, nil)

	repo.On("Create", mock.AnythingOfType("*models.SaleOffer")).Return(nil)

	saleOfferSvc.On("GetDetailedByID", uint(7), mock.AnythingOfType("*uint")).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 7, Price: dtoIn.Price, BuyNowPrice: dtoIn.BuyNowPrice}, nil)

	out, err := svc.Create(dtoIn)
	assert.NoError(t, err)

	assert.Equal(t, uint(7), out.ID)
	assert.Equal(t, *dtoIn.BuyNowPrice, *out.BuyNowPrice)

	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}

func TestAuctionService_Create_Error(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)

	dtoIn := makeValidCreateDTO()
	expected := errors.New("db failure")

	saleOfferSvc.On("PrepareForCreateSaleOffer", mock.AnythingOfType("*sale_offer.CreateSaleOfferDTO")).
		Return(&models.SaleOffer{
			ID:     7,
			UserID: dtoIn.UserID}, expected)

	out, err := svc.Create(dtoIn)
	assert.Nil(t, out)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}
func TestAuctionService_Update_OK(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)
	update := makeValidUpdateDTO()

	// Mock GetByID which is called before Update
	buyNowPrice := uint(1000)
	saleOfferSvc.On("PrepareForUpdateSaleOffer", mock.AnythingOfType("*sale_offer.UpdateSaleOfferDTO"), mock.AnythingOfType("uint")).
		Return(&models.SaleOffer{
			ID:     7,
			UserID: update.ID,
			Auction: &models.Auction{
				OfferID:      7,
				DateEnd:      time.Now().Add(time.Hour),
				BuyNowPrice:  &buyNowPrice,
				InitialPrice: 500,
			}}, nil)

	repo.On("Update", mock.AnythingOfType("*models.SaleOffer")).
		Run(func(args mock.Arguments) {
		}).
		Return(nil)

	buyNowPrice = uint(1000)
	saleOfferSvc.On("GetDetailedByID", uint(7), mock.AnythingOfType("*uint")).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 7, Price: 500, BuyNowPrice: &buyNowPrice}, nil)

	out, svcErr := svc.Update(update, uint(99))
	assert.NoError(t, svcErr)

	assert.Equal(t, uint(7), out.ID)
	repo.AssertExpectations(t)
}
func TestAuctionService_Update_Error(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)
	update := makeValidUpdateDTO()

	expected := errors.New("db failure")

	// Mock GetByID which is called before Update
	buyNowPrice := uint(1000)
	saleOfferSvc.On("PrepareForUpdateSaleOffer", mock.AnythingOfType("*sale_offer.UpdateSaleOfferDTO"), mock.AnythingOfType("uint")).
		Return(&models.SaleOffer{
			ID:     7,
			UserID: update.ID,
			Auction: &models.Auction{
				OfferID:      7,
				DateEnd:      time.Now().Add(time.Hour),
				BuyNowPrice:  &buyNowPrice,
				InitialPrice: 500,
			}}, nil)

	repo.On("Update", mock.AnythingOfType("*models.SaleOffer")).
		Run(func(args mock.Arguments) {
		}).
		Return(nil)

	buyNowPrice = uint(1000)
	saleOfferSvc.On("GetDetailedByID", uint(7), mock.AnythingOfType("*uint")).
		Return(&sale_offer.RetrieveDetailedSaleOfferDTO{ID: 7, Price: 500, BuyNowPrice: &buyNowPrice}, expected)

	out, svcErr := svc.Update(update, uint(99))
	assert.ErrorIs(t, svcErr, expected)

	assert.Equal(t, uint(7), out.ID)
	repo.AssertExpectations(t)
}
func TestAuctionService_Delete_OK(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	saleOfferSvc.On("Delete", uint(8), full.Offer.UserID).Return(nil)
	err := svc.Delete(8, full.Offer.UserID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}

func TestAuctionService_Delete_Unauthorized(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)

	saleOfferSvc.On("Delete", uint(8), full.Offer.UserID+1).Return(errors.New("you are not the owner of this auction"))

	err := svc.Delete(8, full.Offer.UserID+1)
	assert.EqualError(t, err, "you are not the owner of this auction")

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_GetByID_Error(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)

	expected := errors.New("auction not found")

	saleOfferSvc.On("Delete", uint(8), uint(1)).Return(expected)

	err := svc.Delete(8, 1)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}

func TestAuctionService_Delete_Error(t *testing.T) {
	repo := new(mocks.SaleOfferRepositoryInterface)
	saleOfferSvc := new(mocks.SaleOfferServiceInterface)
	purchaseCreator := new(mocks.PurchaseCreatorInterface)
	svc := auction.NewAuctionService(repo, saleOfferSvc, purchaseCreator)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	saleOfferSvc.On("Delete", uint(8), full.Offer.UserID).Return(errors.New("delete failed"))

	err := svc.Delete(8, full.Offer.UserID)
	assert.Error(t, err)

	repo.AssertExpectations(t)
	saleOfferSvc.AssertExpectations(t)
}
