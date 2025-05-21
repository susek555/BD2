package notification

import "gorm.io/gorm"

type NotificationRepositoryInterface interface {
	Create(notification *Notification) error
	GetById(id uint) (*Notification, error)
	GetAll() ([]Notification, error)
}

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepositoryInterface {
	return &NotificationRepository{
		DB: db,
	}
}

func (r *NotificationRepository) Create(notification *Notification) error {
	db := r.DB
	if err := db.Create(notification).Error; err != nil {
		return err
	}
	return nil
}

func (r *NotificationRepository) GetById(id uint) (*Notification, error) {
	db := r.DB
	var notification Notification
	if err := db.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) GetAll() ([]Notification, error) {
	db := r.DB
	var notifications []Notification
	if err := db.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
