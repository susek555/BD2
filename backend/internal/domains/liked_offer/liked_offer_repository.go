package liked_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type LikedOfferRepositoryInterface interface {
	Create(offer *models.LikedOffer) error
	Delete(offerID, userID uint) error
	GetByUserID(id uint) ([]models.LikedOffer, error)
	IsOfferLikedByUser(userID uint, offerID uint) error
}

type LikedOfferRepository struct {
	DB *gorm.DB
}

func NewLikedOfferRepository(db *gorm.DB) LikedOfferRepositoryInterface {
	return &LikedOfferRepository{DB: db}
}

func (r *LikedOfferRepository) Create(offer *models.LikedOffer) error {
	return r.DB.Create(offer).Error
}

func (r *LikedOfferRepository) Delete(offerID, userID uint) error {
	var likedOffer models.LikedOffer
	return r.DB.Where("offer_id = ? AND user_id = ?", offerID, userID).Delete(&likedOffer).Error
}

func (r *LikedOfferRepository) GetByUserID(id uint) ([]models.LikedOffer, error) {
	var likedOffers []models.LikedOffer
	err := r.DB.Where("user_id = ?", id).Find(&likedOffers).Error
	return likedOffers, err
}

func (r *LikedOfferRepository) IsOfferLikedByUser(offerID, userID uint) error {
	var likedOffer models.LikedOffer
	err := r.DB.Where("offer_id = ? AND user_id = ?", offerID, userID).First(&likedOffer).Error
	return err
}
