package notification_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

// Custom mock for NotificationRepository
type mockNotificationRepository struct {
	createFunc  func(notification *models.Notification) error
	getByIDFunc func(id uint) (*models.Notification, error)
	getAllFunc  func() ([]models.Notification, error)
}

func (m *mockNotificationRepository) Create(notif *models.Notification) error {
	if m.createFunc != nil {
		return m.createFunc(notif)
	}
	return nil
}

func (m *mockNotificationRepository) GetByID(id uint) (*models.Notification, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockNotificationRepository) GetAll() ([]models.Notification, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil, nil
}

// Custom mock for ClientNotificationRepository
type mockClientNotificationRepository struct {
	createFunc                 func(clientNotification *models.ClientNotification) error
	getByIDFunc                func(id uint) (*models.ClientNotification, error)
	getAllFunc                 func() ([]models.ClientNotification, error)
	getByUserIDFunc            func(userID uint) ([]models.ClientNotification, error)
	getLatestByUserIDFunc      func(userID uint, count int) ([]models.ClientNotification, error)
	getUnseenCountByUserIdFunc func(userId uint) (uint, error)
	getAllCountByUserIdFunc    func(userId uint) (uint, error)
	getFilteredFunc            func(filter *notification.NotificationFilter) ([]models.ClientNotification, *pagination.PaginationResponse, error)
	updateSeenStatusFunc       func(notificationID, userID uint, seen bool) error
	updateSeenStatusForAllFunc func(userID uint, seen bool) error
}

func (m *mockClientNotificationRepository) Create(clientNotification *models.ClientNotification) error {
	if m.createFunc != nil {
		return m.createFunc(clientNotification)
	}
	return nil
}

func (m *mockClientNotificationRepository) GetByID(id uint) (*models.ClientNotification, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockClientNotificationRepository) GetAll() ([]models.ClientNotification, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil, nil
}

func (m *mockClientNotificationRepository) GetByUserID(userID uint) ([]models.ClientNotification, error) {
	if m.getByUserIDFunc != nil {
		return m.getByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *mockClientNotificationRepository) GetLatestByUserID(userID uint, count int) ([]models.ClientNotification, error) {
	if m.getLatestByUserIDFunc != nil {
		return m.getLatestByUserIDFunc(userID, count)
	}
	return nil, nil
}

func (m *mockClientNotificationRepository) GetUnseenCountByUserId(userId uint) (uint, error) {
	if m.getUnseenCountByUserIdFunc != nil {
		return m.getUnseenCountByUserIdFunc(userId)
	}
	return 0, nil
}

func (m *mockClientNotificationRepository) GetAllCountByUserId(userId uint) (uint, error) {
	if m.getAllCountByUserIdFunc != nil {
		return m.getAllCountByUserIdFunc(userId)
	}
	return 0, nil
}

func (m *mockClientNotificationRepository) GetFiltered(filter *notification.NotificationFilter) ([]models.ClientNotification, *pagination.PaginationResponse, error) {
	if m.getFilteredFunc != nil {
		return m.getFilteredFunc(filter)
	}
	return nil, nil, nil
}

func (m *mockClientNotificationRepository) UpdateSeenStatus(notificationID, userID uint, seen bool) error {
	if m.updateSeenStatusFunc != nil {
		return m.updateSeenStatusFunc(notificationID, userID, seen)
	}
	return nil
}

func (m *mockClientNotificationRepository) UpdateSeenStatusForAll(userID uint, seen bool) error {
	if m.updateSeenStatusForAllFunc != nil {
		return m.updateSeenStatusForAllFunc(userID, seen)
	}
	return nil
}

// Helper functions to create test data
func createSampleNotification() *models.Notification {
	return &models.Notification{
		ID:          1,
		OfferID:     1,
		Title:       "Test Notification",
		Description: "Test Description",
		CreatedAt:   time.Now().UTC(),
	}
}

func createSampleSaleOffer() *models.SaleOffer {
	return &models.SaleOffer{
		ID:     1,
		Price:  25000,
		Status: enums.PUBLISHED,
		Car: &models.Car{
			Model: &models.Model{
				Name: "Test Model",
				Manufacturer: &models.Manufacturer{
					Name: "Test Manufacturer",
				},
			},
		},
		Auction: &models.Auction{
			BuyNowPrice: 30000,
		},
	}
}

func createSampleAuction() *models.Auction {
	return &models.Auction{
		OfferID:     1,
		BuyNowPrice: 30000,
		Offer: &models.SaleOffer{
			Car: &models.Car{
				Model: &models.Model{
					Name: "Test Model",
					Manufacturer: &models.Manufacturer{
						Name: "Test Manufacturer",
					},
				},
			},
		},
	}
}

func createSampleClientNotification() *models.ClientNotification {
	return &models.ClientNotification{
		ID:             1,
		NotificationID: 1,
		UserID:         1,
		Seen:           false,
		Notification:   createSampleNotification(),
	}
}

// TESTS

func TestNotificationService_CreateOutbidNotification_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	testAuction := createSampleAuction()
	amount := int64(27000)

	notificationRepo.createFunc = func(notif *models.Notification) error {
		assert.Equal(t, "Someone outbid you on Test Manufacturer Test Model", notif.Title)
		assert.Equal(t, "New price: 27000", notif.Description)
		assert.NotZero(t, notif.CreatedAt)
		return nil
	}

	err := service.CreateOutbidNotification(testNotification, amount, testAuction)

	assert.NoError(t, err)
}
