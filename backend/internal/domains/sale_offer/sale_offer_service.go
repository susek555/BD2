package sale_offer

import (
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
	IsOfferLikedByUser(userID, offerID uint) bool
}

type SaleOfferService struct {
	repo           SaleOfferRepositoryInterface
	manRepo        manufacturer.ManufacturerRepositoryInterface
	likedOfferRepo liked_offer.LikedOfferReposisotryInterface
}

func NewSaleOfferService(saleOfferRepository SaleOfferRepositoryInterface, manufacturerRepo manufacturer.ManufacturerRepositoryInterface, likedOfferRepo liked_offer.LikedOfferReposisotryInterface) SaleOfferServiceInterface {
	return &SaleOfferService{repo: saleOfferRepository, manRepo: manufacturerRepo, likedOfferRepo: likedOfferRepo}
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
	filter.Constriants.Manufacturers = mapping.MapSliceToDTOs(manufacturers, (*manufacturer.Manufacturer).MapToName)
	offers, pagResponse, err := s.repo.GetFiltered(filter)
	if err != nil {
		return nil, err
	}
	offerDTOs := s.mapSliceWithIsLikedField(offers, filter.UserID)
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) GetByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	offerDTO := offer.MapToDetailedDTO()
	if userID == nil {
		offerDTO.IsLiked = false
	} else {
		offerDTO.IsLiked = s.IsOfferLikedByUser(*userID, offer.ID)
	}
	return offerDTO, nil
}

func (s *SaleOfferService) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error) {
	offers, pagResponse, err := s.repo.GetByUserID(id, pagRequest)
	if err != nil {
		return nil, err
	}
	offerDTOs := s.mapSliceWithIsLikedField(offers, &id)
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) mapSliceWithIsLikedField(offers []SaleOffer, userID *uint) []RetrieveSaleOfferDTO {
	offerDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto := offer.MapToDTO()
		if userID == nil {
			dto.IsLiked = false
		} else {
			dto.IsLiked = s.IsOfferLikedByUser(*userID, offer.ID)
		}
		offerDTOs = append(offerDTOs, *dto)
	}
	return offerDTOs
}

func (s *SaleOfferService) IsOfferLikedByUser(offerID, userID uint) bool {
	return s.likedOfferRepo.IsOfferLikedByUser(userID, offerID)
}
