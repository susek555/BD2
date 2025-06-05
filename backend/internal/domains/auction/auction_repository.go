package auction

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
	"gorm.io/gorm"
)

//go:generate mockery --name=AuctionRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionRepositoryInterface interface {
	Create(auction *models.Auction) error
	Update(auction *models.Auction) error
	BuyNow(auction *models.Auction, userID uint) error
	GetByID(id uint) (*models.Auction, error)
	GetAllActiveAuctions() ([]views.SaleOfferView, error)
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

func (a *AuctionRepository) GetByID(id uint) (*models.Auction, error) {
	var auction models.Auction
	err := a.DB.Preload("Offer").Preload("Offer.Car").First(&auction, id).Error
	return &auction, err
}

func (a *AuctionRepository) BuyNow(auction *models.Auction, userID uint) error {
	auction.Offer.Status = enums.SOLD
	if err := a.Update(auction); err != nil {
		return err
	}
	return a.SaveToPurchases(auction.OfferID, userID, auction.BuyNowPrice)
}

func (a *AuctionRepository) GetAllActiveAuctions() ([]views.SaleOfferView, error) {
	var auctions []views.SaleOfferView
	err := a.DB.Table("sale_offer_view").Where("is_auction IS TRUE").Find(&auctions).Error
	if err != nil {
		return nil, err
	}
	return auctions, nil
}

func (a *AuctionRepository) Update(auction *models.Auction) error {
	return a.DB.Save(auction).Error
}

func (a *AuctionRepository) SaveToPurchases(offerID uint, buyerID uint, finalPrice uint) error {
	purchase := models.Purchase{
		OfferID:    offerID,
		BuyerID:    buyerID,
		FinalPrice: finalPrice,
		IssueDate:  time.Now(),
	}
	return a.DB.Create(&purchase).Error
}
