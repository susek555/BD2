package image

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type ImageRepositoryInterface interface {
	Create(image *models.Image) error
	BatchCreate(images []models.Image) error
	GetByURL(url string) (*models.Image, error)
	GetByOfferID(offerID uint) ([]models.Image, error)
	Delete(id uint) error
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

func (r *ImageRepository) BatchCreate(images []models.Image) error {
	return r.DB.Create(&images).Error
}

func (r *ImageRepository) GetByURL(url string) (*models.Image, error) {
	var image models.Image
	err := r.DB.Where("url = ?", url).Find(&image).Error
	return &image, err
}

func (r *ImageRepository) GetByOfferID(offerID uint) ([]models.Image, error) {
	var images []models.Image
	err := r.DB.Where("offer_id = ?", offerID).Find(&images).Error
	return images, err
}

func (r *ImageRepository) Delete(id uint) error {
	var image models.Image
	err := r.DB.Delete(&image, id).Error
	return err
}

func (r *ImageRepository) DeleteByOfferID(offerID uint) error {
	return r.DB.Where("offer_id = ?", offerID).Delete(&models.Image{}).Error
}
