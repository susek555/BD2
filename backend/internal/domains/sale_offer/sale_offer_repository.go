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
	var saleOffers []SaleOffer
	var totalRecords int64
	query := r.DB.
		Joins("JOIN cars on cars.id = sale_offers.car_id").
		Joins("LEFT JOIN auctions on auctions.offer_id = sale_offers.id").
		Preload("Auction").
		Preload("Car").
		Preload("Car.Model").
		Preload("Car.Model.Manufacturer")
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, nil, err
	}
	if err := query.Model(&SaleOffer{}).Count(&totalRecords).Error; err != nil {
		return nil, nil, err
	}
	if err := query.Scopes(pagination.Paginate(&filter.Pagination)).Find(&saleOffers).Error; err != nil {
		return nil, nil, err
	}

	return saleOffers, &pagination.PaginationResponse{TotalRecords: totalRecords, TotalPages: len(saleOffers)/filter.Pagination.PageSize + 1}, nil
}
