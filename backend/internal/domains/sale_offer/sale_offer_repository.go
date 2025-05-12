package sale_offer

import (
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/gorm"
)

type SaleOfferRepositoryInterface interface {
	Create(offer *SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]SaleOffer, paginator.Cursor, error)
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

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter) ([]SaleOffer, paginator.Cursor, error) {
	var saleOffers []SaleOffer
	query := r.DB.
		Joins("LEFT JOIN auctions on auctions.offer_id = sale_offers.id").
		Preload("Auction").
		Preload("Car").
		Preload("Car.Model").
		Preload("Car.Model.Manufacturer")
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, paginator.Cursor{}, err
	}
	p := GetOfferPaginator(filter.PagingQuery, filter.OrderKey)
	result, cursor, err := p.Paginate(query, &saleOffers)
	if err != nil {
		return nil, paginator.Cursor{}, err
	}
	if result.Error != nil {
		return nil, paginator.Cursor{}, result.Error
	}
	return saleOffers, cursor, nil
}
