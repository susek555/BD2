package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

//go:generate mockery --name=AuctionServiceInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionServiceInterface interface {
	Create(auction *CreateAuctionDTO) (*RetrieveAuctionDTO, error)
	Update(auction *UpdateAuctionDTO, userID uint) (*RetrieveAuctionDTO, error)
	GetByID(id uint, userID *uint) (*RetrieveAuctionDTO, error)
	BuyNow(auctionID, userID uint) (*models.Auction, error)
	UpdatePrice(auctionID uint, newPrice uint) error
	GetByIDNonDTO(id uint) (*models.Auction, error)
	Delete(id, userID uint) error
}

type AuctionService struct {
	auctionRepo      AuctionRepositoryInterface
	saleOfferService sale_offer.SaleOfferServiceInterface
}

func NewAuctionService(repo AuctionRepositoryInterface, service sale_offer.SaleOfferServiceInterface) AuctionServiceInterface {
	return &AuctionService{
		auctionRepo:      repo,
		saleOfferService: service,
	}
}

func (s *AuctionService) Create(auction *CreateAuctionDTO) (*RetrieveAuctionDTO, error) {
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
	return s.GetByID(offer.ID, &offer.UserID)
}

func (s *AuctionService) Update(auction *UpdateAuctionDTO, userID uint) (*RetrieveAuctionDTO, error) {
	if auction.UpdateSaleOfferDTO != nil {
		if _, err := s.saleOfferService.Update(auction.UpdateSaleOfferDTO, userID); err != nil {
			return nil, err
		}
	}
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
	return s.GetByID(updatedAuction.OfferID, &userID)
}

func (s *AuctionService) GetByID(id uint, userID *uint) (*RetrieveAuctionDTO, error) {
	offer, err := s.saleOfferService.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	auction, err := s.auctionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auction)
	dto.RetrieveSaleOfferDTO = offer
	return dto, nil
}

func (s *AuctionService) BuyNow(auctionID, userID uint) (*models.Auction, error) {
	auction, err := s.auctionRepo.GetByID(auctionID)
	if err != nil {
		return nil, err
	}
	if auction.Offer.UserID == userID {
		return nil, ErrAuctionOwnedByUser
	}
	if auction.Offer.Status != enums.PUBLISHED {
		return nil, ErrAuctionNotPublished
	}
	if err := s.auctionRepo.BuyNow(auction, userID); err != nil {
		return nil, err
	}
	_ = s.UpdatePrice(auctionID, auction.BuyNowPrice)
	return auction, err
}

func (s *AuctionService) UpdatePrice(auctionID uint, newPrice uint) error {
	auction, err := s.auctionRepo.GetByID(auctionID)
	if err != nil {
		return err
	}
	if newPrice < auction.Offer.Price || newPrice < auction.InitialPrice {
		return ErrNewPriceLessThanOfferPrice
	}
	err = s.auctionRepo.UpdatePrice(auctionID, newPrice)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuctionService) GetByIDNonDTO(id uint) (*models.Auction, error) {
	auction, err := s.auctionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return auction, nil
}

func (s *AuctionService) Delete(id, userID uint) error {
	auction, err := s.auctionRepo.GetByID(id)
	if err != nil {
		return err
	}
	if !auction.BelongsToUser(userID) {
		return ErrAuctionNotOwned
	}
	return s.auctionRepo.Delete(id)
}
