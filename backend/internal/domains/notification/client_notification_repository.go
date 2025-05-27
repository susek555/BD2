package notification

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"gorm.io/gorm"
)

type ClientNotificationRepositoryInterface interface {
	Create(clientNotification *models.ClientNotification) error
	GetById(id uint) (*models.ClientNotification, error)
	GetAll() ([]models.ClientNotification, error)
	GetByUserId(userId uint) ([]models.ClientNotification, error)
	GetLatestByUserId(userId uint, count int) ([]models.ClientNotification, error)
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

func (r *ClientNotificationRepository) GetById(id uint) (*models.ClientNotification, error) {
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

func (r *ClientNotificationRepository) GetByUserId(userId uint) ([]models.ClientNotification, error) {
	db := r.DB
	var clientNotifications []models.ClientNotification
	if err := db.Where("user_id = ?", userId).Find(&clientNotifications).Error; err != nil {
		return nil, err
	}
	return clientNotifications, nil
}

func (r *ClientNotificationRepository) GetLatestByUserId(userId uint, count int) ([]models.ClientNotification, error) {
	db := r.DB
	var clientNotifications []models.ClientNotification
	err := db.
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(count).
		Preload("Notification").
		Find(&clientNotifications).
		Error
	if err != nil {
		return nil, err
	}
	return clientNotifications, nil
}
