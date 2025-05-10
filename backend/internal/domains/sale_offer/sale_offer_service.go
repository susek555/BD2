package sale_offer

type SaleOfferServiceInterface interface {
	Create(in CreateSaleOfferDTO) error
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
