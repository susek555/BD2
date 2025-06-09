package scheduler

import (
	"container/heap"
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
)

type Scheduler struct {
	mu            sync.Mutex
	heap          timerHeap
	eventsCh      chan AuctionEvent
	closer        AuctionCloserInterface
	saleOfferRepo SaleOfferRepositoryInterface
}

//go:generate mockery --name=SchedulerInterface --output=../../test/mocks --case=snake --with-expecter
type SchedulerInterface interface {
	AddAuction(auctionID string, end time.Time)
	Run(ctx context.Context)
	LoadAuctions() error
	ForceCloseAuction(auctionID string, buyerID uint, amount uint)
}

func NewScheduler(
	repo BidRetrieverInterface,
	redisClient *redis.Client,
	notificationService notification.NotificationServiceInterface,
	saleOfferRepo SaleOfferRepositoryInterface,
	purchaseCreator PurchaseCreatorInterface,
	saleOfferService SaleOfferRetrieverInterface,
	hub ws.HubInterface,
) SchedulerInterface {
	closer := NewAuctionCloser(repo, saleOfferRepo, purchaseCreator, notificationService, hub, saleOfferService)
	return &Scheduler{
		heap:          make(timerHeap, 0),
		eventsCh:      make(chan AuctionEvent, 1024),
		closer:        closer,
		saleOfferRepo: saleOfferRepo,
	}
}

func (s *Scheduler) LoadAuctions() error {
	offers, err := s.saleOfferRepo.GetAllActiveAuctions()
	if err != nil {
		log.Println("scheduler: error loading auctions:", err)
		return err
	}
	for _, offer := range offers {
		auctionID := strconv.FormatUint(uint64(offer.ID), 10)
		if offer.DateEnd.Local().Before(time.Now()) {
			log.Printf("scheduler: skipping auction %s with end time %s, already ended", auctionID, offer.DateEnd)
			continue
		}
		item := &Item{
			AuctionID: auctionID,
			EndAt:     *offer.DateEnd,
		}
		s.mu.Lock()
		heap.Push(&s.heap, item)
		s.mu.Unlock()
		log.Printf("scheduler: loaded auction %s with end time %s", auctionID, offer.DateEnd)
	}
	return nil
}

func (s *Scheduler) AddAuction(auctionID string, endAt time.Time) {
	id, _ := strconv.Atoi(auctionID)
	s.eventsCh <- AuctionEvent{
		Kind: EventAddTimer,
		At:   endAt,
		Cmd: CloseCmd{
			AuctionID: uint(id),
			Reason:    ReasonTimer,
		},
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	var timer *time.Timer = time.NewTimer(time.Hour * 24 * 365)

	for {
		s.mu.Lock()
		if len(s.heap) > 0 {
			timer.Reset(max(time.Until(s.heap[0].EndAt), 0))
		} else {
			timer.Reset(time.Hour * 24 * 365)
		}
		s.mu.Unlock()

		select {
		case <-ctx.Done():
			timer.Stop()
			return

		case ev := <-s.eventsCh:
			switch ev.Kind {
			case EventAddTimer:
				s.mu.Lock()
				heap.Push(&s.heap, &Item{
					AuctionID: strconv.Itoa(int(ev.Cmd.AuctionID)),
					EndAt:     ev.At,
				})
				s.mu.Unlock()

			case EventForceClose:
				s.removeFromHeap(strconv.Itoa(int(ev.Cmd.AuctionID)))
				s.closer.CloseAuction(ev.Cmd)
			}

		case <-timer.C:
			s.mu.Lock()
			if len(s.heap) == 0 {
				s.mu.Unlock()
				continue
			}
			next := heap.Pop(&s.heap).(*Item)
			s.mu.Unlock()

			id, _ := strconv.Atoi(next.AuctionID)
			s.closer.CloseAuction(CloseCmd{
				AuctionID: uint(id),
				Reason:    ReasonTimer,
			})
		}
	}
}

func (s *Scheduler) removeFromHeap(auctionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, item := range s.heap {
		if item.AuctionID == auctionID {
			heap.Remove(&s.heap, i)
			log.Printf("scheduler: removed auction %s from heap", auctionID)
			return
		}
	}
	log.Printf("scheduler: auction %s not found in heap", auctionID)
}

func (s *Scheduler) ForceCloseAuction(auctionID string, buyerID uint, amount uint) {
	auctionIDInt, _ := strconv.Atoi(auctionID)
	s.eventsCh <- AuctionEvent{
		Kind: EventForceClose,
		Cmd: CloseCmd{
			AuctionID: uint(auctionIDInt),
			Reason:    ReasonBuyNow,
			WinnerID:  &buyerID,
			Amount:    &amount,
		},
	}
	log.Printf("scheduler: force closing auction %s by buyer %d with amount %d", auctionID, buyerID, amount)
}
