package notification

import "github.com/susek555/BD2/car-dealer-api/internal/domains/models"

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

