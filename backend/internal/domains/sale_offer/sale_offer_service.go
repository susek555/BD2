package sale_offer

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/mapping"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type ImageRetrieverInterface interface {
	GetByOfferID(offerID uint) ([]models.Image, error)
}

type ManufacturerRetrieverInterface interface {
	GetAll() ([]models.Manufacturer, error)
}

type ModelRetrieverInterface interface {
	GetByManufacturerAndModelName(manufacturerName, modelName string) (*models.Model, error)
}

//go:generate mockery --name=SaleOfferServiceInterface --output=../../test/mocks --case=snake --with-expecter
type SaleOfferServiceInterface interface {
	Create(in *CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error)
	Update(in *UpdateSaleOfferDTO, userID uint) (*RetrieveDetailedSaleOfferDTO, error)
	Publish(id uint, userID uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error)
	GetFiltered(filter *OfferFilter) (*RetrieveOffersWithPagination, error)
	GetModelID(manufacturerName, modelName string) (uint, error)
	DetermineNewModelID(offer *models.SaleOffer, dto *UpdateSaleOfferDTO) (uint, error)
	Buy(offerID uint, userID uint) (*models.SaleOffer, error)
	GetByIdNonDTO(id uint, userID uint) (*models.SaleOffer, error)
}
type SaleOfferService struct {
	saleOfferRepo   SaleOfferRepositoryInterface
	manRetriever    ManufacturerRetrieverInterface
	modelRetriever  ModelRetrieverInterface
	imageRetriever  ImageRetrieverInterface
	accessEvaluator OfferAccessEvaluatorInterface
}

func NewSaleOfferService(
	saleOfferRepository SaleOfferRepositoryInterface,
	manufacturerRetriever ManufacturerRetrieverInterface,
	modelRetriever ModelRetrieverInterface,
	imageRetriever ImageRetrieverInterface,
	accessEvaluator OfferAccessEvaluatorInterface,
) SaleOfferServiceInterface {
	return &SaleOfferService{
		saleOfferRepo:   saleOfferRepository,
		manRetriever:    manufacturerRetriever,
		modelRetriever:  modelRetriever,
		imageRetriever:  imageRetriever,
		accessEvaluator: accessEvaluator,
	}
}

func (s *SaleOfferService) Create(in *CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return nil, err
	}
	modelID, err := s.GetModelID(in.ManufacturerName, in.ModelName)
	if err != nil {
		return nil, err
	}
	offer.Car.ModelID = modelID
	offer.Status = enums.PENDING
	if err := s.saleOfferRepo.Create(offer); err != nil {
		return nil, err
	}
	return s.GetByID(offer.ID, &offer.UserID)
}

func (s *SaleOfferService) Update(in *UpdateSaleOfferDTO, userID uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferRepo.GetByID(in.ID)
	if err != nil {
		return nil, err
	}
	if !offer.BelongsToUser(userID) {
		return nil, ErrOfferNotOwned
	}
	modelID, err := s.DetermineNewModelID(offer, in)
	if err != nil {
		return nil, err
	}
	updatedOffer, err := in.UpdatedOfferFromDTO(offer)
	if err != nil {
		return nil, err
	}
	updatedOffer.Car.ModelID = modelID
	if err = s.saleOfferRepo.Update(updatedOffer); err != nil {
		return nil, err
	}
	return s.GetByID(offer.ID, &offer.UserID)
}

func (s *SaleOfferService) Publish(id uint, userID uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !offer.BelongsToUser(userID) {
		return nil, ErrOfferNotOwned
	}
	if offer.Status != enums.READY {
		return nil, ErrOfferNotReadyToPublish
	}
	offer.Status = enums.PUBLISHED
	if err := s.saleOfferRepo.Update(offer); err != nil {
		return nil, err
	}
	return s.GetByID(id, &userID)
}

func (s *SaleOfferService) GetByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	offerDTO := MapToDetailedDTO(offer)
	userContext, err := s.getUserContextFields(offer, userID)
	if err != nil {
		return nil, err
	}
	offerDTO.UserContext = *userContext
	urls, err := s.getOfferImagesURLs(offer)
	if err != nil {
		return nil, err
	}
	offerDTO.ImagesUrls = urls
	if userID != nil && offer.BelongsToUser(*userID) {
		offerDTO.Status = offer.Status
	}
	return offerDTO, nil
}

func (s *SaleOfferService) GetByIdNonDTO(id uint, userID uint) (*models.SaleOffer, error) {
	offer, err := s.saleOfferRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return offer, nil
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

func (s *SaleOfferService) GetModelID(manufacturerName, modelName string) (uint, error) {
	model, err := s.modelRetriever.GetByManufacturerAndModelName(manufacturerName, modelName)
	if err != nil {
		return 0, ErrInvalidManufacturerModelPair
	}
	return model.ID, nil
}

func (s *SaleOfferService) DetermineNewModelID(offer *models.SaleOffer, dto *UpdateSaleOfferDTO) (uint, error) {
	if dto.Model == nil {
		return 0, nil
	}
	manufacturerName := offer.Car.Model.Manufacturer.Name
	if dto.Manufacturer != nil {
		manufacturerName = *dto.Manufacturer
	}
	modelID, err := s.GetModelID(manufacturerName, *dto.Model)
	if err != nil {
		return 0, err
	}
	return modelID, err
}

func (s *SaleOfferService) Buy(offerID uint, userID uint) (*models.SaleOffer, error) {
	offer, err := s.saleOfferRepo.GetByID(offerID)
	if err != nil {
		return nil, err
	}
	if offer.BelongsToUser(userID) {
		return nil, ErrOfferOwnedByUser
	}
	offer, err = s.saleOfferRepo.BuyOffer(offerID, userID)
	return offer, err
}
func (s *SaleOfferService) getOfferImagesURLs(offer *models.SaleOffer) ([]string, error) {
	images, err := s.imageRetriever.GetByOfferID(offer.ID)
	if err != nil {
		return nil, err
	}
	return mapping.MapSliceToDTOs(images, func(m *models.Image) *string { return &m.Url }), nil
}

func (s *SaleOfferService) getUserContextFields(offer *models.SaleOffer, userID *uint) (*UserContext, error) {
	isLiked := s.accessEvaluator.IsOfferLikedByUser(offer, userID)
	canModify, err := s.accessEvaluator.CanBeModifiedByUser(offer, userID)
	if err != nil {
		return nil, err
	}
	return &UserContext{IsLiked: isLiked, CanModify: canModify}, nil
}

func (s *SaleOfferService) mapOfferSliceWithAdditionalFields(offers []models.SaleOffer, userID *uint) ([]RetrieveSaleOfferDTO, error) {
	offerDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto := MapToDTO(&offer)
		userContext, err := s.getUserContextFields(&offer, userID)
		if err != nil {
			return nil, err
		}
		dto.UserContext = *userContext
		urls, err := s.getOfferImagesURLs(&offer)
		if err != nil {
			return nil, err
		}
		if len(urls) > 0 {
			dto.MainURL = urls[0]
		}
		if userID != nil && offer.BelongsToUser(*userID) {
			dto.Status = offer.Status
		}
		offerDTOs = append(offerDTOs, *dto)
	}
	return offerDTOs, nil
}

func isAuction(offer *models.SaleOffer) bool {
	return offer.Auction != nil
}
