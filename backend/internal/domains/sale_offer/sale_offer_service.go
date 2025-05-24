package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type ImageRetrieverInterface interface {
	GetImagesByOfferID(offerID uint) ([]models.Image, error)
}

type LikedOfferCheckerInterface interface {
	IsOfferLikedByUser(offerID, userID uint) bool
}

type BidRetrieverInterface interface {
	GetByAuctionId(auctionID uint) ([]models.Bid, error)
}

type ManufacturerRetrieverInterface interface {
	GetAll() ([]models.Manufacturer, error)
}

type ModelRetrieverInterface interface {
	GetByManufacturerName(name string) ([]models.Model, error)
}

type SaleOfferServiceInterface interface {
	Create(in *CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error)
	Update(in *UpdateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error)
	GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error)
	GetByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error)
}

type SaleOfferService struct {
	saleOfferRepo  SaleOfferRepositoryInterface
	manRetriever   ManufacturerRetrieverInterface
	modelRetriever ModelRetrieverInterface
	imageRetriever ImageRetrieverInterface
	bidRetriever   BidRetrieverInterface
	likedChecker   LikedOfferCheckerInterface
}

func NewSaleOfferService(
	saleOfferRepository SaleOfferRepositoryInterface,
	manufacturerRetriever ManufacturerRetrieverInterface,
	modelRetriever ModelRetrieverInterface,
	bidRetriever BidRetrieverInterface,
	imageRetriever ImageRetrieverInterface,
	likedChecker LikedOfferCheckerInterface,
) SaleOfferServiceInterface {
	return &SaleOfferService{
		saleOfferRepo:  saleOfferRepository,
		manRetriever:   manufacturerRetriever,
		modelRetriever: modelRetriever,
		bidRetriever:   bidRetriever,
		imageRetriever: imageRetriever,
		likedChecker:   likedChecker,
	}
}

func (s *SaleOfferService) Create(in *CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return nil, err
	}
	modelID, err := s.getModelID(in.Manufacturer, in.Model)
	if err != nil {
		return nil, err
	}
	offer.Car.ModelID = modelID
	if err := s.saleOfferRepo.Create(offer); err != nil {
		return nil, err
	}
	return s.GetByID(offer.ID, &offer.UserID)
}

func (s *SaleOfferService) Update(in *UpdateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferRepo.GetByID(in.ID)
	if err != nil {
		return nil, err
	}
	updatedOffer, err := in.UpdateSaleOfferFromDTO(offer)
	if err != nil {
		return nil, err
	}
	if err := s.updateModelID(updatedOffer, in); err != nil {
		return nil, err
	}
	if err = s.saleOfferRepo.Update(updatedOffer); err != nil {
		return nil, err
	}
	return s.GetByID(updatedOffer.ID, &offer.UserID)
}

func (s *SaleOfferService) GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error) {
	manufacturers, err := s.manRetriever.GetAll()
	if err != nil {
		return nil, err
	}
	filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(manufacturers, manufacturer.MapToName)
	offers, pagResponse, err := s.saleOfferRepo.GetFiltered(filter)
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
	offer, err := s.saleOfferRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	offerDTO := MapToDetailedDTO(offer)
	if err = s.setUserFields(&offerDTO.UserContext, offer.ID, userID); err != nil {
		return nil, err
	}
	if err = s.setImagesUrls(offerDTO); err != nil {
		return nil, err
	}
	return offerDTO, nil
}

func (s *SaleOfferService) GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error) {
	offers, pagResponse, err := s.saleOfferRepo.GetByUserID(id, pagRequest)
	if err != nil {
		return nil, err
	}
	offerDTOs, err := s.mapOfferSliceWithAdditionalFields(offers, &id)
	if err != nil {
		return nil, err
	}
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) canBeModifiedByUser(offerID uint, userID *uint) (bool, error) {
	if userID == nil {
		return false, nil
	}
	offer, err := s.saleOfferRepo.GetByID(offerID)
	if err != nil {
		return false, err
	}
	if offer.Auction == nil {
		return true, nil
	}
	bids, err := s.bidRetriever.GetByAuctionId(offerID)
	if err != nil {
		return false, err
	}
	return len(bids) == 0, nil
}

func (s *SaleOfferService) isOfferLikedByUser(offerID uint, userID *uint) bool {
	if userID == nil {
		return false
	}
	return s.likedChecker.IsOfferLikedByUser(offerID, *userID)
}

func (s *SaleOfferService) getModelID(manufacturerName, modelName string) (uint, error) {
	models, err := s.modelRetriever.GetByManufacturerName(manufacturerName)
	if err != nil {
		return 0, ErrInvalidManufacturer
	}
	for _, model := range models {
		if model.Name == modelName {
			return model.ID, nil
		}
	}
	return 0, ErrInvalidModel
}

func (s *SaleOfferService) updateModelID(updatedOffer *models.SaleOffer, in *UpdateSaleOfferDTO) error {
	if in.Manufacturer == nil && in.Model != nil {
		modelID, err := s.getModelID(updatedOffer.Car.Model.Manufacturer.Name, *in.Model)
		if err != nil {
			return err
		}
		updatedOffer.Car.ModelID = modelID
	} else if in.Manufacturer != nil && in.Model != nil {
		modelID, err := s.getModelID(*in.Manufacturer, *in.Model)
		if err != nil {
			return err
		}
		updatedOffer.Car.ModelID = modelID
	}
	return nil
}

func (s *SaleOfferService) setImagesUrls(dto *RetrieveDetailedSaleOfferDTO) error {
	images, err := s.imageRetriever.GetImagesByOfferID(dto.ID)
	if err != nil {
		return err
	}
	urls := mapping.MapSliceToDTOs(images, func(m *models.Image) *string { return &m.Url })
	dto.ImagesUrls = urls
	return nil
}

func (s *SaleOfferService) setUserFields(userContext *UserContext, offerID uint, userID *uint) error {
	userContext.IsLiked = s.isOfferLikedByUser(offerID, userID)
	stmt, err := s.canBeModifiedByUser(offerID, userID)
	if err != nil {
		return err
	}
	userContext.CanModify = stmt
	return nil
}

func (s *SaleOfferService) mapOfferSliceWithAdditionalFields(offers []models.SaleOffer, userID *uint) ([]RetrieveSaleOfferDTO, error) {
	offerDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto := MapToDTO(&offer)
		if err := s.setUserFields(&dto.UserContext, offer.ID, userID); err != nil {
			return nil, err
		}
		offerDTOs = append(offerDTOs, *dto)
	}
	return offerDTOs, nil
}
