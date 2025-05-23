package models

type Image struct {
	ID      uint       `json:"ID" gorm:"primaryKey"`
	OfferID uint       `json:"offer_id"`
	Url     string     `json:"url"`
	Offer   *SaleOffer `gorm:"foreignKey:OfferID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
