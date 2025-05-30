package models

type Image struct {
	ID       uint       `json:"ID" gorm:"primaryKey"`
	OfferID  uint       `json:"offer_id"`
	Url      string     `json:"url"`
	PublicID string     `json:"public_id"`
	Offer    *SaleOffer `gorm:"foreignKey:OfferID;references:ID"`
}
