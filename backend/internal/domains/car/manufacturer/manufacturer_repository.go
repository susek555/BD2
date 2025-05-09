package manufacturer

import "gorm.io/gorm"

type ManufacturerRepositoryInterface interface {
	GetAll() ([]Manufacturer, error)
}

type ManufacturerRepository struct {
	DB *gorm.DB
}

func NewManufacturerRepository(db *gorm.DB) ManufacturerRepositoryInterface {
	return &ManufacturerRepository{DB: db}
}

func (r *ManufacturerRepository) GetAll() ([]Manufacturer, error) {
	var manufacturers []Manufacturer
	err := r.DB.Find(&manufacturers).Error
	return manufacturers, err
}
