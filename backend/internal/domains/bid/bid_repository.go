package bid

import (
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=BidRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type BidRepositoryInterface interface {
	Create(bid *models.Bid) error
	GetById(id uint) (*models.Bid, error)
	GetByBidderId(bidderID uint) ([]models.Bid, error)
	GetAll() ([]models.Bid, error)
	GetByAuctionId(auctionID uint) ([]models.Bid, error)
	GetHighestBid(auctionID uint) (*models.Bid, error)
	GetHighestBidByUserId(auctionID, userID uint) (*models.Bid, error)
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

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("auction_id = ?", bid.AuctionID).
			Order("amount DESC").
			Limit(1).
			Take(&highest).Error
		if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
			return err
		}
		if highest.Amount >= bid.Amount {
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

func (b *BidRepository) GetById(id uint) (*models.Bid, error) {
	db := b.DB
	var bid models.Bid
	if err := db.First(&bid, id).Error; err != nil {
		return nil, err
	}
	return &bid, nil
}

func (b *BidRepository) GetByBidderId(bidderID uint) ([]models.Bid, error) {
	db := b.DB
	var bids []models.Bid
	if err := db.Where("bidder_id = ?", bidderID).Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetByAuctionId(auctionID uint) ([]models.Bid, error) {
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

func (b *BidRepository) GetHighestBidByUserId(auctionID, userID uint) (*models.Bid, error) {
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
