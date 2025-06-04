package notification

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type ClientNotificationRepositoryInterface interface {
	Create(clientNotification *models.ClientNotification) error
	GetByID(id uint) (*models.ClientNotification, error)
	GetAll() ([]models.ClientNotification, error)
	GetByUserID(userID uint) ([]models.ClientNotification, error)
	GetLatestByUserID(userID uint, count int) ([]models.ClientNotification, error)
}

type ClientNotificationRepository struct {
	DB *gorm.DB
}

func NewClientNotificationRepository(db *gorm.DB) ClientNotificationRepositoryInterface {
	return &ClientNotificationRepository{
		DB: db,
	}
}

func (r *ClientNotificationRepository) Create(clientNotification *models.ClientNotification) error {
	db := r.DB
	if err := db.Create(clientNotification).Error; err != nil {
		return err
	}
	return nil
}

func (r *ClientNotificationRepository) GetByID(id uint) (*models.ClientNotification, error) {
	db := r.DB
	var clientNotification models.ClientNotification
	if err := db.First(&clientNotification, id).Error; err != nil {
		return nil, err
	}
	return &clientNotification, nil
}

func (r *ClientNotificationRepository) GetAll() ([]models.ClientNotification, error) {
	db := r.DB
	var clientNotifications []models.ClientNotification
	if err := db.Find(&clientNotifications).Error; err != nil {
		return nil, err
	}
	return clientNotifications, nil
}

func (r *ClientNotificationRepository) GetByUserID(userID uint) ([]models.ClientNotification, error) {
	db := r.DB
	var clientNotifications []models.ClientNotification
	if err := db.Where("user_id = ?", userID).Find(&clientNotifications).Error; err != nil {
		return nil, err
	}
	return clientNotifications, nil
}

func (r *ClientNotificationRepository) GetLatestByUserID(userID uint, count int) ([]models.ClientNotification, error) {
	db := r.DB
	var clientNotifications []models.ClientNotification
	err := db.
		Joins("JOIN notifications ON notifications.id = client_notifications.notification_id").
		Where("client_notifications.user_id = ?", userID).
		Order("notifications.created_at DESC").
		Limit(count).
		Preload("Notification").
		Find(&clientNotifications).
		Error
	if err != nil {
		return nil, err
	}
	return clientNotifications, nil
}
