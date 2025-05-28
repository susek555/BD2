package notification

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

func MapToClientNotification(not *models.Notification, userID uint) *models.ClientNotification {
	return &models.ClientNotification{
		NotificationID: not.ID,
		UserID:         userID,
		Seen:           false,
	}
}

func MapToNotification(clientNotification *models.ClientNotification) *models.Notification {
	if clientNotification.Notification == nil {
		return nil
	}
	return clientNotification.Notification
}

func MapNotificationToDTO(notification *models.Notification) *RetrieveNotificationDTO {
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
	}
}

func MapToNotifications(clientNotifications []models.ClientNotification) []models.Notification {
	notifications := make([]models.Notification, len(clientNotifications))
	for i, cn := range clientNotifications {
		notifications[i] = *MapToNotification(&cn)
	}
	return notifications
}
