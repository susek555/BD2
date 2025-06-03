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
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type BidRepo interface {
	GetHighestBid(auctionID uint) (*models.Bid, error)
}

type Scheduler struct {
	mu                  sync.Mutex
	heap                timerHeap
	repo                BidRepo
	notificationService notification.NotificationServiceInterface
	redisClient         *redis.Client

	addCh        chan *Item
	forceCloseCh chan BuyNowRecord

	saleOfferRepository sale_offer.SaleOfferRepositoryInterface
	hub                 ws.HubInterface
}

//go:generate mockery --name=SchedulerInterface --output=../../test/mocks --case=snake --with-expecter
type SchedulerInterface interface {
	AddAuction(auctionID string, end time.Time)
	Run(ctx context.Context)
	LoadAuctions() error
	ForceCloseAuction(auctionID string, buyerID uint, amount uint)
}

func NewScheduler(repo BidRepo, redisClient *redis.Client, notificationService notification.NotificationServiceInterface, saleOfferRepo sale_offer.SaleOfferRepositoryInterface, hub ws.HubInterface) SchedulerInterface {
	return &Scheduler{
		heap:         make(timerHeap, 0),
		addCh:        make(chan *Item, 1024),
		forceCloseCh: make(chan BuyNowRecord, 1024),

		notificationService: notificationService,
		repo:                repo,
		redisClient:         redisClient,
		saleOfferRepository: saleOfferRepo,
		hub:                 hub,
	}
}

func (s *Scheduler) LoadAuctions() error {
	offers, err := s.saleOfferRepository.GetAllActiveAuctions()
	if err != nil {
		log.Println("scheduler: error loading auctions:", err)
		return err
	}
	for _, offer := range offers {
		auctionID := strconv.FormatUint(uint64(offer.ID), 10)
		if offer.Auction.DateEnd.Local().Before(time.Now()) {
			log.Printf("scheduler: skipping auction %s with end time %s, already ended", auctionID, offer.Auction.DateEnd)
			continue
		}
		item := &Item{
			AuctionID: auctionID,
			EndAt:     offer.Auction.DateEnd,
		}
		s.mu.Lock()
		heap.Push(&s.heap, item)
		s.mu.Unlock()
		log.Printf("scheduler: loaded auction %s with end time %s", auctionID, offer.Auction.DateEnd)
	}
	return nil
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
		s.mu.Lock()
		var delay time.Duration
		if len(s.heap) > 0 {
			delay = max(time.Until(s.heap[0].EndAt), 0)
		} else {
			delay = time.Hour * 24 * 365
		}
		s.mu.Unlock()

		if timer == nil {
			timer = time.NewTimer(delay)
		} else {
			timer.Stop()
			timer.Reset(delay)
		}

		select {
		case <-ctx.Done():
			timer.Stop()
			return

		case item := <-s.addCh:
			s.mu.Lock()
			heap.Push(&s.heap, item)
			s.mu.Unlock()
			continue

		case buyNowRecord := <-s.forceCloseCh:
			s.removeFromHeap(buyNowRecord.AuctionID)
			s.closeAuctionByID(buyNowRecord.AuctionID, &buyNowRecord.BuyerID, &buyNowRecord.Price)

		case <-timer.C:
			s.mu.Lock()
			if len(s.heap) == 0 {
				s.mu.Unlock()
				continue
			}
			next := heap.Pop(&s.heap).(*Item)
			s.mu.Unlock()

			s.closeAuctionByID(next.AuctionID, nil, nil)
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
	buyNowRecord := BuyNowRecord{
		AuctionID: auctionID,
		Price:     amount,
		BuyerID:   buyerID,
	}
	select {
	case s.forceCloseCh <- buyNowRecord:
	default:
		log.Printf("scheduler: force close channel is full, skipping auction %s", auctionID)
		return
	}
}

func (s *Scheduler) closeAuctionByID(auctionID string, buyerID *uint, amount *uint) {
	auctionIDInt, err := strconv.Atoi(auctionID)
	if err != nil {
		log.Printf("scheduler: invalid auctionID %q: %v", auctionID, err)
		return
	}
	if buyerID != nil && amount != nil {
		log.Printf("scheduler: force closing auction %d by buyer %d with amount %d", auctionIDInt, *buyerID, *amount)
		if err := s.saleOfferRepository.UpdateStatus(uint(auctionIDInt), enums.SOLD); err != nil {
			log.Printf("scheduler: error updating status for auction %d: %v", auctionIDInt, err)
			return
		}
		if err := s.saleOfferRepository.SaveToPurchases(uint(auctionIDInt), *buyerID, *amount); err != nil {
			log.Printf("scheduler: error saving purchase for auction %d: %v", auctionIDInt, err)
			return
		}
		notif := models.Notification{OfferID: uint(auctionIDInt)}
		offer, err := s.saleOfferRepository.GetByID(uint(auctionIDInt))
		if err != nil {
			log.Printf("scheduler: cannot load offer %d: %v", auctionIDInt, err)
			return
		}
		winnerID := strconv.FormatUint(uint64(*buyerID), 10)
		if err := s.notificationService.CreateEndAuctionNotification(&notif, winnerID, int64(*amount), offer); err != nil {
			log.Printf("scheduler: notif create failed: %v", err)
			_ = s.saleOfferRepository.UpdateStatus(uint(auctionIDInt), enums.EXPIRED)
			return
		}
		s.hub.SaveNotificationForClients(auctionID, 0, &notif)
		s.hub.SendFourLatestNotificationsToClient(auctionID, "0")
		return
	}
	highest, err := s.repo.GetHighestBid(uint(auctionIDInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("scheduler: no bids found for auction %d, marking as expired", auctionIDInt)
			_ = s.saleOfferRepository.UpdateStatus(uint(auctionIDInt), enums.EXPIRED)
			return
		}
		log.Printf("scheduler: error fetching highest bid for auction %d: %v", auctionIDInt, err)
		return
	}
	winnerID := strconv.FormatUint(uint64(highest.BidderID), 10)
	amount = &highest.Amount

	notif := models.Notification{OfferID: uint(auctionIDInt)}
	offer, err := s.saleOfferRepository.GetByID(uint(auctionIDInt))
	if err != nil {
		log.Printf("scheduler: cannot load offer %d: %v", auctionIDInt, err)
		return
	}

	if err := s.notificationService.CreateEndAuctionNotification(&notif, winnerID, int64(*amount), offer); err != nil {
		log.Printf("scheduler: notif create failed: %v", err)
		_ = s.saleOfferRepository.UpdateStatus(uint(auctionIDInt), enums.EXPIRED)
		return
	}

	s.hub.SaveNotificationForClients(auctionID, 0, &notif)
	s.hub.SendFourLatestNotificationsToClient(auctionID, "0")

	_ = s.saleOfferRepository.UpdateStatus(uint(auctionIDInt), enums.SOLD)
	_ = s.saleOfferRepository.SaveToPurchases(uint(auctionIDInt), highest.BidderID, highest.Amount)
}
