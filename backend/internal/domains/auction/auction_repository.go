package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"gorm.io/gorm"
)

type AuctionRepositoryInterface interface {
	//generic.CRUDRepository[sale_offer.Auction]
	Create(auction *sale_offer.Auction) error
	GetAll() ([]sale_offer.Auction, error)
	GetById(id uint) (*sale_offer.Auction, error)
}

type AuctionRepository struct {
	DB *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) AuctionRepositoryInterface {
	return &AuctionRepository{DB: db}
}

func (a *AuctionRepository) Create(auction *sale_offer.Auction) error {
	return a.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(auction.Offer.Car).Error; err != nil {
			return err
		}
		if err := tx.Create(auction.Offer).Error; err != nil {
			return err
		}
		if err := tx.Create(auction).Error; err != nil {
			return err
		}
		return nil
	})
}

func (a *AuctionRepository) GetAll() ([]sale_offer.Auction, error) {
	db := a.DB
	var auctions []sale_offer.Auction
	err := db.Preload("Offer").
		Preload("Offer.Car").
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
		First(&auction, id).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}
