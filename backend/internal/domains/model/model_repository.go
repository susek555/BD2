package model

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"gorm.io/gorm"
)

type ModelRepositoryInterface interface {
	GetByManufacturerID(id uint) ([]models.Model, error)
	GetByManufacturerName(name string) ([]models.Model, error)
	GetByManufacturerAndModelName(manufacturerName, modelName string) (*models.Model, error)
}

type ModelRepository struct {
	DB *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepositoryInterface {
	return &ModelRepository{DB: db}
}

func (r *ModelRepository) GetByManufacturerID(id uint) ([]models.Model, error) {
	var models []models.Model
	err := r.DB.Where("manufacturer_id = ?", id).Find(&models).Error
	return models, err
}

func (r *ModelRepository) GetByManufacturerName(name string) ([]models.Model, error) {
	var models []models.Model
	err := r.DB.
		Preload("Manufacturer").
		Joins("JOIN manufacturers ON manufacturers.id = models.manufacturer_id").
		Where("manufacturers.name = ?", name).
		Find(&models).Error
	return models, err
}

func (r *ModelRepository) GetByManufacturerAndModelName(manufacturerName, modelName string) (*models.Model, error) {
	var model models.Model
	err := r.DB.
		Preload("Manufacturer").
		Joins("JOIN manufacturers ON manufacturers.id = models.manufacturer_id").
		Where("manufacturers.name = ? AND models.name = ?", manufacturerName, modelName).
		First(&model).Error
	return &model, err
}
