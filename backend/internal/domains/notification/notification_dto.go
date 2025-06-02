package notification

type RetrieveNotificationDTO struct {
	ID          uint   `json:"id"`
	OfferID     uint   `json:"offer_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Seen        bool   `json:"seen"`
}
