package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *models.SaleOffer) error
	Update(offer *models.SaleOffer) error
	BuyOffer(offer *models.SaleOffer, buyerID uint) error
	GetByID(id uint) (*models.SaleOffer, error)
	GetViewByID(id uint) (*views.SaleOfferView, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)
	GetFiltered(filter *OfferFilter, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)
	UpdateStatus(offerID uint, status enums.Status) error
	SaveToPurchases(offerID uint, buyerID uint, finalPrice uint) error
	Delete(id uint) error
}

type SaleOfferRepository struct {
	DB *gorm.DB
}

func NewSaleOfferRepository(db *gorm.DB) SaleOfferRepositoryInterface {
	return &SaleOfferRepository{DB: db}
}

func (r *SaleOfferRepository) Create(offer *models.SaleOffer) error {
	return r.DB.Create(offer).Error
}

func (r *SaleOfferRepository) Update(offer *models.SaleOffer) error {
	return r.DB.Save(offer).Error
}

func (r *SaleOfferRepository) GetByID(id uint) (*models.SaleOffer, error) {
	var offer models.SaleOffer
	err := r.DB.Preload("Car").First(&offer, id).Error
	return &offer, err
}

func (r *SaleOfferRepository) GetViewByID(id uint) (*views.SaleOfferView, error) {
	var offerView views.SaleOfferView
	err := r.DB.Table("sale_offer_view").First(&offerView, id).Error
	return &offerView, err
}

func (r *SaleOfferRepository) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
	saleOffers, paginationResponse, err := pagination.PaginateResults[views.SaleOfferView](pagRequest, r.DB.Table("sale_offer_view").Where("user_id = ?", id))
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
	query := r.DB.Table("sale_offer_view").Where("status = ?", enums.PUBLISHED)
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, nil, err
	}
	saleOffers, paginationResponse, err := pagination.PaginateResults[views.SaleOfferView](pagRequest, query)
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) BuyOffer(offer *models.SaleOffer, buyerID uint) error {
	offer.Status = enums.SOLD
	if err := r.Update(offer); err != nil {
		return err
	}
	return r.SaveToPurchases(offer.ID, buyerID, offer.Price)
}

func (r *SaleOfferRepository) UpdateStatus(offerID uint, status enums.Status) error {
	return r.DB.Model(&models.SaleOffer{}).
		Where("id = ?", offerID).
		Update("status", status).
		Error
}

func (r *SaleOfferRepository) SaveToPurchases(offerID uint, buyerID uint, finalPrice uint) error {
	purchase := models.Purchase{
		OfferID:    offerID,
		BuyerID:    buyerID,
		FinalPrice: finalPrice,
		IssueDate:  time.Now(),
	}
	return r.DB.Create(&purchase).Error
}

func (r *SaleOfferRepository) Delete(id uint) error {
	return r.DB.Delete(&models.SaleOffer{}, id).Error
}
