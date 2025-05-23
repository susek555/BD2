package models

type Notification struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	OfferID     uint       `json:"offer_id"`
	Offer       *SaleOffer `json:"sale_offer,omitempty" gorm:"foreignKey:OfferID;references:ID"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Date        string     `json:"date"`
}
