//go:build unit
// +build unit

package auction_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
)

func TestAuctionService_Create_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	auc := sale_offer.Auction{OfferID: 1}

	repo.On("Create", &auc).Return(nil)

	err := svc.Create(&auc)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuctionService_Create_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	auc := sale_offer.Auction{OfferID: 1}
	expectedErr := errors.New("db failure")

	repo.On("Create", &auc).Return(expectedErr)

	err := svc.Create(&auc)

	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}

func TestAuctionService_Update_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	auc := sale_offer.Auction{OfferID: 2}

	repo.On("Update", &auc).Return(nil)

	err := svc.Update(&auc)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuctionService_Update_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	auc := sale_offer.Auction{OfferID: 2}
	expectedErr := errors.New("update failed")

	repo.On("Update", &auc).Return(expectedErr)

	err := svc.Update(&auc)

	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_OK(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	id := uint(3)
	repo.On("Delete", id).Return(nil)

	err := svc.Delete(id)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuctionService_Delete_Error(t *testing.T) {
	repo := new(mocks.AuctionRepositoryInterface)
	svc := auction.NewAuctionService(repo)

	id := uint(3)
	expectedErr := errors.New("delete failed")

	repo.On("Delete", id).Return(expectedErr)

	err := svc.Delete(id)

	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}
