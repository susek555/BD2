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
	GetByID(id uint) (*models.SaleOffer, error)
	GetViewByID(id uint) (*views.SaleOfferView, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)
	GetFiltered(filter *OfferFilter, pagRequest *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)
	GetAllActiveAuctions() ([]models.SaleOffer, error)
	BuyOffer(offerID uint, buyerID uint) (*models.SaleOffer, error)
	UpdateStatus(offerID uint, status enums.Status) error
	SaveToPurchases(offerID uint, buyerID uint, finalPrice uint) error
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

func (r *SaleOfferRepository) GetAllActiveAuctions() ([]models.SaleOffer, error) {
	var auctions []models.SaleOffer
	err := r.DB.
		Preload("Auction").
		Joins("JOIN auctions ON auctions.offer_id = sale_offers.id").
		Where("auctions.offer_id IS NOT NULL").
		Where("auctions.date_end > NOW()").
		Find(&auctions).
		Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (r *SaleOfferRepository) BuyOffer(offerID uint, buyerID uint) (*models.SaleOffer, error) {
	offer, err := r.GetByID(offerID)
	if err != nil {
		return nil, err
	}
	if offer.Status == enums.SOLD {
		return nil, ErrOfferAlreadySold
	}
	err = r.UpdateStatus(offerID, enums.SOLD)
	if err != nil {
		return nil, err
	}
	return offer, r.SaveToPurchases(offerID, buyerID, offer.Price)
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

func (r *SaleOfferRepository) buildBaseQuery() *gorm.DB {
	query := r.DB.
		Preload("Auction").
		Preload("User").
		Preload("Car").
		Preload("Car.Model").
		Preload("Car.Model.Manufacturer")
	return query
}
