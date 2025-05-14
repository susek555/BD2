package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]SaleOffer, *pagination.PaginationResponse, error)
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
	query, err := r.buildQuery(filter)
	if err != nil {
		return nil, nil, err
	}
	totalRecords, err := r.countTotalRecords(query)
	if err != nil {
		return nil, nil, err
	}
	paginationFunc, paginationResponse, err := pagination.Paginate(&filter.Pagination, totalRecords)
	if err != nil {
		return nil, nil, err
	}
	var saleOffers []SaleOffer
	if err := query.Scopes(paginationFunc).Find(&saleOffers).Error; err != nil {
		return nil, nil, err
	}
	return saleOffers, paginationResponse, nil
}

func (r *SaleOfferRepository) buildQuery(filter *OfferFilter) (*gorm.DB, error) {
	query := r.DB.
		Joins("JOIN cars on cars.id = sale_offers.car_id").
		Joins("LEFT JOIN auctions on auctions.offer_id = sale_offers.id").
		Preload("Auction").
		Preload("User").
		Preload("Car").
		Preload("Car.Model").
		Preload("Car.Model.Manufacturer")
	return filter.ApplyOfferFilters(query)
}

func (r *SaleOfferRepository) countTotalRecords(query *gorm.DB) (int64, error) {
	var totalRecords int64
	err := query.Model(&SaleOffer{}).Count(&totalRecords).Error
	return totalRecords, err
}
