package image

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"gorm.io/gorm"
)

type ImageRepositoryInterface interface {
	Create(image *models.Image) error
	GetImagesByOfferID(offerID uint) ([]models.Image, error)
	DeleteByOfferID(offerID uint) error
}

type ImageRepository struct {
	DB *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepositoryInterface {
	return &ImageRepository{DB: db}
}

func (r *ImageRepository) Create(image *models.Image) error {
	return r.DB.Create(&image).Error
}

func (r *ImageRepository) GetImagesByOfferID(offerID uint) ([]models.Image, error) {
	var images []models.Image
	err := r.DB.Where("offer_id = ?", offerID).Find(&images).Error
	return images, err
}

func (r *ImageRepository) DeleteByOfferID(offerID uint) error {
	return r.DB.Where("offer_id = ?", offerID).Delete(&models.Image{}).Error
}
