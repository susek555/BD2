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
	GetByID(id uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error)
}

type SaleOfferService struct {
	repo           SaleOfferRepositoryInterface
	manRepo        manufacturer.ManufacturerRepositoryInterface
	likedOfferRepo liked_offer.LikedOfferReposisotry
}

func NewSaleOfferService(saleOfferRepository SaleOfferRepositoryInterface, manufacturerRepo manufacturer.ManufacturerRepositoryInterface) SaleOfferServiceInterface {
	return &SaleOfferService{repo: saleOfferRepository, manRepo: manufacturerRepo}
}

func (s *SaleOfferService) Create(in CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(offer); err != nil {
		return nil, err
	}
	return s.GetByID(offer.ID)
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
	offerDTOs := s.mapSliceWithIsLikedField(offers)
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) GetByID(id uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	offerDTO := offer.MapToDetailedDTO()
	offerDTO.IsLiked = s.likedOfferRepo.IsOfferLikedByUser(offer.ID, offer.UserID)
	return offerDTO, nil
}

func (s *SaleOfferService) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error) {
	offers, pagResponse, err := s.repo.GetByUserID(id, pagRequest)
	if err != nil {
		return nil, err
	}
	offerDTOs := s.mapSliceWithIsLikedField(offers)
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) mapSliceWithIsLikedField(offers []SaleOffer) []RetrieveSaleOfferDTO {
	offerDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto := offer.MapToDTO()
		dto.IsLiked = s.likedOfferRepo.IsOfferLikedByUser(offer.ID, offer.UserID)
		offerDTOs = append(offerDTOs, *dto)
	}
	return offerDTOs
}
