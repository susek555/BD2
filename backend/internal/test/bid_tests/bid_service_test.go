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
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	repo.On("Create", mock.MatchedBy(func(b *models.Bid) bool {
		return b.AuctionID == 1 && b.BidderID == 1
	})).Return(nil)

	saleOfferRetriever.On("GetByID", uint(1)).Return(&models.SaleOffer{
		ID:     1,
		Status: enums.PUBLISHED,
		Auction: &models.Auction{
			DateEnd:      time.Now().Add(24 * time.Hour),
			BuyNowPrice:  100,
			InitialPrice: 50,
		}}, nil)

	auctionPriceUpdater.On("UpdatePrice", mock.AnythingOfType("*models.SaleOffer"), uint(0)).Return(nil)

	dto := &bid.CreateBidDTO{
		AuctionID: 1,
		Amount:    0,
	}
	_, err := svc.Create(dto, 1)
	assert.NoError(t, err)
}
func TestBidService_Create_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("insert failed")

	repo.On("Create", mock.MatchedBy(func(b *models.Bid) bool {
		return b.AuctionID == 1
	})).Return(expectedErr)

	saleOfferRetriever.On("GetByID", uint(1)).Return(&models.SaleOffer{
		ID:     1,
		Status: enums.PUBLISHED,
		Auction: &models.Auction{
			DateEnd:      time.Now().Add(24 * time.Hour),
			BuyNowPrice:  100,
			InitialPrice: 50,
		}}, nil)

	auctionPriceUpdater.On("UpdatePrice", uint(1), uint(0)).Return(nil)

	dto := &bid.CreateBidDTO{
		AuctionID: 1,
		Amount:    0,
	}
	_, err := svc.Create(dto, 1)

	assert.ErrorIs(t, err, expectedErr)
}
func TestBidService_Create_SerializesPerAuction(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	const calls = 2
	const aucID = 777
	var running int32 // how many inside - 0 or 1

	repo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		if atomic.AddInt32(&running, 1) > 1 {
			t.Errorf("mutex did not work - Repo.Create")
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&running, -1)
	}).Return(nil).Times(calls)

	saleOfferRetriever.On("GetByID", uint(777)).Return(&models.SaleOffer{
		ID:     777,
		Status: enums.PUBLISHED,
		Auction: &models.Auction{
			DateEnd:      time.Now().Add(24 * time.Hour),
			BuyNowPrice:  100,
			InitialPrice: 50,
		}}, nil)

	auctionPriceUpdater.On("UpdatePrice", mock.AnythingOfType("*models.SaleOffer"), mock.Anything).Run(func(args mock.Arguments) {
		if atomic.AddInt32(&running, 1) > 1 {
			t.Errorf("mutex did not work - AuctionService.UpdatePrice")
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
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

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
	auctionPriceUpdater.On("UpdatePrice", uint(1), uint(0)).Return(nil)

	got, err := svc.GetHighestBid(10)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBid_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("db error")

	repo.On("GetHighestBid", uint(10)).Return(nil, expectedErr)
	auctionPriceUpdater.On("UpdatePrice", uint(1), uint(0)).Return(nil)

	got, err := svc.GetHighestBid(10)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetAll_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	modelsBids := []models.Bid{
		{ID: 1, AuctionID: 10, BidderID: 1, Amount: 100},
		{ID: 2, AuctionID: 11, BidderID: 2, Amount: 200},
	}
	expected := []bid.RetrieveBidDTO{
		{AuctionID: 10, BidderID: 1, Amount: 100},
		{AuctionID: 11, BidderID: 2, Amount: 200},
	}

	repo.On("GetAll").Return(modelsBids, nil)

	got, err := svc.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetAll_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("db error")

	repo.On("GetAll").Return(nil, expectedErr)

	got, err := svc.GetAll()

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetByID_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	modelsBid := &models.Bid{
		ID:        1,
		AuctionID: 10,
		BidderID:  1,
		Amount:    100,
	}

	expected := &bid.RetrieveBidDTO{
		AuctionID: 10,
		BidderID:  1,
		Amount:    100,
	}

	repo.On("GetByID", uint(1)).Return(modelsBid, nil)

	got, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("bid not found")

	repo.On("GetByID", uint(999)).Return(nil, expectedErr)

	got, err := svc.GetByID(999)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetByBidderID_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	modelsBids := []models.Bid{
		{ID: 1, AuctionID: 10, BidderID: 1, Amount: 100},
		{ID: 2, AuctionID: 11, BidderID: 1, Amount: 150},
	}
	expected := []bid.RetrieveBidDTO{
		{AuctionID: 10, BidderID: 1, Amount: 100},
		{AuctionID: 11, BidderID: 1, Amount: 150},
	}

	repo.On("GetByBidderID", uint(1)).Return(modelsBids, nil)

	got, err := svc.GetByBidderID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetByBidderID_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("db error")

	repo.On("GetByBidderID", uint(1)).Return(nil, expectedErr)

	got, err := svc.GetByBidderID(1)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetByAuctionID_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	modelsBids := []models.Bid{
		{ID: 1, AuctionID: 10, BidderID: 1, Amount: 100},
		{ID: 2, AuctionID: 10, BidderID: 2, Amount: 120},
	}
	expected := []bid.RetrieveBidDTO{
		{AuctionID: 10, BidderID: 1, Amount: 100},
		{AuctionID: 10, BidderID: 2, Amount: 120},
	}

	repo.On("GetByAuctionID", uint(10)).Return(modelsBids, nil)

	got, err := svc.GetByAuctionID(10)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetByAuctionID_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("db error")

	repo.On("GetByAuctionID", uint(10)).Return(nil, expectedErr)

	got, err := svc.GetByAuctionID(10)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBidByUserID_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	modelsBid := &models.Bid{
		ID:        1,
		AuctionID: 10,
		BidderID:  1,
		Amount:    100,
	}
	expected := &bid.RetrieveBidDTO{
		AuctionID: 10,
		BidderID:  1,
		Amount:    100,
	}

	repo.On("GetHighestBidByUserID", uint(10), uint(1)).Return(modelsBid, nil)

	got, err := svc.GetHighestBidByUserID(10, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBidByUserID_NotFound(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	saleOfferRetriever := new(mocks.SaleOfferRetrieverInterface)
	auctionPriceUpdater := new(mocks.AuctionPriceUpdaterInterface)
	svc := bid.NewBidService(repo, saleOfferRetriever, auctionPriceUpdater)

	expectedErr := errors.New("no bid found for user")

	repo.On("GetHighestBidByUserID", uint(10), uint(999)).Return(nil, expectedErr)

	got, err := svc.GetHighestBidByUserID(10, 999)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}
