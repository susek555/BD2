package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
)

type SaleOfferServiceInterface interface {
	Create(in CreateSaleOfferDTO) error
	GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error)
	GetByID(id uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByUserID(id uint) ([]RetrieveSaleOfferDTO, error)
}

type SaleOfferService struct {
	repo    SaleOfferRepositoryInterface
	manRepo manufacturer.ManufacturerRepositoryInterface
}

func NewSaleOfferService(saleOfferRepository SaleOfferRepositoryInterface, manufacturerRepo manufacturer.ManufacturerRepositoryInterface) SaleOfferServiceInterface {
	return &SaleOfferService{repo: saleOfferRepository, manRepo: manufacturerRepo}
}

func (s *SaleOfferService) Create(in CreateSaleOfferDTO) error {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return err
	}
	return s.repo.Create(offer)
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
	offersDTOs := mapping.MapSliceToDTOs(offers, (*SaleOffer).MapToDTO)
	return &RetrieveOffersWithPagination{Offers: offersDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) GetByID(id uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	offerDTO := offer.MapToDetailedDTO()
	return offerDTO, nil
}

func (s *SaleOfferService) GetByUserID(id uint) ([]RetrieveSaleOfferDTO, error) {
	offers, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
	}
	offersDTOs := mapping.MapSliceToDTOs(offers, (*SaleOffer).MapToDTO)
	return offersDTOs, nil
}
