package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type SaleOfferServiceInterface interface {
	Create(in CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error)
	GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error)
	GetByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error)
	LikeOffer(offerID, userID uint) error
	DislikeOffer(offerID, userID uint) error
	IsOfferLikedByUser(offerID, userID uint) bool
	CanBeModifiedByUser(offerID, userID uint) (bool, error)
}

type SaleOfferService struct {
	repo           SaleOfferRepositoryInterface
	manRepo        manufacturer.ManufacturerRepositoryInterface
	likedOfferRepo liked_offer.LikedOfferRepositoryInterface
	bidRepo        bid.BidRepositoryInterface
}

func NewSaleOfferService(saleOfferRepository SaleOfferRepositoryInterface, manufacturerRepo manufacturer.ManufacturerRepositoryInterface, likedOfferRepo liked_offer.LikedOfferRepositoryInterface, bidRepo bid.BidRepositoryInterface) SaleOfferServiceInterface {
	return &SaleOfferService{repo: saleOfferRepository, manRepo: manufacturerRepo, likedOfferRepo: likedOfferRepo, bidRepo: bidRepo}
}

func (s *SaleOfferService) Create(in CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(offer); err != nil {
		return nil, err
	}
	return s.GetByID(offer.ID, &offer.UserID)
}

func (s *SaleOfferService) GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error) {
	manufacturers, err := s.manRepo.GetAll()
	if err != nil {
		return nil, err
	}
	filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(manufacturers, (*manufacturer.Manufacturer).MapToName)
	offers, pagResponse, err := s.repo.GetFiltered(filter)
	if err != nil {
		return nil, err
	}
	offerDTOs, err := s.mapOfferSliceWithAdditionalFields(offers, filter.UserID)
	if err != nil {
		return nil, err
	}
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) GetByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	offerDTO := offer.MapToDetailedDTO()
	s.fillUserFields(&offerDTO.UserContext, offer.ID, userID)
	return offerDTO, nil
}

func (s *SaleOfferService) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error) {
	offers, pagResponse, err := s.repo.GetByUserID(id, pagRequest)
	if err != nil {
		return nil, err
	}
	offerDTOs, err := s.mapOfferSliceWithAdditionalFields(offers, &id)
	if err != nil {
		return nil, err
	}
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) LikeOffer(offerID, userID uint) error {
	offer, err := s.repo.GetByID(offerID)
	if err != nil {
		return err
	}
	if offer.UserID == userID {
		return ErrLikeOwnOffer
	}
	return s.likedOfferRepo.Create(&liked_offer.LikedOffer{OfferID: offerID, UserID: userID})
}

func (s *SaleOfferService) DislikeOffer(offerID, userID uint) error {
	if !s.IsOfferLikedByUser(offerID, userID) {
		return ErrDislikeNotLikedOffer
	}
	return s.likedOfferRepo.Delete(offerID, userID)
}

func (s *SaleOfferService) IsOfferLikedByUser(offerID, userID uint) bool {
	return s.likedOfferRepo.IsOfferLikedByUser(userID, offerID)
}

func (s *SaleOfferService) CanBeModifiedByUser(offerID, userID uint) (bool, error) {
	offer, err := s.repo.GetByID(offerID)
	if err != nil {
		return false, err
	}
	if offer.Auction == nil {
		return true, nil
	}
	bids, err := s.bidRepo.GetByAuctionId(offerID)
	if err != nil {
		return false, err
	}
	return len(bids) == 0, nil
}

func (s *SaleOfferService) fillUserFields(userContext *UserContext, offerID uint, userID *uint) error {
	if userID == nil {
		userContext.IsLiked = false
		userContext.CanModify = false
	} else {
		userContext.IsLiked = s.IsOfferLikedByUser(offerID, *userID)
		smt, err := s.CanBeModifiedByUser(offerID, *userID)
		if err != nil {
			return err
		}
		userContext.CanModify = smt
	}
	return nil
}

func (s *SaleOfferService) mapOfferSliceWithAdditionalFields(offers []SaleOffer, userID *uint) ([]RetrieveSaleOfferDTO, error) {
	offerDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto := offer.MapToDTO()
		if err := s.fillUserFields(&dto.UserContext, offer.ID, userID); err != nil {
			return nil, err
		}
		offerDTOs = append(offerDTOs, *dto)
	}
	return offerDTOs, nil
}
