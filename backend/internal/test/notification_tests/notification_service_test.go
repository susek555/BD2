package notification_tests

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
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
	amount := uint(27000)

	notificationRepo.createFunc = func(notif *models.Notification) error {
		assert.Equal(t, "Someone outbid you on Test Manufacturer Test Model", notif.Title)
		assert.Equal(t, "New price: 27000", notif.Description)
		assert.NotZero(t, notif.CreatedAt)
		return nil
	}

	err := service.CreateOutbidNotification(testNotification, amount, testAuction)

	assert.NoError(t, err)
}

func TestNotificationService_CreateOutbidNotification_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)
	testNotification := createSampleNotification()
	testAuction := createSampleAuction()
	amount := uint(27000)
	expectedError := errors.New("repository error")

	notificationRepo.createFunc = func(notif *models.Notification) error {
		return expectedError
	}

	err := service.CreateOutbidNotification(testNotification, amount, testAuction)

	assert.ErrorIs(t, err, expectedError)
}

func TestNotificationService_CreateEndAuctionNotification_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	testSaleOffer := createSampleSaleOffer()
	winner := "testuser123"
	winningBid := uint(28000)

	notificationRepo.createFunc = func(notif *models.Notification) error {
		assert.Equal(t, "Auction ended for Test Manufacturer Test Model", notif.Title)
		assert.Equal(t, "The auction for Test Manufacturer Test Model has ended. Winner: testuser123 Winning bid: 28000", notif.Description)
		assert.NotZero(t, notif.CreatedAt)
		return nil
	}

	err := service.CreateEndAuctionNotification(testNotification, winner, winningBid, testSaleOffer)

	assert.NoError(t, err)
}

func TestNotificationService_CreateEndAuctionNotification_NoBidsError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	testSaleOffer := createSampleSaleOffer()
	winner := ""
	winningBid := uint(0)

	err := service.CreateEndAuctionNotification(testNotification, winner, winningBid, testSaleOffer)

	assert.ErrorIs(t, err, notification.ErrNoBids)
}

func TestNotificationService_CreateEndAuctionNotification_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	testSaleOffer := createSampleSaleOffer()
	winner := "testuser123"
	winningBid := uint(28000)
	expectedError := errors.New("repository error")

	notificationRepo.createFunc = func(notif *models.Notification) error {
		return expectedError
	}

	err := service.CreateEndAuctionNotification(testNotification, winner, winningBid, testSaleOffer)

	assert.ErrorIs(t, err, expectedError)
}

func TestNotificationService_CreateBuyNotification_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	testSaleOffer := createSampleSaleOffer()
	buyerID := "buyer123"

	notificationRepo.createFunc = func(notif *models.Notification) error {
		assert.Equal(t, "The offer for Test Manufacturer Test Model has been bought", notif.Title)
		assert.Equal(t, "The offer has been bought by buyer123 for 25000", notif.Description)
		assert.NotZero(t, notif.CreatedAt)
		return nil
	}

	err := service.CreateBuyNotification(testNotification, buyerID, testSaleOffer)

	assert.NoError(t, err)
}

func TestNotificationService_CreateBuyNowNotification_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	testSaleOffer := createSampleSaleOffer()
	buyerID := "buyer123"

	notificationRepo.createFunc = func(notif *models.Notification) error {
		assert.Equal(t, "The auction for Test Manufacturer Test Model has been bought", notif.Title)
		assert.Equal(t, "The auction has been bought by buyer123 for 30000", notif.Description)
		assert.NotZero(t, notif.CreatedAt)
		return nil
	}

	err := service.CreateBuyNowNotification(testNotification, buyerID, testSaleOffer)

	assert.NoError(t, err)
}

func TestNotificationService_GetNotificationByID_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	expectedNotification := createSampleNotification()
	notificationID := uint(1)

	notificationRepo.getByIDFunc = func(id uint) (*models.Notification, error) {
		assert.Equal(t, notificationID, id)
		return expectedNotification, nil
	}

	result, err := service.GetNotificationByID(notificationID)

	assert.NoError(t, err)
	assert.Equal(t, expectedNotification, result)
}

func TestNotificationService_GetNotificationByID_NotFound(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	notificationID := uint(999)

	notificationRepo.getByIDFunc = func(id uint) (*models.Notification, error) {
		return nil, gorm.ErrRecordNotFound
	}

	result, err := service.GetNotificationByID(notificationID)

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, result)
}

func TestNotificationService_GetFilteredNotifications_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	filter := &notification.NotificationFilter{
		ReceiverID: func() *uint { id := uint(1); return &id }(),
	}

	expectedClientNotifications := []models.ClientNotification{
		*createSampleClientNotification(),
	}
	expectedPagination := &pagination.PaginationResponse{
		TotalRecords: 1,
		TotalPages:   1,
	}

	clientNotificationRepo.getFilteredFunc = func(f *notification.NotificationFilter) ([]models.ClientNotification, *pagination.PaginationResponse, error) {
		assert.Equal(t, filter, f)
		return expectedClientNotifications, expectedPagination, nil
	}

	result, err := service.GetFilteredNotifications(filter)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPagination, result.PaginationResponse)
	assert.Len(t, result.Notifications, 1)
}

func TestNotificationService_GetFilteredNotifications_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	filter := &notification.NotificationFilter{}
	expectedError := errors.New("repository error")

	clientNotificationRepo.getFilteredFunc = func(f *notification.NotificationFilter) ([]models.ClientNotification, *pagination.PaginationResponse, error) {
		return nil, nil, expectedError
	}

	result, err := service.GetFilteredNotifications(filter)

	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, result)
}

func TestNotificationService_UpdateSeenStatus_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	notificationID := uint(1)
	userID := uint(1)
	seen := true

	clientNotificationRepo.updateSeenStatusFunc = func(nID, uID uint, s bool) error {
		assert.Equal(t, notificationID, nID)
		assert.Equal(t, userID, uID)
		assert.Equal(t, seen, s)
		return nil
	}

	err := service.UpdateSeenStatus(notificationID, userID, seen)

	assert.NoError(t, err)
}

func TestNotificationService_UpdateSeenStatus_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	notificationID := uint(1)
	userID := uint(1)
	seen := true
	expectedError := errors.New("repository error")

	clientNotificationRepo.updateSeenStatusFunc = func(nID, uID uint, s bool) error {
		return expectedError
	}

	err := service.UpdateSeenStatus(notificationID, userID, seen)

	assert.ErrorIs(t, err, expectedError)
}

func TestNotificationService_UpdateSeenStatusForAll_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	userID := uint(1)
	seen := true

	clientNotificationRepo.updateSeenStatusForAllFunc = func(uID uint, s bool) error {
		assert.Equal(t, userID, uID)
		assert.Equal(t, seen, s)
		return nil
	}

	err := service.UpdateSeenStatusForAll(userID, seen)

	assert.NoError(t, err)
}

func TestNotificationService_UpdateSeenStatusForAll_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	userID := uint(1)
	seen := true
	expectedError := errors.New("repository error")

	clientNotificationRepo.updateSeenStatusForAllFunc = func(uID uint, s bool) error {
		return expectedError
	}

	err := service.UpdateSeenStatusForAll(userID, seen)

	assert.ErrorIs(t, err, expectedError)
}

func TestNotificationService_GetLatestNotificationsByUserID_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	userID := uint(1)
	count := uint(4)

	expectedClientNotifications := []models.ClientNotification{
		*createSampleClientNotification(),
	}
	expectedUnseenCount := uint(2)
	expectedAllCount := uint(5)

	clientNotificationRepo.getLatestByUserIDFunc = func(uID uint, c int) ([]models.ClientNotification, error) {
		assert.Equal(t, userID, uID)
		assert.Equal(t, 4, c) // Service hardcodes count to 4
		return expectedClientNotifications, nil
	}

	clientNotificationRepo.getUnseenCountByUserIdFunc = func(uID uint) (uint, error) {
		assert.Equal(t, userID, uID)
		return expectedUnseenCount, nil
	}

	clientNotificationRepo.getAllCountByUserIdFunc = func(uID uint) (uint, error) {
		assert.Equal(t, userID, uID)
		return expectedAllCount, nil
	}

	result, err := service.GetLatestNotificationsByUserID(userID, count)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUnseenCount, result.UnseenNotifsCount)
	assert.Equal(t, expectedAllCount, result.AllNotifsCount)
	assert.Len(t, result.Notifications, 1)
}

func TestNotificationService_GetLatestNotificationsByUserID_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	userID := uint(1)
	count := uint(4)
	expectedError := errors.New("repository error")

	clientNotificationRepo.getLatestByUserIDFunc = func(uID uint, c int) ([]models.ClientNotification, error) {
		return nil, expectedError
	}

	result, err := service.GetLatestNotificationsByUserID(userID, count)

	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, result)
}

func TestNotificationService_GetLatestNotificationsByUserID_UnseenCountError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	userID := uint(1)
	count := uint(4)
	expectedError := errors.New("unseen count error")

	clientNotificationRepo.getLatestByUserIDFunc = func(uID uint, c int) ([]models.ClientNotification, error) {
		return []models.ClientNotification{}, nil
	}

	clientNotificationRepo.getUnseenCountByUserIdFunc = func(uID uint) (uint, error) {
		return 0, expectedError
	}

	result, err := service.GetLatestNotificationsByUserID(userID, count)

	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, result)
}

func TestNotificationService_GetLatestNotificationsByUserID_AllCountError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	userID := uint(1)
	count := uint(4)
	expectedError := errors.New("all count error")

	clientNotificationRepo.getLatestByUserIDFunc = func(uID uint, c int) ([]models.ClientNotification, error) {
		return []models.ClientNotification{}, nil
	}

	clientNotificationRepo.getUnseenCountByUserIdFunc = func(uID uint) (uint, error) {
		return 0, nil
	}

	clientNotificationRepo.getAllCountByUserIdFunc = func(uID uint) (uint, error) {
		return 0, expectedError
	}

	result, err := service.GetLatestNotificationsByUserID(userID, count)

	assert.ErrorIs(t, err, expectedError)
	assert.Nil(t, result)
}

func TestNotificationService_SaveNotificationToClient_Success(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	userID := uint(1)

	clientNotificationRepo.createFunc = func(clientNotification *models.ClientNotification) error {
		assert.Equal(t, testNotification.ID, clientNotification.NotificationID)
		assert.Equal(t, userID, clientNotification.UserID)
		assert.False(t, clientNotification.Seen)
		return nil
	}

	err := service.SaveNotificationToClient(testNotification, userID)

	assert.NoError(t, err)
}

func TestNotificationService_SaveNotificationToClient_RepositoryError(t *testing.T) {
	notificationRepo := &mockNotificationRepository{}
	clientNotificationRepo := &mockClientNotificationRepository{}
	service := notification.NewNotificationService(notificationRepo, clientNotificationRepo)

	testNotification := createSampleNotification()
	userID := uint(1)
	expectedError := errors.New("repository error")

	clientNotificationRepo.createFunc = func(clientNotification *models.ClientNotification) error {
		return expectedError
	}
	err := service.SaveNotificationToClient(testNotification, userID)

	assert.ErrorIs(t, err, expectedError)
}
