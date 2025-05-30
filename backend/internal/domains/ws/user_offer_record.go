package ws

type UserOfferRecord struct {
	OfferID string `json:"offer_id" gorm:"column:offer_id"`
	UserID  string `json:"user_id" gorm:"column:user_id"`
}
