package manufacturer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"gorm.io/gorm"
)

type ManufacturerRepositoryInterface interface {
	GetAll() ([]models.Manufacturer, error)
}

type ManufacturerRepository struct {
	DB *gorm.DB
}

func NewManufacturerRepository(db *gorm.DB) ManufacturerRepositoryInterface {
	return &ManufacturerRepository{DB: db}
}

func (r *ManufacturerRepository) GetAll() ([]models.Manufacturer, error) {
	var manufacturers []models.Manufacturer
	err := r.DB.Find(&manufacturers).Error
	return manufacturers, err
}
