package scheduler

import (
	"log"
	"strconv"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type CloseReason int

const (
	ReasonTimer CloseReason = iota
	ReasonBuyNow
	ReasonBidOverBuyNow
)

type CloseCmd struct {
	AuctionID uint
	Reason    CloseReason
	WinnerID  *uint
	Amount    *uint
}

type AuctionCloser interface {
	CloseAuction(cmd CloseCmd)
}

type auctionCloser struct {
	bidRepo             BidRepo
	saleRepo            sale_offer.SaleOfferRepositoryInterface
	notificationService notification.NotificationServiceInterface
	hub                 ws.HubInterface
}

func NewAuctionCloser(bidRepo BidRepo, saleRepo sale_offer.SaleOfferRepositoryInterface, notificationService notification.NotificationServiceInterface, hub ws.HubInterface) AuctionCloser {
	return &auctionCloser{
		bidRepo:             bidRepo,
		saleRepo:            saleRepo,
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
				_ = c.saleRepo.UpdateStatus(auctionID, enums.EXPIRED)
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
		_ = c.saleRepo.SaveToPurchases(auctionID, uint(winnerIDint), uint(amount))
	}

	_ = c.saleRepo.UpdateStatus(auctionID, enums.SOLD)

	n := models.Notification{OfferID: auctionID}
	if err := c.notificationService.CreateEndAuctionNotification(&n, winnerID, amount, offer); err != nil {
		log.Printf("closer: notif err: %v", err)
		_ = c.saleRepo.UpdateStatus(auctionID, enums.EXPIRED)
		return
	}
	idStr := strconv.FormatUint(uint64(auctionID), 10)
	c.hub.SaveNotificationForClients(idStr, 0, &n)
	c.hub.SendFourLatestNotificationsToClient(idStr, "0")
}
