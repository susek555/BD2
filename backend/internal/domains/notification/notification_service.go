package notification

import (
	"fmt"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
)

type NotificationServiceInterface interface {
	CreateOutbidNotification(notification *models.Notification, amount int64, offer *models.Auction) error
	CreateEndAuctionNotification(notification *models.Notification, winner string, winningBid int64, offer *models.SaleOffer) error
	GetNotificationByID(id uint) (*models.Notification, error)
}

type NotificationService struct {
	NotificationRepository NotificationRepositoryInterface
}

func NewNotificationService(notificationRepository NotificationRepositoryInterface) NotificationServiceInterface {
	return &NotificationService{
		NotificationRepository: notificationRepository,
	}
}

func (s *NotificationService) CreateOutbidNotification(notification *models.Notification, amount int64, offer *models.Auction) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf("Someone outbid you on %s %s", offer.Offer.Car.Model.Manufacturer.Name, offer.Offer.Car.Model.Name)
	notification.Description = fmt.Sprintf("New price: %d", amount)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateEndAuctionNotification(notification *models.Notification, winner string, winningBid int64, offer *models.SaleOffer) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf("Auction ended for %s %s", offer.Car.Model.Manufacturer.Name, offer.Car.Model.Name)
	notification.Description = fmt.Sprintf("The auction for %s %s has ended. Winner: %s \n Winning bid: %d", offer.Car.Model.Manufacturer.Name, offer.Car.Model.Name, winner, winningBid)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	notification, err := s.NotificationRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return notification, nil
}
