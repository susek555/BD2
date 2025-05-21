package notification

import "gorm.io/gorm"

type ClientNotificationRepositoryInterface interface {
	Create(clientNotification *ClientNotification) error
	GetById(id uint) (*ClientNotification, error)
	GetAll() ([]ClientNotification, error)
}

type ClientNotificationRepository struct {
	DB *gorm.DB
}

func NewClientNotificationRepository(db *gorm.DB) ClientNotificationRepositoryInterface {
	return &ClientNotificationRepository{
		DB: db,
	}
}

func (r *ClientNotificationRepository) Create(clientNotification *ClientNotification) error {
	db := r.DB
	if err := db.Create(clientNotification).Error; err != nil {
		return err
	}
	return nil
}

func (r *ClientNotificationRepository) GetById(id uint) (*ClientNotification, error) {
	db := r.DB
	var clientNotification ClientNotification
	if err := db.First(&clientNotification, id).Error; err != nil {
		return nil, err
	}
	return &clientNotification, nil
}

func (r *ClientNotificationRepository) GetAll() ([]ClientNotification, error) {
	db := r.DB
	var clientNotifications []ClientNotification
	if err := db.Find(&clientNotifications).Error; err != nil {
		return nil, err
	}
	return clientNotifications, nil
}
