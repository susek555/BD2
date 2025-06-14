package notification

import (
	"fmt"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

//go:generate mockery --name=NotificationServiceInterface --output=../../test/mocks --case=snake --with-expecter

type NotificationServiceInterface interface {
	CreateOutbidNotification(notification *models.Notification, amount uint, offer SaleOfferInterface) error
	CreateEndAuctionNotification(notification *models.Notification, winner string, winningBid uint, offer SaleOfferInterface) error
	CreateBuyNotification(notification *models.Notification, buyerID string, offer SaleOfferInterface) error
	CreateBuyNowNotification(notification *models.Notification, buyerID string, offer SaleOfferInterface) error
	GetNotificationByID(id uint) (*models.Notification, error)
	GetFilteredNotifications(filter *NotificationFilter) (*RetrieveNotificationsWithPagination, error)
	UpdateSeenStatus(notificationID, userID uint, seen bool) error
	UpdateSeenStatusForAll(userID uint, seen bool) error
	GetLatestNotificationsByUserID(userID uint, count uint) (*NotificationsDTO, error)
	SaveNotificationToClient(notification *models.Notification, userID uint) error
}

type NotificationService struct {
	NotificationRepository       NotificationRepositoryInterface
	ClientNotificationRepository ClientNotificationRepositoryInterface
}

func NewNotificationService(notificationRepository NotificationRepositoryInterface, clientNotification ClientNotificationRepositoryInterface) NotificationServiceInterface {
	return &NotificationService{
		NotificationRepository:       notificationRepository,
		ClientNotificationRepository: clientNotification,
	}
}

func (s *NotificationService) CreateOutbidNotification(notification *models.Notification, amount uint, offer SaleOfferInterface) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(OutbidTitleTemplate, offer.GetBrand(), offer.GetModel())
	notification.Description = fmt.Sprintf(OutbidDescriptionTemplate, amount)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateEndAuctionNotification(notification *models.Notification, winner string, winningBid uint, offer SaleOfferInterface) error {
	if winningBid == 0 && winner == "" {
		return ErrNoBids
	}
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(EndAuctionTitleTemplate, offer.GetBrand(), offer.GetModel())
	notification.Description = fmt.Sprintf(EndAuctionDescriptionTemplate, offer.GetBrand(), offer.GetModel(), winner, winningBid)
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateBuyNotification(notification *models.Notification, buyerID string, offer SaleOfferInterface) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(BuyOfferTitleTemplate, offer.GetBrand(), offer.GetModel())
	notification.Description = fmt.Sprintf(BuyOfferDescriptionTemplate, buyerID, offer.GetPrice())
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) CreateBuyNowNotification(notification *models.Notification, buyerID string, offer SaleOfferInterface) error {
	notification.CreatedAt = time.Now().UTC()
	notification.Title = fmt.Sprintf(BuyNowTitleTemplate, offer.GetBrand(), offer.GetModel())
	notification.Description = fmt.Sprintf(BuyNowDescriptionTemplate, buyerID, offer.GetPrice())
	return s.NotificationRepository.Create(notification)
}

func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	notification, err := s.NotificationRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) GetFilteredNotifications(filter *NotificationFilter) (*RetrieveNotificationsWithPagination, error) {
	notifications, pagResponse, err := s.ClientNotificationRepository.GetFiltered(filter)
	if err != nil {
		return nil, err
	}
	notificationsDTO := mapping.MapSliceToDTOs(notifications, MapToNotificationDTO)
	return &RetrieveNotificationsWithPagination{
		Notifications:      notificationsDTO,
		PaginationResponse: pagResponse,
	}, nil
}

func (s *NotificationService) UpdateSeenStatus(notificationID, userID uint, seen bool) error {
	if err := s.ClientNotificationRepository.UpdateSeenStatus(notificationID, userID, seen); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) UpdateSeenStatusForAll(userID uint, seen bool) error {
	if err := s.ClientNotificationRepository.UpdateSeenStatusForAll(userID, seen); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) GetLatestNotificationsByUserID(userID uint, count uint) (*NotificationsDTO, error) {
	notifications, err := s.ClientNotificationRepository.GetLatestByUserID(userID, 4)
	if err != nil {
		return nil, err
	}
	unseenCount, err := s.ClientNotificationRepository.GetUnseenCountByUserId(userID)
	if err != nil {
		return nil, err
	}
	allCount, err := s.ClientNotificationRepository.GetAllCountByUserId(userID)
	if err != nil {
		return nil, err
	}
	return MapToNotificationsDTO(notifications, unseenCount, allCount), nil
}

func (s *NotificationService) SaveNotificationToClient(notification *models.Notification, userID uint) error {
	clientNotification := MapToClientNotification(notification, userID)
	if err := s.ClientNotificationRepository.Create(clientNotification); err != nil {
		return err
	}
	return nil
}
