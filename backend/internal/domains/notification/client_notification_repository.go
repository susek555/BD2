package notification

type ClientNotificationRepositoryInterface interface {
	Create(clientNotification *ClientNotification) (*ClientNotification, error)
	GetById(id uint) (*ClientNotification, error)
	GetAll() ([]ClientNotification, error)
}
