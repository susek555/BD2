package bid_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"

	// ----------------------------------------------------------

	bid "github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
)

// ---------- CREATE ----------

func TestBidService_Create_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	svc := bid.NewBidService(repo)

	b := models.Bid{AuctionID: 1}

	repo.On("Create", &b).Return(nil)

	dto := &bid.CreateBidDTO{
		AuctionID: b.AuctionID,
		Amount:    b.Amount,
		UserID:    1,
	}
	_, err := svc.Create(dto)
}
func TestBidService_Create_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	svc := bid.NewBidService(repo)

	b := models.Bid{AuctionID: 1}
	expectedErr := errors.New("insert failed")

	repo.On("Create", &b).Return(expectedErr)

	dto := &bid.CreateBidDTO{
		AuctionID: b.AuctionID,
		Amount:    b.Amount,
	}
	_, err := svc.Create(dto)

	assert.ErrorIs(t, err, expectedErr)
}
func TestBidService_Create_SerializesPerAuction(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	svc := bid.NewBidService(repo)

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

	var wg sync.WaitGroup
	wg.Add(calls)

	for i := 0; i < calls; i++ {
		go func() {
			defer wg.Done()
			dto := &bid.CreateBidDTO{
				AuctionID: aucID,
			}
			_, _ = svc.Create(dto)
		}()
	}

	wg.Wait()
	repo.AssertExpectations(t)

		go func() {
			defer wg.Done()
			_ = svc.Create(&models.Bid{AuctionID: aucID})
		}()

	wg.Wait()
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBid_OK(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	svc := bid.NewBidService(repo)

	expected := &models.Bid{AuctionID: 10, Amount: 1000}

	repo.On("GetHighestBid", uint(10)).Return(expected, nil)

	got, err := svc.GetHighestBid(10)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestBidService_GetHighestBid_Error(t *testing.T) {
	repo := new(mocks.BidRepositoryInterface)
	svc := bid.NewBidService(repo)

	expectedErr := errors.New("db error")

	repo.On("GetHighestBid", uint(10)).Return(nil, expectedErr)

	got, err := svc.GetHighestBid(10)

	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, got)
	repo.AssertExpectations(t)
}
