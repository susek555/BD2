package auction_test

import (
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
)

func makeFullAuctionEntity(id uint, end time.Time, buyNow uint) *models.Auction {
	const (
		mockUserID = 99
	)
	man := models.Manufacturer{ID: 2, Name: "Tesla"}
	mod := models.Model{ID: 1, Name: "ModelS", Manufacturer: man}
	car := &sale_offer.Car{
		OfferID:        id,
		Mileage:        1000,
		ProductionYear: 2021,
		Color:          car_params.RED,
		ModelID:        mod.ID,
		Model:          mod,
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
			Margin:             100,
			DateOfIssue:        time.Now(),
			Vin:                "VIN123456789",
			ProductionYear:     2020,
			Mileage:            12345,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   time.Now().Format("2006-01-02"),
			Color:              car_params.RED,
			FuelType:           car_params.DIESEL,
			Transmission:       car_params.MANUAL,
			NumberOfGears:      6,
			Drive:              car_params.FWD,
			ModelID:            7,
		},
		DateEnd:     future.Format("15:04 02/01/2006"),
		BuyNowPrice: 5000,
	}
}

func TestAuctionService_Create_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	dtoIn := makeValidCreateDTO()
	entity, err := dtoIn.MapToAuction()
	assert.NoError(t, err)

	repo.On("Create", entity).Run(func(args mock.Arguments) {
		a := args.Get(0).(*models.Auction)
		full := makeFullAuctionEntity(7, a.DateEnd, a.BuyNowPrice)
		*a = *full
	}).Return(nil)

	out, svcErr := svc.Create(dtoIn)
	assert.NoError(t, svcErr)

	// Now MapToDTO saw a fully‐hydrated Auction
	assert.Equal(t, uint(7), out.ID)
	assert.Equal(t, "alice", out.Username)
	assert.Equal(t, "Tesla ModelS", out.Name)
	assert.Equal(t, dtoIn.DateEnd, out.DateEnd)
	assert.Equal(t, dtoIn.BuyNowPrice, out.BuyNowPrice)

	repo.AssertExpectations(t)
}

func TestAuctionService_Create_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	dtoIn := makeValidCreateDTO()
	entity, _ := dtoIn.MapToAuction()
	expected := errors.New("db failure")

	repo.On("Create", entity).Return(expected)

	out, err := svc.Create(dtoIn)
	assert.Nil(t, out)
	assert.ErrorIs(t, err, expected)
	repo.AssertExpectations(t)
}

func TestAuctionService_GetAll_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	now := time.Now()
	a1 := makeFullAuctionEntity(1, now.Add(time.Hour), 100)
	a2 := makeFullAuctionEntity(2, now.Add(2*time.Hour), 200)

	repo.On("GetAll").Return([]models.Auction{*a1, *a2}, nil)

	all, err := svc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 2)

	assert.Equal(t, uint(1), all[0].ID)
	assert.Equal(t, "alice", all[0].Username)
	assert.Equal(t, uint(100), all[0].BuyNowPrice)

	assert.Equal(t, uint(2), all[1].ID)
	assert.Equal(t, uint(200), all[1].BuyNowPrice)

	repo.AssertExpectations(t)
}

func TestAuctionService_GetAll_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	expected := errors.New("fail")
	repo.On("GetAll").Return(nil, expected)

	all, err := svc.GetAll()
	assert.Nil(t, all)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
}

func TestAuctionService_GetById_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	now := time.Now()
	full := makeFullAuctionEntity(3, now.Add(time.Hour), 300)
	repo.On("GetById", uint(3)).Return(full, nil)

	out, err := svc.GetById(3)
	assert.NoError(t, err)

	assert.Equal(t, uint(3), out.ID)
	assert.Equal(t, "alice", out.Username)
	assert.Equal(t, uint(300), out.BuyNowPrice)

	repo.AssertExpectations(t)
}

func TestAuctionService_GetById_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	expected := errors.New("not found")
	repo.On("GetById", uint(3)).Return(nil, expected)

	out, err := svc.GetById(3)
	assert.Nil(t, out)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
}

func TestAuctionService_Update_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	update := auction.UpdateAuctionDTO{
		Id:               5,
		CreateAuctionDTO: *makeValidCreateDTO(),
	}
	update.DateEnd = time.Now().Add(24 * time.Hour).Format("15:04 02/01/2006")

	entity, err := update.MapToAuction()
	assert.NoError(t, err)

	repo.On("Update", entity).Run(func(args mock.Arguments) {
		a := args.Get(0).(*models.Auction)
		full := makeFullAuctionEntity(5, a.DateEnd, a.BuyNowPrice)
		*a = *full
	}).Return(nil)

	out, svcErr := svc.Update(&update)
	assert.NoError(t, svcErr)

	assert.Equal(t, uint(5), out.ID)
	assert.Equal(t, "alice", out.Username)
	assert.Equal(t, uint(5), out.ID)
	assert.Equal(t, uint(5000), out.Price)
	assert.Equal(t, uint(5000), out.BuyNowPrice)

	repo.AssertExpectations(t)
}

func TestAuctionService_Update_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	update := auction.UpdateAuctionDTO{
		Id:               7,
		CreateAuctionDTO: *makeValidCreateDTO(),
	}
	entity, _ := update.MapToAuction()
	expected := errors.New("update failed")

	repo.On("Update", entity).Return(expected)

	out, svcErr := svc.Update(&update)
	assert.Nil(t, out)
	assert.ErrorIs(t, svcErr, expected)

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	repo.On("GetById", uint(8)).Return(full, nil)
	repo.On("Delete", uint(8)).Return(nil)

	err := svc.Delete(8, full.Offer.UserID)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_Unauthorized(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	repo.On("GetById", uint(8)).Return(full, nil)

	err := svc.Delete(8, full.Offer.UserID+1)
	assert.EqualError(t, err, "you are not the owner of this auction")

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_GetById_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	expected := errors.New("not found")
	repo.On("GetById", uint(8)).Return(nil, expected)

	err := svc.Delete(8, 1)
	assert.ErrorIs(t, err, expected)

	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	full := makeFullAuctionEntity(8, time.Now().Add(time.Hour), 400)
	repo.On("GetById", uint(8)).Return(full, nil)
	repo.On("Delete", uint(8)).Return(errors.New("delete failed"))

	err := svc.Delete(8, full.Offer.UserID)
	assert.Error(t, err)

	repo.AssertExpectations(t)
}
