package model

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"gorm.io/gorm"
)

type ModelRepositoryInterface interface {
	GetByManufacturerID(id uint) ([]models.Model, error)       // Corrected spelling
	GetByManufacturerName(name string) ([]models.Model, error) // Corrected spelling
}

type ModelRepository struct {
	DB *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepositoryInterface {
	return &ModelRepository{DB: db}
}

func (r *ModelRepository) GetByManufacturerID(id uint) ([]models.Model, error) { // Corrected spelling
	var models []models.Model
	err := r.DB.Where("manufacturer_id = ?", id).Find(&models).Error
	return models, err
}

func (r *ModelRepository) GetByManufacturerName(name string) ([]models.Model, error) { // Corrected spelling
	var models []models.Model
	err := r.DB.
		Preload("Manufacturer").
		Joins("JOIN manufacturers ON manufacturers.id = models.manufacturer_id").
		Where("manufacturers.name = ?", name).
		Find(&models).Error
	return models, err
}
