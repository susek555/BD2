package notification

type NotificationServiceInterface interface {
	CreateOutbidNotification(notification *Notification) error
	CreateEndAuctionNotification(notification *Notification) error
}

type NotificationService struct {
	NotificationRepository NotificationRepositoryInterface
}

func NewNotificationService(notificationRepository NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{
		NotificationRepository: notificationRepository,
	}
}
