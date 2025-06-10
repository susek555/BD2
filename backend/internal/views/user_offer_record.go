package views

type UserOfferRecord struct {
	OfferID uint `json:"offer_id" gorm:"column:offer_id"`
	UserID  uint `json:"user_id" gorm:"column:user_id"`
}
