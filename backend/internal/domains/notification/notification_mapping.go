package notification

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

func MapToClientNotification(not *models.Notification, userID uint) *models.ClientNotification {
	return &models.ClientNotification{
		NotificationID: not.ID,
		UserID:         userID,
		Seen:           false,
	}
}

func MapNotificationToDTO(notification *models.Notification, seen bool) *RetrieveNotificationDTO {
	if notification == nil {
		return nil
	}
	createdAt := notification.CreatedAt.Format(formats.DateTimeLayout)
	return &RetrieveNotificationDTO{
		ID:          notification.ID,
		Title:       notification.Title,
		Description: notification.Description,
		CreatedAt:   createdAt,
		OfferID:     notification.OfferID,
		Seen:        seen,
	}
}

func MapToNotificationDTO(clientNotification *models.ClientNotification) *RetrieveNotificationDTO {
	if clientNotification == nil {
		return nil
	}
	return MapNotificationToDTO(clientNotification.Notification, clientNotification.Seen)
}

func MapToNotificationsDTO(clientNotifications []models.ClientNotification, unseenNotifsCount uint, allNotifsCount uint) *NotificationsDTO {
	notifications := make([]RetrieveNotificationDTO, len(clientNotifications))
	for i, cn := range clientNotifications {
		notifications[i] = *MapNotificationToDTO(cn.Notification, cn.Seen)
	}
	return &NotificationsDTO{
		Notifications:     notifications,
		UnseenNotifsCount: unseenNotifsCount,
		AllNotifsCount:    allNotifsCount,
	}
}
