package auction

import (
	"errors"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

//go:generate mockery --name=AuctionServiceInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionServiceInterface interface {
	Create(auction *CreateAuctionDTO) (*RetrieveAuctionDTO, error)
	GetAll() ([]RetrieveAuctionDTO, error)
	GetByID(id uint) (*RetrieveAuctionDTO, error)
	Update(auction *UpdateAuctionDTO, userID uint) (*RetrieveAuctionDTO, error)
	Delete(id, userID uint) error
	BuyNow(auctionID, userID uint) (*models.Auction, error)
	UpdatePrice(auctionID uint, newPrice uint) error
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
	return s.GetByID(auctionEntity.OfferID)
}

func (s *AuctionService) GetAll() ([]RetrieveAuctionDTO, error) {
	auctions, err := s.auctionRepo.GetAll()
	if err != nil {
		return nil, err
	}
	var auctionsDTO []RetrieveAuctionDTO
	for _, auction := range auctions {
		dto := MapToDTO(&auction)
		auctionsDTO = append(auctionsDTO, *dto)
	}
	return auctionsDTO, nil
}

func (s *AuctionService) GetByID(id uint) (*RetrieveAuctionDTO, error) {
	auction, err := s.auctionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auction)
	return dto, nil
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
	return s.GetByID(updatedAuction.OfferID)
}

func (s *AuctionService) Delete(ID, userID uint) error {
	auction, err := s.auctionRepo.GetByID(ID)
	if err != nil {
		return err
	}
	if auction.Offer.UserID != userID {
		return errors.New("you are not the owner of this auction")
	}
	err = s.auctionRepo.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuctionService) BuyNow(auctionID, userID uint) (*models.Auction, error) {
	auction, err := s.auctionRepo.GetByID(auctionID)
	if err != nil {
		return nil, err
	}
	if auction.Offer.UserID == userID {
		return nil, ErrAuctionOwnedByUser
	}
	auction, err = s.auctionRepo.BuyNow(auctionID, userID)
	s.UpdatePrice(auctionID, auction.BuyNowPrice)
	return auction, err
}

func (s *AuctionService) UpdatePrice(auctionID uint, newPrice uint) error {
	auction, err := s.auctionRepo.GetById(auctionID)
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
