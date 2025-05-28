package auction

import (
	"errors"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
)

type AuctionServiceInterface interface {
	Create(auction *CreateAuctionDTO) (*RetrieveAuctionDTO, error)
	GetAll() ([]RetrieveAuctionDTO, error)
	GetById(id uint) (*RetrieveAuctionDTO, error)
	Update(auction *UpdateAuctionDTO) (*RetrieveAuctionDTO, error)
	Delete(id, userId uint) error
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
	auctionEntity, err := auction.MapToAuction()
	if err != nil {
		return nil, err
	}
	modelID, err := s.saleOfferService.GetModelID(auction.ManufacturerName, auction.ModelName)
	if err != nil {
		return nil, err
	}
	auctionEntity.Offer.Car.ModelID = modelID
	auctionEntity.Offer.Status = models.PENDING
	if auctionEntity.BuyNowPrice < 1 {
		return nil, ErrBuyNowPriceLessThan1
	}
	if auctionEntity.BuyNowPrice < auctionEntity.Offer.Price {
		return nil, ErrBuyNowPriceLessThanOfferPrice
	}
	err = s.auctionRepo.Create(auctionEntity)
	if err != nil {
		return nil, err
	}
	dto, err := s.GetById(auctionEntity.OfferID)
	if err != nil {
		return nil, err
	}
	return dto, nil
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

func (s *AuctionService) GetById(id uint) (*RetrieveAuctionDTO, error) {
	auction, err := s.auctionRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auction)
	return dto, nil
}

func (s *AuctionService) Update(auction *UpdateAuctionDTO) (*RetrieveAuctionDTO, error) {
	auctionEntity, err := auction.MapToAuction()
	if err != nil {
		return nil, err
	}
	err = s.auctionRepo.Update(auctionEntity)
	if err != nil {
		return nil, err
	}
	dto := MapToDTO(auctionEntity)
	return dto, nil
}

func (s *AuctionService) Delete(id, userId uint) error {
	auction, err := s.auctionRepo.GetById(id)
	if err != nil {
		return err
	}
	if auction.Offer.UserID != userId {
		return errors.New("you are not the owner of this auction")
	}
	err = s.auctionRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
