package models

type LikedOffer struct {
	UserID    uint       `json:"user_id"`
	OfferID   uint       `json:"offer_id"`
	SaleOffer *SaleOffer `gorm:"foreignKey:OfferID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User      *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
