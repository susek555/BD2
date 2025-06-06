package bid

import (
	"errors"

	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=BidRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type BidRepositoryInterface interface {
	Create(bid *models.Bid) error
	GetByID(id uint) (*models.Bid, error)
	GetByBidderID(bidderID uint) ([]models.Bid, error)
	GetAll() ([]models.Bid, error)
	GetByAuctionID(auctionID uint) ([]models.Bid, error)
	GetHighestBid(auctionID uint) (*models.Bid, error)
	GetHighestBidByUserID(auctionID, userID uint) (*models.Bid, error)
}

type BidRepository struct {
	DB *gorm.DB
}

var (
	ErrBidTooLow = errors.New("bid is lower than current highest")
)

func NewBidRepository(db *gorm.DB) BidRepositoryInterface {
	return &BidRepository{
		DB: db,
	}
}

func (b *BidRepository) Create(bid *models.Bid) error {
	return b.DB.Transaction(func(tx *gorm.DB) error {
		var highest models.Bid
		var auction models.Auction
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("auction_id = ?", bid.AuctionID).
			Order("amount DESC").
			Preload("Auction").
			Preload("Auction.Offer").
			Limit(1).
			Take(&highest).Error
		if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
			return err
		}
		if highest.Amount >= bid.Amount {
			return ErrBidTooLow
		}
		err = tx.
			Model(&auction).
			Where("offer_id = ?", bid.AuctionID).
			Preload("Offer").
			First(&auction).Error
		if err != nil {
			return err
		}
		if auction.Offer.Price > bid.Amount {
			return ErrBidTooLow
		}

		err = tx.Create(bid).Error
		if err != nil {
			return err
		}
		if err := tx.
			Preload("Auction").
			Preload("Auction.Offer.Car.Model.Manufacturer").
			Preload("Bidder").
			First(bid, bid.ID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (b *BidRepository) GetAll() ([]models.Bid, error) {
	db := b.DB
	var bids []models.Bid
	err := db.Find(&bids).Error
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetByID(id uint) (*models.Bid, error) {
	db := b.DB
	var bid models.Bid
	if err := db.First(&bid, id).Error; err != nil {
		return nil, err
	}
	return &bid, nil
}

func (b *BidRepository) GetByBidderID(bidderID uint) ([]models.Bid, error) {
	db := b.DB
	var bids []models.Bid
	if err := db.Where("bidder_id = ?", bidderID).Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetByAuctionID(auctionID uint) ([]models.Bid, error) {
	db := b.DB
	var bids []models.Bid
	if err := db.Where("auction_id = ?", auctionID).Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetHighestBid(auctionID uint) (*models.Bid, error) {
	db := b.DB
	var bid models.Bid
	err := db.
		Where("auction_id = ?", auctionID).
		Order("amount desc").
		First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (b *BidRepository) GetHighestBidByUserID(auctionID, userID uint) (*models.Bid, error) {
	db := b.DB
	var bid models.Bid
	err := db.
		Where("auction_id = ?", auctionID).
		Where("bidder_id = ?", userID).
		Order("amount desc").
		First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
