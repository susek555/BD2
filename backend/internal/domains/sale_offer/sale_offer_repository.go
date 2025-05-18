package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]SaleOffer, *pagination.PaginationResponse, error)
	GetByID(id uint) (*SaleOffer, error)
	GetByUserID(id uint, pagination *pagination.PaginationRequest) ([]SaleOffer, *pagination.PaginationResponse, error)
}

type SaleOfferRepository struct {
	DB *gorm.DB
}

func NewSaleOfferRepository(db *gorm.DB) SaleOfferRepositoryInterface {
	return &SaleOfferRepository{DB: db}
}

func (r *SaleOfferRepository) Create(offer *SaleOffer) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(offer).Error
}

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter) ([]SaleOffer, *pagination.PaginationResponse, error) {
	query := r.buildBaseQuery().
		Joins("JOIN cars on cars.offer_id = sale_offers.id").
		Joins("LEFT JOIN auctions on auctions.offer_id = sale_offers.id")
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, nil, err
	}
	saleOffers, paginationResponse, err := pagination.PaginateResults[SaleOffer](&filter.Pagination, query)
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) GetByID(id uint) (*SaleOffer, error) {
	var offer SaleOffer
	err := r.buildBaseQuery().First(&offer, id).Error
	return &offer, err
}

func (r *SaleOfferRepository) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) ([]SaleOffer, *pagination.PaginationResponse, error) {
	saleOffers, paginationResponse, err := pagination.PaginateResults[SaleOffer](pagRequest, r.buildBaseQuery().Where("user_id = ?", id))
	if err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
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
