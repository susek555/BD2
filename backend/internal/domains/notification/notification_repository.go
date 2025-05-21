package notification

type NotificationRepositoryInterface interface {
	Create(notification *Notification) (*Notification, error)
	GetById(id uint) (*Notification, error)
	GetAll() ([]Notification, error)
}