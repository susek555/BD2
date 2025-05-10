package auction

import "github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"

type AuctionServiceInterface interface {
	Create(auction *sale_offer.Auction) error
	GetAll() ([]sale_offer.Auction, error)
	GetByID(id uint) (*sale_offer.Auction, error)
}

type AuctionService struct {
	repo AuctionRepositoryInterface
}

func NewAuctionService(repo AuctionRepositoryInterface) AuctionServiceInterface {
	return &AuctionService{
		repo: repo,
	}
}

func (s *AuctionService) Create(auction *sale_offer.Auction) error {
	return s.repo.Create(auction)
}

func (s *AuctionService) GetAll() ([]sale_offer.Auction, error) {
	return s.repo.GetAll()
}

func (s *AuctionService) GetByID(id uint) (*sale_offer.Auction, error) {
	return s.repo.GetById(id)
}
