package bid

import "gorm.io/gorm"

type BidRepositoryInterface interface {
	Create(bid *Bid) error
	//GetByID(id uint) (*Bid, error)
	//GetByBidderID(bidderID uint) ([]Bid, error)
	//GetAll() ([]Bid, error)
	//GetByAuctionID(auctionID uint) ([]Bid, error)
	//GetHighestBid(auctionID uint) (*Bid, error)
	//GetHighestBidByUserId(userID uint) (*Bid, error)
}

type BidRepository struct {
	DB *gorm.DB
}

func NewBidRepository(db *gorm.DB) BidRepositoryInterface {
	return &BidRepository{
		DB: db,
	}
}

func (b *BidRepository) Create(bid *Bid) error {
	if err := b.DB.Create(bid).Error; err != nil {
		return err
	}
	return nil
}
