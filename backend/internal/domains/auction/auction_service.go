package auction

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

//go:generate mockery --name=PurchaseCreatorInterface --output=../../test/mocks --case=snake --with-expecter
type PurchaseCreatorInterface interface {
	Create(purchase *models.Purchase) error
}

//go:generate mockery --name=SaleOfferServiceInterface --output=../../test/mocks --case=snake --with-expecter
type SaleOfferServiceInterface interface {
	PrepareForCreateSaleOffer(in *sale_offer.CreateSaleOfferDTO) (*models.SaleOffer, error)
	PrepareForUpdateSaleOffer(in *sale_offer.UpdateSaleOfferDTO, userID uint) (*models.SaleOffer, error)
	PrepareForBuySaleOffer(id uint, userID uint) (*models.SaleOffer, error)
	GetDetailedByID(id uint, userID *uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	Delete(id uint, userID uint) error
}

//go:generate mockery --name=AuctionServiceInterface --output=../../test/mocks --case=snake --with-expecter
type AuctionServiceInterface interface {
	Create(auction *CreateAuctionDTO) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	Update(auction *UpdateAuctionDTO, userID uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error)
	BuyNow(auctionID, userID uint) (*models.SaleOffer, error)
	UpdatePrice(auction *models.SaleOffer, newPrice uint) error
	Delete(id, userID uint) error
}

type AuctionService struct {
	saleOfferRepo    sale_offer.SaleOfferRepositoryInterface
	saleOfferService SaleOfferServiceInterface
	purchaseCreator  PurchaseCreatorInterface
}

func NewAuctionService(repo sale_offer.SaleOfferRepositoryInterface, service SaleOfferServiceInterface, purchaseCreator PurchaseCreatorInterface) AuctionServiceInterface {
	return &AuctionService{
		saleOfferRepo:    repo,
		saleOfferService: service,
		purchaseCreator:  purchaseCreator,
	}
}

func (s *AuctionService) Create(in *CreateAuctionDTO) (*sale_offer.RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferService.PrepareForCreateSaleOffer(&in.CreateSaleOfferDTO)
	if err != nil {
		return nil, err
	}
	auction, err := s.PrepareForCreateAuction(in)
	if err != nil {
		return nil, err
	}
	offer.Auction = auction
	if err := s.saleOfferRepo.Create(offer); err != nil {
		return nil, err
	}
	return s.saleOfferService.GetDetailedByID(offer.ID, &offer.UserID)
}

func (s *AuctionService) Update(in *UpdateAuctionDTO, userID uint) (*sale_offer.RetrieveDetailedSaleOfferDTO, error) {
	updatedOffer, err := s.saleOfferService.PrepareForUpdateSaleOffer(&in.UpdateSaleOfferDTO, userID)
	if err != nil {
		return nil, err
	}
	updatedAuction, err := s.PrepareForUpdateAuction(in, updatedOffer)
	if err != nil {
		return nil, err
	}
	updatedOffer = updatedAuction
	if err := s.saleOfferRepo.Update(updatedOffer); err != nil {
		return nil, err
	}
	return s.saleOfferService.GetDetailedByID(updatedOffer.ID, &updatedOffer.UserID)
}

func (s *AuctionService) BuyNow(id uint, userID uint) (*models.SaleOffer, error) {
	offer, err := s.saleOfferService.PrepareForBuySaleOffer(id, userID)
	if err != nil {
		return nil, err
	}
	offer.Status = enums.SOLD
	offer.Price = offer.Auction.BuyNowPrice
	if err := s.saleOfferRepo.Update(offer); err != nil {
		return nil, err
	}
	purchaseModel := &models.Purchase{OfferID: offer.ID, BuyerID: userID, FinalPrice: offer.Auction.BuyNowPrice, IssueDate: time.Now()}
	if err := s.purchaseCreator.Create(purchaseModel); err != nil {
		return nil, err
	}
	return offer, err
}

func (s *AuctionService) UpdatePrice(offer *models.SaleOffer, newPrice uint) error {
	if newPrice < offer.Price || newPrice < offer.Auction.InitialPrice {
		return ErrNewPriceLessThanOfferPrice
	}
	offer.Price = newPrice
	if err := s.saleOfferRepo.Update(offer); err != nil {
		return err
	}
	return nil
}

func (s *AuctionService) Delete(id uint, userID uint) error {
	return s.saleOfferService.Delete(id, userID)
}

func (s *AuctionService) PrepareForCreateAuction(in *CreateAuctionDTO) (*models.Auction, error) {
	auction, err := in.MapToAuction()
	if err != nil {
		return nil, err
	}
	if auction.BuyNowPrice < 1 {
		return nil, ErrBuyNowPriceLessThan1
	}
	if auction.BuyNowPrice < in.CreateSaleOfferDTO.Price {
		return nil, ErrBuyNowPriceLessThanOfferPrice
	}
	return auction, nil
}

func (s *AuctionService) PrepareForUpdateAuction(in *UpdateAuctionDTO, auction *models.SaleOffer) (*models.SaleOffer, error) {
	updatedAuction, err := in.UpdatedAuctionFromDTO(auction)
	if err != nil {
		return nil, err
	}
	return updatedAuction, nil
}
