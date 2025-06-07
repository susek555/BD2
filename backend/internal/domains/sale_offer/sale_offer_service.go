package sale_offer

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
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
	GetByID(id uint) (*models.Model, error)
}

type PurchaseCreatorInterface interface {
	Create(purchase *models.Purchase) error
}

type SaleOfferPreparatornterface interface {
	PrepareForCreateSaleOffer(in *CreateSaleOfferDTO) (*models.SaleOffer, error)
	PrepareForUpdateSaleOffer(in *UpdateSaleOfferDTO, userID uint) (*models.SaleOffer, error)
	PrepareForBuySaleOffer(id uint, userID uint) (*models.SaleOffer, error)
}

type SaleOfferManagerInterface interface {
	Create(in *CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error)
	Update(in *UpdateSaleOfferDTO, userID uint) (*RetrieveDetailedSaleOfferDTO, error)
	Publish(id uint, userID uint) (*RetrieveDetailedSaleOfferDTO, error)
	Buy(id uint, userID uint) (*models.SaleOffer, error)
	Delete(id uint, userID uint) error
}

type SaleOfferRetrieverInterface interface {
	GetByID(id uint, userID *uint) (*RetrieveSaleOfferDTO, error)
	GetDetailedByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error)
	GetByUserID(id uint, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error)
	GetFiltered(filter *OfferFilter, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error)
}

type SaleOfferServiceInterface interface {
	SaleOfferPreparatornterface
	SaleOfferManagerInterface
	SaleOfferRetrieverInterface
}

type SaleOfferService struct {
	saleOfferRepo   SaleOfferRepositoryInterface
	manRetriever    ManufacturerRetrieverInterface
	modelRetriever  ModelRetrieverInterface
	imageRetriever  ImageRetrieverInterface
	accessEvaluator OfferAccessEvaluatorInterface
	purchaseCreator PurchaseCreatorInterface
}

func NewSaleOfferService(
	saleOfferRepository SaleOfferRepositoryInterface,
	manufacturerRetriever ManufacturerRetrieverInterface,
	modelRetriever ModelRetrieverInterface,
	imageRetriever ImageRetrieverInterface,
	accessEvaluator OfferAccessEvaluatorInterface,
	purchaseCreator PurchaseCreatorInterface,
) SaleOfferServiceInterface {
	return &SaleOfferService{
		saleOfferRepo:   saleOfferRepository,
		manRetriever:    manufacturerRetriever,
		modelRetriever:  modelRetriever,
		imageRetriever:  imageRetriever,
		accessEvaluator: accessEvaluator,
		purchaseCreator: purchaseCreator,
	}
}

func (s *SaleOfferService) Create(in *CreateSaleOfferDTO) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.PrepareForCreateSaleOffer(in)
	if err != nil {
		return nil, err
	}
	if err := s.saleOfferRepo.Create(offer); err != nil {
		return nil, err
	}
	return s.GetDetailedByID(offer.ID, &offer.UserID)
}

func (s *SaleOfferService) Update(in *UpdateSaleOfferDTO, userID uint) (*RetrieveDetailedSaleOfferDTO, error) {
	updatedOffer, err := s.PrepareForUpdateSaleOffer(in, userID)
	if err != nil {
		return nil, err
	}
	if err = s.saleOfferRepo.Update(updatedOffer); err != nil {
		return nil, err
	}
	return s.GetDetailedByID(updatedOffer.ID, &updatedOffer.UserID)
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
	if err := s.saleOfferRepo.UpdateStatus(offer, enums.PUBLISHED); err != nil {
		return nil, err
	}
	return s.GetDetailedByID(id, &userID)
}

func (s *SaleOfferService) Buy(id uint, userID uint) (*models.SaleOffer, error) {
	offer, err := s.PrepareForBuySaleOffer(id, userID)
	if err != nil {
		return nil, err
	}
	if offer.IsAuction {
		return nil, ErrOfferIsAuction
	}
	if err := s.saleOfferRepo.UpdateStatus(offer, enums.SOLD); err != nil {
		return nil, err
	}
	purchaseModel := &models.Purchase{OfferID: offer.ID, BuyerID: userID, FinalPrice: offer.Price, IssueDate: time.Now()}
	if err := s.purchaseCreator.Create(purchaseModel); err != nil {
		return nil, err
	}
	return s.saleOfferRepo.GetByID(id)
}

func (s *SaleOfferService) Delete(id uint, userID uint) error {
	offer, err := s.saleOfferRepo.GetByID(id)
	if err != nil {
		return err
	}
	if err := s.authorizeModificationByUser(offer, userID); err != nil {
		return err
	}
	return s.saleOfferRepo.Delete(id)
}

func (s *SaleOfferService) GetByID(id uint, userID *uint) (*RetrieveSaleOfferDTO, error) {
	offer, err := s.saleOfferRepo.GetViewByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapOfferWithAdditionalFields(offer, userID)
}

func (s *SaleOfferService) GetDetailedByID(id uint, userID *uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offer, err := s.saleOfferRepo.GetViewByID(id)
	if err != nil {
		return nil, err
	}
	return s.mapOfferWithAdditionalFieldsDetailed(offer, userID)
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

func (s *SaleOfferService) GetFiltered(filter *OfferFilter, pagRequest *pagination.PaginationRequest) (*RetrieveOffersWithPagination, error) {
	manufacturers, err := s.manRetriever.GetAll()
	if err != nil {
		return nil, err
	}
	filter.Constraints.Manufacturers = mapping.MapSliceToDTOs(manufacturers, manufacturer.MapToName)
	offers, pagResponse, err := s.saleOfferRepo.GetFiltered(filter, pagRequest)
	if err != nil {
		return nil, err
	}
	offerDTOs, err := s.mapOfferSliceWithAdditionalFields(offers, filter.UserID)
	if err != nil {
		return nil, err
	}
	return &RetrieveOffersWithPagination{Offers: offerDTOs, PaginationResponse: pagResponse}, nil
}

func (s *SaleOfferService) PrepareForCreateSaleOffer(in *CreateSaleOfferDTO) (*models.SaleOffer, error) {
	offer, err := in.MapToSaleOffer()
	if err != nil {
		return nil, err
	}
	modelID, err := s.getModelID(in.ManufacturerName, in.ModelName)
	if err != nil {
		return nil, err
	}
	offer.Car.ModelID = modelID
	offer.Status = enums.PENDING
	return offer, nil
}

func (s *SaleOfferService) PrepareForUpdateSaleOffer(in *UpdateSaleOfferDTO, userID uint) (*models.SaleOffer, error) {
	offer, err := s.saleOfferRepo.GetByID(in.ID)
	if err != nil {
		return nil, err
	}
	if !offer.BelongsToUser(userID) {
		return nil, ErrOfferNotOwned
	}
	if err := s.authorizeModificationByUser(offer, userID); err != nil {
		return nil, err
	}
	modelID, err := s.determineNewModelID(offer, in)
	if err != nil {
		return nil, err
	}
	updatedOffer, err := in.UpdateOfferFromDTO(offer)
	if err != nil {
		return nil, err
	}
	updatedOffer.Car.ModelID = modelID
	return updatedOffer, nil
}

func (s *SaleOfferService) PrepareForBuySaleOffer(id uint, userID uint) (*models.SaleOffer, error) {
	offer, err := s.saleOfferRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if offer.BelongsToUser(userID) {
		return nil, ErrOfferOwnedByUser
	}
	if offer.Status != enums.PUBLISHED {
		return nil, ErrOfferNotPublished
	}
	return offer, nil
}

func (s *SaleOfferService) getModelID(manufacturerName, modelName string) (uint, error) {
	model, err := s.modelRetriever.GetByManufacturerAndModelName(manufacturerName, modelName)
	if err != nil {
		return 0, ErrInvalidManufacturerModelPair
	}
	return model.ID, nil
}

func (s *SaleOfferService) determineNewModelID(offer *models.SaleOffer, dto *UpdateSaleOfferDTO) (uint, error) {
	if dto.ModelName == nil {
		return offer.Car.ModelID, nil
	}
	manufacturerName := ""
	if dto.ManufacturerName != nil {
		manufacturerName = *dto.ManufacturerName
	} else {
		model, err := s.modelRetriever.GetByID(offer.Car.ModelID)
		if err != nil {
			return 0, err
		}
		manufacturerName = model.Manufacturer.Name
	}
	modelID, err := s.getModelID(manufacturerName, *dto.ModelName)
	if err != nil {
		return 0, err
	}
	return modelID, nil
}

func (s *SaleOfferService) mapOfferSliceWithAdditionalFields(offers []views.SaleOfferView, userID *uint) ([]RetrieveSaleOfferDTO, error) {
	offerDTOs := make([]RetrieveSaleOfferDTO, 0, len(offers))
	for _, offer := range offers {
		dto, err := s.mapOfferWithAdditionalFields(&offer, userID)
		if err != nil {
			return nil, err
		}
		offerDTOs = append(offerDTOs, *dto)
	}
	return offerDTOs, nil
}

func (s *SaleOfferService) mapOfferWithAdditionalFields(offer *views.SaleOfferView, userID *uint) (*RetrieveSaleOfferDTO, error) {
	offerDTO := MapViewToDTO(offer)
	userContext, err := s.getUserContextFields(offer, userID)
	if err != nil {
		return nil, err
	}
	offerDTO.UserContext = *userContext
	urls, err := s.getOfferImagesURLs(offer)
	if err != nil {
		return nil, err
	}
	if len(urls) > 0 {
		offerDTO.MainURL = urls[0]
	}
	return offerDTO, nil
}

func (s *SaleOfferService) mapOfferWithAdditionalFieldsDetailed(offer *views.SaleOfferView, userID *uint) (*RetrieveDetailedSaleOfferDTO, error) {
	offerDTO := MapViewToDetailedDTO(offer)
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

func (s *SaleOfferService) getUserContextFields(offer *views.SaleOfferView, userID *uint) (*UserContext, error) {
	isLiked := s.accessEvaluator.IsOfferLikedByUser(offer, userID)
	canModify, err := s.accessEvaluator.CanBeModifiedByUser(offer, userID)
	if err != nil {
		return nil, err
	}
	return &UserContext{IsLiked: isLiked, CanModify: canModify}, nil
}

func (s *SaleOfferService) getOfferImagesURLs(offer *views.SaleOfferView) ([]string, error) {
	images, err := s.imageRetriever.GetByOfferID(offer.ID)
	if err != nil {
		return nil, err
	}
	return mapping.MapSliceToDTOs(images, func(m *models.Image) *string { return &m.Url }), nil
}

func (s *SaleOfferService) authorizeModificationByUser(offer SaleOfferEntityInterface, userID uint) error {
	if !offer.BelongsToUser(userID) {
		return ErrOfferNotOwned
	}
	canBeModifed, err := s.accessEvaluator.CanBeModifiedByUser(offer, &userID)
	if err != nil {
		return err
	}
	if !canBeModifed {
		return ErrOfferModification
	}
	return nil
}
