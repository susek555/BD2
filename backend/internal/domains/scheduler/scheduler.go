package scheduler

import (
	"container/heap"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auctionws"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"gorm.io/gorm"
)

type Scheduler struct {
	heap        timerHeap
	repo        bid.BidRepositoryInterface
	redisClient *redis.Client
}

type SchedulerInterface interface {
	AddAuction(auctionID string, end time.Time)
	Run(ctx context.Context)
}

func NewScheduler(db *gorm.DB, repo bid.BidRepositoryInterface, redisClient *redis.Client) *Scheduler {
	return &Scheduler{
		heap:        make(timerHeap, 0),
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *Scheduler) AddAuction(auctionID string, endAt time.Time) {
	item := &Item{
		AuctionID: auctionID,
		EndAt:     endAt,
	}
	heap.Push(&s.heap, item)
}

func (s *Scheduler) Run(ctx context.Context) {
	for ctx.Err() == nil {
		if len(s.heap) == 0 {
			time.Sleep(time.Second)
			continue
		}
		next := s.heap[0]
		delay := time.Until(next.EndAt)
		if delay > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
		}
		heap.Pop(&s.heap)
		// TODO: Move the auction to transaction
		auctionIDuint64, err := strconv.ParseUint(next.AuctionID, 10, 32)
		if err != nil {
			continue
		}
		auctionIDuint := uint(auctionIDuint64)

		highestBid, err := s.repo.GetHighestBid(auctionIDuint)
		if err != nil {
			continue
		}
		env := auctionws.NewEndAuctionEnvelope(next.AuctionID, highestBid.Bidder.Username, int64(highestBid.Amount))
		auctionws.PublishAuctionEvent(ctx, s.redisClient, next.AuctionID, env)
	}
}
