package bid

import (
	"sync"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type BidServiceInterface interface {
	Create(bidDTO *CreateBidDTO, bidderID uint) (*ProcessingBidDTO, error)
	GetAll() ([]RetrieveBidDTO, error)
	GetByID(id uint) (*RetrieveBidDTO, error)
	GetByBidderID(bidderID uint) ([]RetrieveBidDTO, error)
	GetByAuctionID(auctionID uint) ([]RetrieveBidDTO, error)
	GetHighestBid(auctionID uint) (*RetrieveBidDTO, error)
	GetHighestBidByUserID(auctionID, userID uint) (*RetrieveBidDTO, error)
}

type BidService struct {
	Repo           BidRepositoryInterface
	AuctionService auction.AuctionServiceInterface
}

func NewBidService(repo BidRepositoryInterface, auctionService auction.AuctionServiceInterface) BidServiceInterface {
	return &BidService{
		Repo:           repo,
		AuctionService: auctionService,
	}
}

var auctionLocks sync.Map

func (service *BidService) Create(bidDTO *CreateBidDTO, bidderID uint) (*ProcessingBidDTO, error) {
	bid := bidDTO.MapToBid(bidderID)
	l, _ := auctionLocks.LoadOrStore(bid.AuctionID, &sync.Mutex{})
	m := l.(*sync.Mutex)

	m.Lock()
	defer m.Unlock()
	err := service.Repo.Create(bid)
	if err != nil {
		return nil, err
	}
	err = service.AuctionService.UpdatePrice(bid.AuctionID, bid.Amount)
	if err != nil {
		return nil, err
	}
	return MapToProcessingDTO(bid), nil
}

func (service *BidService) GetAll() ([]RetrieveBidDTO, error) {
	bids, err := service.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	return mapping.MapSliceToDTOs(bids, MapToDTO), nil
}

func (service *BidService) GetByID(id uint) (*RetrieveBidDTO, error) {
	bid, err := service.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return MapToDTO(bid), nil
}

func (service *BidService) GetByBidderID(bidderID uint) ([]RetrieveBidDTO, error) {
	bids, err := service.Repo.GetByBidderID(bidderID)
	if err != nil {
		return nil, err
	}

	return mapping.MapSliceToDTOs(bids, MapToDTO), nil
}

func (service *BidService) GetByAuctionID(auctionID uint) ([]RetrieveBidDTO, error) {
	bids, err := service.Repo.GetByAuctionID(auctionID)
	if err != nil {
		return nil, err
	}

	return mapping.MapSliceToDTOs(bids, MapToDTO), nil
}

func (service *BidService) GetHighestBid(auctionID uint) (*RetrieveBidDTO, error) {
	bid, err := service.Repo.GetHighestBid(auctionID)
	if err != nil {
		return nil, err
	}
	return MapToDTO(bid), nil
}

func (service *BidService) GetHighestBidByUserID(auctionID, userID uint) (*RetrieveBidDTO, error) {
	bid, err := service.Repo.GetHighestBidByUserID(auctionID, userID)
	if err != nil {
		return nil, err
	}
	return MapToDTO(bid), nil
}
