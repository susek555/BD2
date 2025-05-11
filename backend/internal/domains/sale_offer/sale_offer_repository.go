package sale_offer

import "gorm.io/gorm"

type SaleOfferRepositoryInterface interface {
	Create(offer *SaleOffer) error
	GetFiltered(filter *OfferFilter) ([]SaleOffer, error)
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

func (r *SaleOfferRepository) GetFiltered(filter *OfferFilter) ([]SaleOffer, error) {
	var saleOffers []SaleOffer
	query := r.DB.
		Preload("Auction").
		Preload("Car").
		Preload("Car.Model").
		Preload("Car.Model.Manufacturer").
		Joins("JOIN cars ON cars.id = sale_offers.car_id").
		Joins("JOIN auctions on auctions.offer_id = sale_offers.id")
	query, err := filter.ApplyOfferFilters(query)
	if err != nil {
		return nil, err
	}
	err = query.Find(&saleOffers).Error
	if err != nil {
		return nil, err
	}
	return saleOffers, nil
}
