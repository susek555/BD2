package scheduler

import (
	"log"
	"strconv"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
	"gorm.io/gorm"
)

type CloseReason int

const (
	ReasonTimer CloseReason = iota
	ReasonBuyNow
	ReasonBidOverBuyNow
)

type BidRetrieverInterface interface {
	GetHighestBid(auctionID uint) (*models.Bid, error)
}

type SaleOfferRepositoryInterface interface {
	UpdateStatus(offer *models.SaleOffer, status enums.Status) error
	GetByID(id uint) (*models.SaleOffer, error)
	GetAllActiveAuctions() ([]views.SaleOfferView, error)
}

type PurchaseCreatorInterface interface {
	Create(purchase *models.Purchase) error
}

type CloseCmd struct {
	AuctionID uint
	Reason    CloseReason
	WinnerID  *uint
	Amount    *uint
}

type AuctionCloserInterface interface {
	CloseAuction(cmd CloseCmd)
}

type auctionCloser struct {
	bidRepo             BidRetrieverInterface
	saleRepo            SaleOfferRepositoryInterface
	purchaseCreator     PurchaseCreatorInterface
	notificationService notification.NotificationServiceInterface
	hub                 ws.HubInterface
}

func NewAuctionCloser(bidRepo BidRetrieverInterface, saleRepo SaleOfferRepositoryInterface, purchaseCreator PurchaseCreatorInterface, notificationService notification.NotificationServiceInterface, hub ws.HubInterface) AuctionCloserInterface {
	return &auctionCloser{
		bidRepo:             bidRepo,
		saleRepo:            saleRepo,
		purchaseCreator:     purchaseCreator,
		notificationService: notificationService,
		hub:                 hub,
	}
}

func (c *auctionCloser) CloseAuction(cmd CloseCmd) {
	auctionID := cmd.AuctionID
	log.Printf("closer: closing auction %d (reason %d)", auctionID, cmd.Reason)

	offer, err := c.saleRepo.GetByID(auctionID)
	if err != nil {
		log.Printf("closer: cannot load offer %d: %v", auctionID, err)
		return
	}
	if offer.Status != enums.PUBLISHED {
		log.Printf("closer: auction %d already %s — skip", auctionID, offer.Status)
		return
	}

	var winnerID string
	var amount int64

	if cmd.WinnerID != nil && cmd.Amount != nil {
		winnerID = strconv.FormatUint(uint64(*cmd.WinnerID), 10)
		amount = int64(*cmd.Amount)
	} else {
		highest, err := c.bidRepo.GetHighestBid(auctionID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				_ = c.saleRepo.UpdateStatus(offer, enums.EXPIRED)
				return
			}
			log.Printf("closer: GetHighestBid err: %v", err)
			return
		}
		winnerID = strconv.FormatUint(uint64(highest.BidderID), 10)
		amount = int64(highest.Amount)
	}

	if winnerID != "" {
		winnerIDint, err := strconv.Atoi(winnerID)
		if err != nil {
			log.Printf("closer: invalid winnerID %q: %v", winnerID, err)
			return
		}
		purchaseModel := &models.Purchase{OfferID: auctionID, BuyerID: uint(winnerIDint), FinalPrice: uint(amount), IssueDate: time.Now()}
		_ = c.purchaseCreator.Create(purchaseModel)
	}

	_ = c.saleRepo.UpdateStatus(offer, enums.SOLD)

	n := models.Notification{OfferID: auctionID}
	if err := c.notificationService.CreateEndAuctionNotification(&n, winnerID, amount, offer); err != nil {
		log.Printf("closer: notif err: %v", err)
		_ = c.saleRepo.UpdateStatus(offer, enums.EXPIRED)
		return
	}
	idStr := strconv.FormatUint(uint64(auctionID), 10)
	c.hub.SaveNotificationForClients(idStr, 0, &n)
	c.hub.SendFourLatestNotificationsToClient(idStr, "0")
}
