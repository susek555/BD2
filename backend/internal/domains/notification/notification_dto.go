package notification

import "github.com/susek555/BD2/car-dealer-api/pkg/pagination"

type RetrieveNotificationDTO struct {
	ID          uint   `json:"id"`
	OfferID     uint   `json:"offer_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Seen        bool   `json:"seen"`
}

type NotificationsDTO struct {
	Notifications     []RetrieveNotificationDTO `json:"notifications"`
	UnseenNotifsCount uint                      `json:"unseen_notifs_count"`
	AllNotifsCount    uint                      `json:"all_notifs_count"`
}

type RetrieveNotificationsWithPagination struct {
	Notifications      []RetrieveNotificationDTO      `json:"notifications"`
	PaginationResponse *pagination.PaginationResponse `json:"pagination"`
}
