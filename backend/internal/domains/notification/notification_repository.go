package notification

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	Create(notification *models.Notification) error
	GetByID(id uint) (*models.Notification, error)
	GetAll() ([]models.Notification, error)
}

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepositoryInterface {
	return &NotificationRepository{
		DB: db,
	}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	db := r.DB
	if err := db.Create(notification).Error; err != nil {
		return err
	}
	return nil
}

func (r *NotificationRepository) GetByID(id uint) (*models.Notification, error) {
	db := r.DB
	var notification models.Notification
	if err := db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) GetAll() ([]models.Notification, error) {
	db := r.DB
	var notifications []models.Notification
	if err := db.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
