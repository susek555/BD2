package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=AuctionRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionRepositoryInterface interface {
	generic.CRUDRepository[models.Auction]
	BuyNow(auctionID, userID uint) (*models.Auction, error)
	UpdatePrice(auctionID uint, newPrice uint) error
}

type AuctionRepository struct {
	DB *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) AuctionRepositoryInterface {
	return &AuctionRepository{DB: db}
}

func (a *AuctionRepository) Create(auction *models.Auction) error {
	return a.DB.Create(auction).Error
}

func (a *AuctionRepository) GetAll() ([]models.Auction, error) {
	db := a.DB
	var auctions []models.Auction
	err := db.Preload("Offer").
		Preload("Offer.Car").
		Preload("Offer.Car.Model").
		Preload("Offer.Car.Model.Manufacturer").
		Preload("Offer.User").
		Find(&auctions).Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (a *AuctionRepository) GetByID(id uint) (*models.Auction, error) {
	db := a.DB
	var auction models.Auction
	err := db.Preload("Offer").
		Preload("Offer.Car").
		Preload("Offer.Car.Model").
		Preload("Offer.Car.Model.Manufacturer").
		Preload("Offer.User").
		First(&auction, id).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

func (a *AuctionRepository) Update(auction *models.Auction) error {
	return a.DB.Save(auction).Error
}

func (a *AuctionRepository) Delete(id uint) error {
	return a.DB.
		Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&models.Car{}, id).Error; err != nil {
				return err
			}
			return nil
		})
}

func (a *AuctionRepository) BuyNow(auctionID, userID uint) (*models.Auction, error) {
	auction, err := a.GetByID(auctionID)
	if err != nil {
		return nil, err
	}
	if auction.Offer.Status == enums.SOLD {
		return nil, ErrAuctionAlreadySold
	}
	err = a.DB.Model(&models.SaleOffer{}).
		Where("id = ?", auctionID).
		Update("status", enums.SOLD).
		Error
	if err != nil {
		return nil, err
	}
	err = a.DB.Model(&models.Purchase{}).
		Create(&models.Purchase{
			OfferID:    auction.OfferID,
			BuyerID:    userID,
			FinalPrice: auction.BuyNowPrice,
			IssueDate:  auction.DateEnd,
		}).Error
	return auction, err
}

func (a *AuctionRepository) UpdatePrice(auctionID uint, newPrice uint) error {
	return a.DB.Model(&models.SaleOffer{}).
		Where("id = ?", auctionID).
		Update("price", newPrice).
		Error
}
