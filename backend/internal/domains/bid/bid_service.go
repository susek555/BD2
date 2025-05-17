package bid

import (
	"sync"

	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type BidServiceInterface interface {
	Create(bidDTO *CreateBidDTO, bidderID uint) (*RetrieveBidDTO, error)
	GetAll() ([]RetrieveBidDTO, error)
	GetById(id uint) (*RetrieveBidDTO, error)
	GetByBidderId(bidderId uint) ([]RetrieveBidDTO, error)
	GetByAuctionId(auctionID uint) ([]RetrieveBidDTO, error)
	GetHighestBid(auctionId uint) (*RetrieveBidDTO, error)
	GetHighestBidByUserId(auctionId, userId uint) (*RetrieveBidDTO, error)
}

type BidService struct {
	Repo BidRepositoryInterface
}

func NewBidService(repo BidRepositoryInterface) BidServiceInterface {
	return &BidService{
		Repo: repo,
	}
}

var auctionLocks sync.Map

func (service *BidService) Create(bidDTO *CreateBidDTO, bidderID uint) (*RetrieveBidDTO, error) {
	bid := bidDTO.MapToBid(bidderID)
	l, _ := auctionLocks.LoadOrStore(bid.AuctionID, &sync.Mutex{})
	m := l.(*sync.Mutex)

	m.Lock()
	defer m.Unlock()

	err := service.Repo.Create(bid)
	if err != nil {
		return nil, err
	}
	return bid.MapToDTO(), nil
}

func (service *BidService) GetAll() ([]RetrieveBidDTO, error) {
	bids, err := service.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	return mapping.MapSliceToDTOs(bids, (*Bid).MapToDTO), nil
}

func (service *BidService) GetById(id uint) (*RetrieveBidDTO, error) {
	bid, err := service.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return bid.MapToDTO(), nil
}

func (service *BidService) GetByBidderId(bidderId uint) ([]RetrieveBidDTO, error) {
	bids, err := service.Repo.GetByBidderId(bidderId)
	if err != nil {
		return nil, err
	}

	return mapping.MapSliceToDTOs(bids, (*Bid).MapToDTO), nil
}

func (service *BidService) GetByAuctionId(auctionID uint) ([]RetrieveBidDTO, error) {
	bids, err := service.Repo.GetByAuctionId(auctionID)
	if err != nil {
		return nil, err
	}

	return mapping.MapSliceToDTOs(bids, (*Bid).MapToDTO), nil
}

func (service *BidService) GetHighestBid(auctionId uint) (*RetrieveBidDTO, error) {
	bid, err := service.Repo.GetHighestBid(auctionId)
	if err != nil {
		return nil, err
	}
	return bid.MapToDTO(), nil
}

func (service *BidService) GetHighestBidByUserId(auctionId, userId uint) (*RetrieveBidDTO, error) {
	bid, err := service.Repo.GetHighestBidByUserId(auctionId, userId)
	if err != nil {
		return nil, err
	}
	return bid.MapToDTO(), nil
}
