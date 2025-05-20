package liked_offer

import "gorm.io/gorm"

type LikedOfferRepositoryInterface interface {
	Create(offer *LikedOffer) error
	Delete(offerID, userID uint) error
	GetByUserID(id uint) ([]LikedOffer, error)
	IsOfferLikedByUser(userID uint, offerID uint) bool
}

type LikedOfferRepository struct {
	DB *gorm.DB
}

func NewLikedOfferRepository(db *gorm.DB) LikedOfferRepositoryInterface {
	return &LikedOfferRepository{DB: db}
}

func (r *LikedOfferRepository) Create(offer *LikedOffer) error {
	return r.DB.Create(offer).Error
}

func (r *LikedOfferRepository) Delete(offerID, userID uint) error {
	var likedOffer LikedOffer
	return r.DB.Where("offer_id = ? AND user_id = ?", offerID, userID).Delete(&likedOffer).Error
}

func (r *LikedOfferRepository) GetByUserID(id uint) ([]LikedOffer, error) {
	var likedOffers []LikedOffer
	err := r.DB.Where("user_id = ?", id).Find(&likedOffers).Error
	return likedOffers, err
}

func (r *LikedOfferRepository) IsOfferLikedByUser(offerID, userID uint) bool {
	var likedOffer LikedOffer
	err := r.DB.Where("offer_id = ? AND user_id = ?", offerID, userID).First(&likedOffer).Error
	return err == nil
}
