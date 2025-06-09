package purchase

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type PurchaseRepositoryInterface interface {
	Create(purchase *models.Purchase) error
	GetByID(id uint) (*models.Purchase, error)
}

type PurchaseRepository struct {
	DB *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepositoryInterface {
	return &PurchaseRepository{DB: db}
}

func (r *PurchaseRepository) Create(purchase *models.Purchase) error {
	return r.DB.Create(&purchase).Error
}

func (r *PurchaseRepository) GetByID(id uint) (*models.Purchase, error) {
	var purchase models.Purchase
	if err := r.DB.First(&purchase, id).Error; err != nil {
		return nil, err
	}
	return &purchase, nil
}
