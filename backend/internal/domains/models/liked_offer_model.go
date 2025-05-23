package models

type LikedOffer struct {
	UserID  uint `json:"user_id" gorm:"primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	OfferID uint `json:"offer_id" gorm:"primaryKey;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
