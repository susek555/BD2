package bid_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	bid "github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
)

// ---------- CREATE ----------

func TestBidService_Create_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	auctionService := new(mocks.AuctionServiceInterface)
	auctionRepo := new(mocks.AuctionRepositoryInterface)
	svc := bid.NewBidService(repo, auctionRepo, auctionService)

	repo.On("Create", mock.MatchedBy(func(b *models.Bid) bool {
		return b.AuctionID == 1 && b.BidderID == 1
	})).Return(nil)

	auctionRepo.On("GetByID", uint(1)).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)

	auctionService.On("GetByIDNonDTO", uint(1)).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)
	auctionService.On("UpdatePrice", mock.AnythingOfType("*models.Auction"), uint(0)).Return(nil)

	dto := &bid.CreateBidDTO{
		AuctionID: 1,
		Amount:    0,
	}
	_, err := svc.Create(dto, 1)
	assert.NoError(t, err)
}
func TestBidService_Create_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	auctionService := new(mocks.AuctionServiceInterface)
	auctionRepo := new(mocks.AuctionRepositoryInterface)
	svc := bid.NewBidService(repo, auctionRepo, auctionService)

	expectedErr := errors.New("insert failed")

	repo.On("Create", mock.MatchedBy(func(b *models.Bid) bool {
		return b.AuctionID == 1
	})).Return(expectedErr)

	auctionRepo.On("GetByID", uint(1)).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)

	auctionService.On("GetByIDNonDTO", uint(1)).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)
	auctionService.On("UpdatePrice", uint(1), uint(0)).Return(nil)

	dto := &bid.CreateBidDTO{
		AuctionID: 1,
		Amount:    0,
	}
	_, err := svc.Create(dto, 1)

	assert.ErrorIs(t, err, expectedErr)
}
func TestBidService_Create_SerializesPerAuction(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	auctionService := new(mocks.AuctionServiceInterface)
	auctionRepo := new(mocks.AuctionRepositoryInterface)
	svc := bid.NewBidService(repo, auctionRepo, auctionService)

	const calls = 2
	const aucID = 777
	var running int32 // how many inside - 0 or 1

	repo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		if atomic.AddInt32(&running, 1) > 1 {
			t.Errorf("mutex did not work – Repo.Create")
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&running, -1)
	}).Return(nil).Times(calls)

	auctionRepo.On("GetByID", uint(aucID)).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)

	auctionService.On("GetByIDNonDTO", uint(aucID), mock.Anything).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)

	auctionService.On("UpdatePrice", mock.AnythingOfType("*models.Auction"), mock.Anything).Run(func(args mock.Arguments) {
		if atomic.AddInt32(&running, 1) > 1 {
			t.Errorf("mutex did not work – AuctionService.UpdatePrice")
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&running, -1)
	}).Return(nil).Times(calls)
	var wg sync.WaitGroup
	wg.Add(calls)

	for i := 0; i < calls; i++ {
		go func() {
			defer wg.Done()
			dto := &bid.CreateBidDTO{
				AuctionID: aucID,
			}
			_, _ = svc.Create(dto, 1)
		}()
	}

	wg.Wait()
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBid_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	auctionService := new(mocks.AuctionServiceInterface)
	auctionRepo := new(mocks.AuctionRepositoryInterface)
	svc := bid.NewBidService(repo, auctionRepo, auctionService)

	expected := &bid.RetrieveBidDTO{
		AuctionID: 10,
		BidderID:  1,
		Amount:    100,
	}

	// Return a models.Bid that the service will convert to RetrieveBidDTO
	modelsBid := &models.Bid{
		AuctionID: 10,
		BidderID:  1,
		Amount:    100,
	}
	repo.On("GetHighestBid", uint(10)).Return(modelsBid, nil)
	auctionService.On("UpdatePrice", uint(1), uint(0)).Return(nil)
	auctionService.On("GetByIDNonDTO", uint(1), mock.Anything).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)

	got, err := svc.GetHighestBid(10)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBid_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	auctionService := new(mocks.AuctionServiceInterface)
	auctionRepo := new(mocks.AuctionRepositoryInterface)
	svc := bid.NewBidService(repo, auctionRepo, auctionService)

	expectedErr := errors.New("db error")

	repo.On("GetHighestBid", uint(10)).Return(nil, expectedErr)
	auctionService.On("UpdatePrice", uint(1), uint(0)).Return(nil)
	auctionService.On("GetByIDNonDTO", uint(1), mock.Anything).Return(&models.Auction{
		OfferID:      1,
		DateEnd:      time.Now().Add(24 * time.Hour),
		BuyNowPrice:  100,
		InitialPrice: 50,
		Offer: &models.SaleOffer{
			Status: enums.PUBLISHED,
		},
	}, nil)

	got, err := svc.GetHighestBid(10)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}
