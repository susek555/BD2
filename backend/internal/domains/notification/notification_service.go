package notification

import (
	"fmt"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

//go:generate mockery --name=NotificationServiceInterface --output=../../test/mocks --case=snake --with-expecter

type NotificationServiceInterface interface {
	CreateOutbidNotification(notification *models.Notification, amount int64, offer *models.Auction) error
	CreateEndAuctionNotification(notification *models.Notification, winner string, winningBid int64, offer *models.SaleOffer) error
	CreateBuyNotication(notification *models.Notification, buyerID string, offer *models.SaleOffer) error
	CreateBuyNowNotification(notification *models.Notification, buyerID string, offer *models.Auction) error
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
	notification.Title = fmt.Sprintf(OutbidTitleTemplate, offer.Offer.Car.Model.Manufacturer.Name, offer.Offer.Car.Model.Name)
	notification.Description = fmt.Sprintf(OutbidDescriptionTemplate, amount)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateEndAuctionNotification(notification *models.Notification, winner string, winningBid int64, offer *models.SaleOffer) error {
	if winningBid == 0 && winner == "" {
		return ErrNoBids
	}
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(EndAuctionTitleTemplate, offer.Car.Model.Manufacturer.Name, offer.Car.Model.Name)
	notification.Description = fmt.Sprintf(EndAuctionDescriptionTemplate, offer.Car.Model.Manufacturer.Name, offer.Car.Model.Name, winner, winningBid)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateBuyNotication(notification *models.Notification, buyerID string, offer *models.SaleOffer) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(BuyOfferTitleTemplate, offer.Car.Model.Manufacturer.Name, offer.Car.Model.Name)
	notification.Description = fmt.Sprintf(BuyOfferDescriptionTemplate, buyerID, offer.Price)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateBuyNowNotification(notification *models.Notification, buyerID string, offer *models.Auction) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(BuyNowTitleTemplate, offer.Offer.Car.Model.Manufacturer.Name, offer.Offer.Car.Model.Name)
	notification.Description = fmt.Sprintf(BuyNowDescriptionTemplate, buyerID, offer.BuyNowPrice)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	notification, err := s.NotificationRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return notification, nil
}
