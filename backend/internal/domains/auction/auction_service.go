package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
)

type AuctionServiceInterface interface {
	generic.CRUDService[sale_offer.Auction]
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

func (s *AuctionService) GetById(id uint) (*sale_offer.Auction, error) {
	return s.repo.GetById(id)
}

func (s *AuctionService) Update(auction *sale_offer.Auction) error {
	return s.repo.Update(auction)
}

func (s *AuctionService) Delete(id uint) error {
	return s.repo.Delete(id)
}
