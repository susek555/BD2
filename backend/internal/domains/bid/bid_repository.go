package bid

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BidRepositoryInterface interface {
	Create(bid *Bid) error
	GetById(id uint) (*Bid, error)
	GetByBidderId(bidderID uint) ([]Bid, error)
	GetAll() ([]Bid, error)
	GetByAuctionId(auctionID uint) ([]Bid, error)
	GetHighestBid(auctionID uint) (*Bid, error)
	GetHighestBidByUserId(auctionID, userID uint) (*Bid, error)
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

func (b *BidRepository) Create(bid *Bid) error {
	return b.DB.Transaction(func(tx *gorm.DB) error {
		var highest Bid

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("auction_id = ?", bid.AuctionID).
			Order("amount DESC").
			Limit(1).
			Take(&highest).Error
		if err != nil {
			return err
		}
		if highest.Amount >= bid.Amount {
			return ErrBidTooLow
		}

		return tx.Create(bid).Error
	})
}

func (b *BidRepository) GetAll() ([]Bid, error) {
	db := b.DB
	var bids []Bid
	err := db.Find(&bids).Error
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetById(id uint) (*Bid, error) {
	db := b.DB
	var bid Bid
	if err := db.First(&bid, id).Error; err != nil {
		return nil, err
	}
	return &bid, nil
}

func (b *BidRepository) GetByBidderId(bidderID uint) ([]Bid, error) {
	db := b.DB
	var bids []Bid
	if err := db.Where("bidder_id = ?", bidderID).Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetByAuctionId(auctionID uint) ([]Bid, error) {
	db := b.DB
	var bids []Bid
	if err := db.Where("auction_id = ?", auctionID).Find(&bids).Error; err != nil {
		return nil, err
	}
	return bids, nil
}

func (b *BidRepository) GetHighestBid(auctionID uint) (*Bid, error) {
	db := b.DB
	var bid Bid
	err := db.
		Where("auction_id = ?", auctionID).
		Order("amount desc").
		First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (b *BidRepository) GetHighestBidByUserId(auctionID, userID uint) (*Bid, error) {
	db := b.DB
	var bid Bid
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
