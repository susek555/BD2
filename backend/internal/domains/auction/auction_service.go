package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

type SaleOfferServiceInterface interface {
	Create(in *sale_offer.CreateSaleOfferDTO) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	Update(in *sale_offer.UpdateSaleOfferDTO, userID uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	GetDetailedByID(id uint, userID *uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	Delete(id uint, userID uint) error
}

//go:generate mockery --name=AuctionServiceInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionServiceInterface interface {
	Create(auction *CreateAuctionDTO) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	Update(auction *UpdateAuctionDTO, userID uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	BuyNow(auctionID, userID uint) (*models.Auction, error)
	UpdatePrice(auction *models.Auction, newPrice uint) error
	Delete(id, userID uint) error
}

type AuctionService struct {
	auctionRepo      AuctionRepositoryInterface
	saleOfferService SaleOfferServiceInterface
}

func NewAuctionService(repo AuctionRepositoryInterface, service SaleOfferServiceInterface) AuctionServiceInterface {
	return &AuctionService{
		auctionRepo:      repo,
		saleOfferService: service,
	}
}

func (s *AuctionService) Create(auction *CreateAuctionDTO) (*sale_offer.RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferService.Create(&auction.CreateSaleOfferDTO)
	if err != nil {
		return nil, err
	}
	auctionEntity, err := auction.MapToAuction()
	if err != nil {
		return nil, err
	}
	if auctionEntity.BuyNowPrice < 1 {
		return nil, ErrBuyNowPriceLessThan1
	}
	if auctionEntity.BuyNowPrice < offer.Price {
		return nil, ErrBuyNowPriceLessThanOfferPrice
	}
	auctionEntity.OfferID = offer.ID
	if err := s.auctionRepo.Create(auctionEntity); err != nil {
		return nil, err
	}
	return s.saleOfferService.GetDetailedByID(offer.ID, &offer.UserID)
}

func (s *AuctionService) Update(auction *UpdateAuctionDTO, userID uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error) {
	auctionEntity, err := s.auctionRepo.GetByID(auction.ID)
	if err != nil {
		return nil, err
	}
	updatedAuction, err := auction.UpdatedAuctionFromDTO(auctionEntity)
	if err != nil {
		return nil, err
	}
	err = s.auctionRepo.Update(updatedAuction)
	if err != nil {
		return nil, err
	}
	return s.saleOfferService.Update(&auction.UpdateSaleOfferDTO, userID)
}

func (s *AuctionService) BuyNow(auctionID, userID uint) (*models.Auction, error) {
	auction, err := s.auctionRepo.GetByID(auctionID)
	if err != nil {
		return nil, err
	}
	if !auction.Offer.BelongsToUser(userID) {
		return nil, ErrAuctionOwnedByUser
	}
	if auction.Offer.Status != enums.PUBLISHED {
		return nil, ErrAuctionNotPublished
	}
	if err := s.auctionRepo.BuyNow(auction, userID); err != nil {
		return nil, err
	}
	_ = s.UpdatePrice(auction, auction.BuyNowPrice)
	return auction, err
}

func (s *AuctionService) UpdatePrice(auction *models.Auction, newPrice uint) error {
	if newPrice < auction.Offer.Price || newPrice < auction.InitialPrice {
		return ErrNewPriceLessThanOfferPrice
	}
	auction.Offer.Price = newPrice
	if err := s.auctionRepo.Update(auction); err != nil {
		return err
	}
	return nil
}

func (s *AuctionService) Delete(id uint, userID uint) error {
	return s.saleOfferService.Delete(id, userID)
}
