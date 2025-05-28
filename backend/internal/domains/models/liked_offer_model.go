package models

type LikedOffer struct {
	UserID    uint       `json:"user_id"`
	OfferID   uint       `json:"offer_id"`
	SaleOffer *SaleOffer `gorm:"foreignKey:OfferID;references:ID"`
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
}
