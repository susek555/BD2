package sale_offer

type SaleOfferServiceInterface interface {
	Create(in CreateSaleOfferDTO) error
	GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error)
}

type SaleOfferService struct {
	repo SaleOfferRepositoryInterface
}

func NewSaleOfferService(saleOfferRepository SaleOfferRepositoryInterface) SaleOfferServiceInterface {
	return &SaleOfferService{repo: saleOfferRepository}
}

func (s *SaleOfferService) Create(in CreateSaleOfferDTO) error {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return err
	}
	return s.repo.Create(offer)
}

func (s *SaleOfferService) GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error) {
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
