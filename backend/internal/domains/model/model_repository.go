package model

import "gorm.io/gorm"

type ModelRepositoryInterface interface {
	GetByManufacturerID(id uint) ([]Model, error)       // Corrected spelling
	GetByManufacturerName(name string) ([]Model, error) // Corrected spelling
}

type ModelRepository struct {
	DB *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepositoryInterface {
	return &ModelRepository{DB: db}
}

func (r *ModelRepository) GetByManufacturerID(id uint) ([]Model, error) { // Corrected spelling
	var models []Model
	err := r.DB.Where("manufacturer_id = ?", id).Find(&models).Error
	return models, err
}

func (r *ModelRepository) GetByManufacturerName(name string) ([]Model, error) { // Corrected spelling
	var models []Model
	err := r.DB.
		Preload("Manufacturer").
		Joins("JOIN manufacturers ON manufacturers.id = models.manufacturer_id").
		Where("manufacturers.name = ?", name).
		Find(&models).Error
	return models, err
}
