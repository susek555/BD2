package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"gorm.io/gorm"
)

//go:generate mockery --name=AuctionRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionRepositoryInterface interface {
	generic.CRUDRepository[sale_offer.Auction]
}

type AuctionRepository struct {
	DB *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) AuctionRepositoryInterface {
	return &AuctionRepository{DB: db}
}

func (a *AuctionRepository) Create(auction *sale_offer.Auction) error {
	return a.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Transaction(func(tx *gorm.DB) error {
			err := tx.Create(auction).Error
			if err != nil {
				return err
			}
			return tx.
				Preload("Offer.Car.Model.Manufacturer").
				First(auction, "offer_id = ?", auction.OfferID).Error
		})
}

func (a *AuctionRepository) GetAll() ([]sale_offer.Auction, error) {
	db := a.DB
	var auctions []sale_offer.Auction
	err := db.Preload("Offer").
		Preload("Offer.Car").
		Preload("Offer.Car.Model").
		Preload("Offer.Car.Model.Manufacturer").
		Find(&auctions).Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (a *AuctionRepository) GetById(id uint) (*sale_offer.Auction, error) {
	db := a.DB
	var auction sale_offer.Auction
	err := db.Preload("Offer").
		Preload("Offer.Car").
		Preload("Offer.Car.Model").
		Preload("Offer.Car.Model.Manufacturer").
		First(&auction, id).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

func (a *AuctionRepository) Update(auction *sale_offer.Auction) error {
	return a.DB.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(auction).Error; err != nil {
				return err
			}
			return tx.
				Preload("Offer.Car.Model.Manufacturer").
				First(auction, "offer_id = ?", auction.OfferID).Error
		})
}

func (a *AuctionRepository) Delete(id uint) error {
	return a.DB.
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&sale_offer.Car{}, id).Error; err != nil {
				return err
			}
			return nil
		})
}
