package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]SaleOffer, *pagination.PaginationResponse, error)
	GetByID(id uint) (*SaleOffer, error)
	GetByUserID(id uint) ([]SaleOffer, error)
}

type SaleOfferRepository struct {
	DB *gorm.DB
}

func NewSaleOfferRepository(db *gorm.DB) SaleOfferRepositoryInterface {
	return &SaleOfferRepository{DB: db}
}

func (r *SaleOfferRepository) Create(offer *SaleOffer) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(offer.Car).Error; err != nil {
			return err
		}
		if err := tx.Create(offer).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter) ([]SaleOffer, *pagination.PaginationResponse, error) {
	query := r.buildQuery()
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
	err := r.DB.First(&offer, id).Error
	return &offer, err
}

func (r *SaleOfferRepository) GetByUserID(id uint) ([]SaleOffer, error) {
	var offers []SaleOffer
	err := r.DB.Where("user_id = ?", id).Find(&offers).Error
	return offers, err
}

func (r *SaleOfferRepository) buildQuery() *gorm.DB {
	query := r.DB.
		Joins("JOIN cars on cars.offer_id = sale_offers.id").
		Joins("LEFT JOIN auctions on auctions.offer_id = sale_offers.id").
		Preload("Auction").
		Preload("User").
		Preload("Car").
		Preload("Car.Model").
		Preload("Car.Model.Manufacturer")
	return query
}
