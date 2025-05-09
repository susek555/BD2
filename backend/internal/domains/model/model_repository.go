package model

import "gorm.io/gorm"

type ModelRepositoryInterface interface {
	GetByManufactuerID(id uint) ([]Model, error)
	GetByManufacuterName(name string) ([]Model, error)
}

type ModelRepository struct {
	DB *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepositoryInterface {
	return &ModelRepository{DB: db}
}

func (r *ModelRepository) GetByManufactuerID(id uint) ([]Model, error) {
	var Models []Model
	err := r.DB.Where("manufacturer_id = ?", id).Find(&Models).Error
	return Models, err
}

func (r *ModelRepository) GetByManufacuterName(name string) ([]Model, error) {
	var Models []Model
	err := r.DB.Preload("Manufacturer").
		Joins("JOIN manufacturers ON manufacturers.id = Models.manufacturer_id").
		Where("manufacturers.name = ?", name).
		Find(&Models).Error
	return Models, err
}
