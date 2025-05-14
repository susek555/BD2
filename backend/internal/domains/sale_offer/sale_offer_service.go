package sale_offer

import "github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"

type SaleOfferServiceInterface interface {
	Create(in CreateSaleOfferDTO) error
	GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error)
}

type SaleOfferService struct {
	repo       SaleOfferRepositoryInterface
	manService manufacturer.ManufacturerServiceInterface
}

func NewSaleOfferService(saleOfferRepository SaleOfferRepositoryInterface, manufacturerService manufacturer.ManufacturerServiceInterface) SaleOfferServiceInterface {
	return &SaleOfferService{repo: saleOfferRepository, manService: manufacturerService}
}

func (s *SaleOfferService) Create(in CreateSaleOfferDTO) error {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return err
	}
	return s.repo.Create(offer)
}

func (s *SaleOfferService) GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error) {
	manufacturers, err := s.manService.GetAllAsNames()
	if err != nil {
		return nil, err
	}
	filter.Constriants.Manufacturers = manufacturers
	offers, pagResponse, err := s.repo.GetFiltered(filter)
	if err != nil {
		return nil, err
	}
	offersDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto := offer.MapToDTO()
		offersDTOs = append(offersDTOs, *dto)
	}
	return &RetrieveOffersWithPagination{Offers: offersDTOs, PaginationResponse: pagResponse}, nil
}
