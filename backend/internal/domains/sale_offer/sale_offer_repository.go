package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *models.SaleOffer) error
	Update(offer *models.SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]models.SaleOffer, *pagination.PaginationResponse, error)
	GetByID(id uint) (*models.SaleOffer, error)
	GetByUserID(id uint, pagination *pagination.PaginationRequest) ([]models.SaleOffer, *pagination.PaginationResponse, error)
	GetAllActiveAuctions() ([]models.SaleOffer, error)
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

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter) ([]models.SaleOffer, *pagination.PaginationResponse, error) {
	query := r.buildBaseQuery().
		Where("sale_offers.status = ?", enums.PUBLISHED).
		Joins("JOIN cars on cars.offer_id = sale_offers.id").
		Joins("LEFT JOIN auctions on auctions.offer_id = sale_offers.id")
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, nil, err
	}
	saleOffers, paginationResponse, err := pagination.PaginateResults[models.SaleOffer](&filter.Pagination, query)
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetByID(id uint) (*models.SaleOffer, error) {
	var offer models.SaleOffer
	err := r.buildBaseQuery().First(&offer, id).Error
	return &offer, err
}

func (r *SaleOfferRepository) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]models.SaleOffer, *pagination.PaginationResponse, error) {
	saleOffers, paginationResponse, err := pagination.PaginateResults[models.SaleOffer](pagRequest, r.buildBaseQuery().Where("user_id = ?", id))
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
		Where("sale_offers.status = ?", models.PUBLISHED).
		Where("auctions.date_end > NOW()").
		Find(&auctions).
		Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
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
