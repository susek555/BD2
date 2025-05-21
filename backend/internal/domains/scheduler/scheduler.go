package scheduler

import (
	"container/heap"
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auctionws"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
)

type Scheduler struct {
	mu                  sync.Mutex
	heap                timerHeap
	repo                bid.BidRepositoryInterface
	notificationService notification.NotificationServiceInterface
	redisClient         *redis.Client
	addCh               chan *Item
}

type SchedulerInterface interface {
	AddAuction(auctionID string, end time.Time)
	Run(ctx context.Context)
}

func NewScheduler(repo bid.BidRepositoryInterface, redisClient *redis.Client, notificationService notification.NotificationServiceInterface) SchedulerInterface {
	return &Scheduler{
		heap:                make(timerHeap, 0),
		addCh:               make(chan *Item, 1024),
		notificationService: notificationService,
		repo:                repo,
		redisClient:         redisClient,
	}
}

func (s *Scheduler) AddAuction(auctionID string, endAt time.Time) {
	item := &Item{
		AuctionID: auctionID,
		EndAt:     endAt,
	}
	s.addCh <- item
}

func (s *Scheduler) Run(ctx context.Context) {
	var timer *time.Timer
	for ctx.Err() == nil {
		var timerC <-chan time.Time

		s.mu.Lock()
		if len(s.heap) > 0 {
			delay := max(time.Until(s.heap[0].EndAt), 0)
			if timer == nil {
				timer = time.NewTimer(delay)
			} else {
				timer.Stop()
				timer.Reset(delay)
			}
			timerC = timer.C
		}
		select {
		case <-ctx.Done():
			log.Println("scheduler: shutting down")
			if timer != nil {
				timer.Stop()
			}
			return
		case item := <-s.addCh:
			log.Printf("scheduler: adding auction %s", item.AuctionID)
			heap.Push(&s.heap, item)
			s.mu.Unlock()
		case <-timerC:
			log.Printf("scheduler: closing auction %s", s.heap[0].AuctionID)
			next := heap.Pop(&s.heap).(*Item)
			s.mu.Unlock()
			s.closeAuction(ctx, next.AuctionID)
		}
	}
}

func (s *Scheduler) closeAuction(ctx context.Context, auctionID string) {
	auctionIDInt, err := strconv.Atoi(auctionID)
	if err != nil {
		return
	}
	highest, err := s.repo.GetHighestBid(uint(auctionIDInt))
	var winnerID string
	var amount int64
	if err == nil {
		winnerID = strconv.FormatUint(uint64(highest.BidderID), 10)
		amount = int64(highest.Amount)
	}
	notification := notification.Notification{
		OfferID: uint(auctionIDInt),
	}
	err = s.notificationService.CreateEndAuctionNotification(&notification, winnerID, amount)
	if err != nil {
		log.Println("Error creating notification:", err)
		return
	}
	env := auctionws.NewNotificationEnvelope(&notification)
	if err := auctionws.PublishAuctionEvent(ctx, s.redisClient, auctionID, env); err != nil {
		log.Printf("failed to publish auction event: %v", err)
	}
}
