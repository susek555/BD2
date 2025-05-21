package notification

import (
	"fmt"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
)

type NotificationServiceInterface interface {
	CreateOutbidNotification(notification *Notification, amount int64) error
	CreateEndAuctionNotification(notification *Notification, winner string, winningBid int64) error
	GetNotificationByID(id uint) (*Notification, error)
}

type NotificationService struct {
	NotificationRepository NotificationRepositoryInterface
	SaleOfferService       sale_offer.SaleOfferServiceInterface
}

func NewNotificationService(notificationRepository NotificationRepositoryInterface, saleOfferService sale_offer.SaleOfferServiceInterface) NotificationServiceInterface {
	return &NotificationService{
		NotificationRepository: notificationRepository,
		SaleOfferService:       saleOfferService,
	}
}

func (s *NotificationService) CreateOutbidNotification(notification *Notification, amount int64) error {
	offer, err := s.SaleOfferService.GetByID(notification.OfferID)
	if err != nil {
		return err
	}
	notification.Date = time.Now().Format("2006-01-02 15:04:05")
	notification.Title = fmt.Sprintf("Someone outbid you on %s %s \n", offer.Brand, offer.Model)
	notification.Description = fmt.Sprintf("You were outbid on your offer for %s %s. \n New price: %d \n", offer.Brand, offer.Model, amount)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateEndAuctionNotification(notification *Notification, winner string, winningBid int64) error {
	offer, err := s.SaleOfferService.GetByID(notification.OfferID)
	if err != nil {
		return err
	}
	notification.Date = time.Now().Format("2006-01-02 15:04:05")
	notification.Title = fmt.Sprintf("Auction ended for %s %s \n", offer.Brand, offer.Model)
	notification.Description = fmt.Sprintf("The auction for %s %s has ended. \n Winner: %s \n Winning bid: %d", offer.Brand, offer.Model, winner, winningBid)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) GetNotificationByID(id uint) (*Notification, error) {
	notification, err := s.NotificationRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return notification, nil
}
