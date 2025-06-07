package sale_offer_tests

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/views"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
	"gorm.io/gorm"
)

// Custom mock repository
type mockSaleOfferRepository struct {
	createFunc       func(offer *models.SaleOffer) error
	getByIDFunc      func(id uint) (*models.SaleOffer, error)
	getViewByIDFunc  func(id uint) (*views.SaleOfferView, error)
	updateFunc       func(offer *models.SaleOffer) error
	updateStatusFunc func(offer *models.SaleOffer, status enums.Status) error
	deleteFunc       func(id uint) error
	getFilteredFunc  func(filter sale_offer.OfferFilterIntreface, pagination *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error)
}

func (m *mockSaleOfferRepository) Create(offer *models.SaleOffer) error {
	if m.createFunc != nil {
		return m.createFunc(offer)
	}
	return nil
}

func (m *mockSaleOfferRepository) GetByID(id uint) (*models.SaleOffer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockSaleOfferRepository) GetViewByID(id uint) (*views.SaleOfferView, error) {
	if m.getViewByIDFunc != nil {
		return m.getViewByIDFunc(id)
	}
	return nil, nil
}

func (m *mockSaleOfferRepository) Update(offer *models.SaleOffer) error {
	if m.updateFunc != nil {
		return m.updateFunc(offer)
	}
	return nil
}

func (m *mockSaleOfferRepository) UpdateStatus(offer *models.SaleOffer, status enums.Status) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(offer, status)
	}
	return nil
}

func (m *mockSaleOfferRepository) Delete(id uint) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockSaleOfferRepository) GetFiltered(filter sale_offer.OfferFilterIntreface, pagination *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
	if m.getFilteredFunc != nil {
		return m.getFilteredFunc(filter, pagination)
	}
	return nil, nil, nil
}

func (m *mockSaleOfferRepository) GetAllActiveAuctions() ([]views.SaleOfferView, error) {
	return []views.SaleOfferView{}, nil
}

type mockPurchaseCreator struct {
	createFunc func(purchase *models.Purchase) error
}

func (m *mockPurchaseCreator) Create(purchase *models.Purchase) error {
	if m.createFunc != nil {
		return m.createFunc(purchase)
	}
	return nil
}

type MockManufacturerRetrieverInterface struct {
	getAllFunc func() ([]models.Manufacturer, error)
}

func (m *MockManufacturerRetrieverInterface) GetAll() ([]models.Manufacturer, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return []models.Manufacturer{}, nil
}

type MockModelRetrieverInterface struct {
	getByManufacturerAndModelNameFunc func(manufacturerName, modelName string) (*models.Model, error)
	getByIDFunc                       func(id uint) (*models.Model, error)
}

func (m *MockModelRetrieverInterface) GetByManufacturerAndModelName(manufacturerName, modelName string) (*models.Model, error) {
	if m.getByManufacturerAndModelNameFunc != nil {
		return m.getByManufacturerAndModelNameFunc(manufacturerName, modelName)
	}
	return nil, nil
}

func (m *MockModelRetrieverInterface) GetByID(id uint) (*models.Model, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

type MockImageRetrieverInterface struct {
	getByOfferIDFunc func(offerID uint) ([]models.Image, error)
}

func (m *MockImageRetrieverInterface) GetByOfferID(offerID uint) ([]models.Image, error) {
	if m.getByOfferIDFunc != nil {
		return m.getByOfferIDFunc(offerID)
	}
	return []models.Image{}, nil
}

type MockImageRemoverInterface struct {
	deleteByFolderNameFunc func(folder string) error
}

func (m *MockImageRemoverInterface) DeleteByFolderName(folder string) error {
	if m.deleteByFolderNameFunc != nil {
		return m.deleteByFolderNameFunc(folder)
	}
	return nil
}

type MockOfferAccessEvaluatorInterface struct {
	canBeModifiedByUserFunc func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error)
	isOfferLikedByUserFunc  func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool
}

func (m *MockOfferAccessEvaluatorInterface) CanBeModifiedByUser(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
	if m.canBeModifiedByUserFunc != nil {
		return m.canBeModifiedByUserFunc(offer, userID)
	}
	return true, nil
}

func (m *MockOfferAccessEvaluatorInterface) IsOfferLikedByUser(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
	if m.isOfferLikedByUserFunc != nil {
		return m.isOfferLikedByUserFunc(offer, userID)
	}
	return false
}

func createMockSaleOfferService() (*sale_offer.SaleOfferService, *mockSaleOfferRepository, *MockManufacturerRetrieverInterface, *MockModelRetrieverInterface, *MockImageRetrieverInterface, *MockImageRemoverInterface, *MockOfferAccessEvaluatorInterface, *mockPurchaseCreator) {
	mockRepo := &mockSaleOfferRepository{}
	mockManufacturerRetriever := &MockManufacturerRetrieverInterface{}
	mockModelRetriever := &MockModelRetrieverInterface{}
	mockImageRetriever := &MockImageRetrieverInterface{}
	mockImageRemover := &MockImageRemoverInterface{}
	mockAccessEvaluator := &MockOfferAccessEvaluatorInterface{}
	mockPurchaseCreator := &mockPurchaseCreator{}

	service := sale_offer.NewSaleOfferService(
		mockRepo,
		mockManufacturerRetriever,
		mockModelRetriever,
		mockImageRetriever,
		mockImageRemover,
		mockAccessEvaluator,
		mockPurchaseCreator,
	).(*sale_offer.SaleOfferService)

	return service, mockRepo, mockManufacturerRetriever, mockModelRetriever, mockImageRetriever, mockImageRemover, mockAccessEvaluator, mockPurchaseCreator
}

func createSampleCreateDTO() *sale_offer.CreateSaleOfferDTO {
	return &sale_offer.CreateSaleOfferDTO{
		UserID:             1,
		Description:        "Test description",
		Price:              25000,
		Margin:             enums.LOW_MARGIN,
		Vin:                "123ABC456DEF7890",
		ProductionYear:     2001,
		Mileage:            50000,
		NumberOfDoors:      4,
		NumberOfSeats:      5,
		EnginePower:        150,
		EngineCapacity:     1000,
		RegistrationNumber: "ABC123",
		RegistrationDate:   "2005-01-15",
		Color:              enums.RED,
		FuelType:           enums.PETROL,
		Transmission:       enums.MANUAL,
		NumberOfGears:      6,
		Drive:              enums.FWD,
		ManufacturerName:   "Toyota",
		ModelName:          "Camry",
	}
}

func createSampleUpdateDTO() *sale_offer.UpdateSaleOfferDTO {
	description := "Updated description"
	price := uint(30000)
	return &sale_offer.UpdateSaleOfferDTO{
		ID:          1,
		Description: &description,
		Price:       &price,
	}
}

func createSampleSaleOffer() *models.SaleOffer {
	return &models.SaleOffer{
		ID:          1,
		UserID:      1,
		Description: "Test description",
		Price:       25000,
		Status:      enums.PENDING,
		DateOfIssue: time.Now(),
		IsAuction:   false,
		Car: &models.Car{
			ModelID: 1,
			Vin:     "123ABC456DEF7890",
		},
	}
}

func createSampleSaleOfferView() *views.SaleOfferView {
	return &views.SaleOfferView{
		ID:          1,
		UserID:      1,
		Username:    "testuser",
		Description: "Test car",
		Price:       25000,
		Status:      enums.PUBLISHED,
		IsAuction:   false,
	}
}

func createSampleModel() *models.Model {
	return &models.Model{
		ID:   1,
		Name: "Camry",
		Manufacturer: &models.Manufacturer{
			ID:   1,
			Name: "Toyota",
		},
	}
}

func TestSaleOfferService_Create_Success(t *testing.T) {
	service, mockRepo, _, mockModelRetriever, mockImageRetriever, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	createDTO := createSampleCreateDTO()
	sampleView := createSampleSaleOfferView()
	sampleModel := createSampleModel()

	// Set up mock behavior
	mockModelRetriever.getByManufacturerAndModelNameFunc = func(manufacturer, model string) (*models.Model, error) {
		if manufacturer == "Toyota" && model == "Camry" {
			return sampleModel, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	mockRepo.createFunc = func(offer *models.SaleOffer) error {
		offer.ID = 1
		return nil
	}
	mockRepo.getViewByIDFunc = func(id uint) (*views.SaleOfferView, error) {
		return sampleView, nil
	}

	mockAccessEvaluator.isOfferLikedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
		return false
	}
	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return true, nil
	}

	mockImageRetriever.getByOfferIDFunc = func(offerID uint) ([]models.Image, error) {
		return []models.Image{}, nil
	}

	result, err := service.Create(createDTO)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
}

func TestSaleOfferService_Create_InvalidModel(t *testing.T) {
	service, _, _, mockModelRetriever, _, _, _, _ := createMockSaleOfferService()

	createDTO := createSampleCreateDTO()

	mockModelRetriever.getByManufacturerAndModelNameFunc = func(manufacturer, model string) (*models.Model, error) {
		return nil, gorm.ErrRecordNotFound
	}

	result, err := service.Create(createDTO)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrInvalidManufacturerModelPair, err)
}

func TestSaleOfferService_Create_RepositoryError(t *testing.T) {
	service, mockRepo, _, mockModelRetriever, _, _, _, _ := createMockSaleOfferService()

	createDTO := createSampleCreateDTO()
	sampleModel := createSampleModel()

	mockModelRetriever.getByManufacturerAndModelNameFunc = func(manufacturer, model string) (*models.Model, error) {
		return sampleModel, nil
	}

	mockRepo.createFunc = func(offer *models.SaleOffer) error {
		return errors.New("database error")
	}

	result, err := service.Create(createDTO)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
}

func TestSaleOfferService_Update_Success(t *testing.T) {
	service, mockRepo, _, _, mockImageRetriever, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	updateDTO := createSampleUpdateDTO()
	sampleOffer := createSampleSaleOffer()
	sampleView := createSampleSaleOfferView()

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}
	mockRepo.updateFunc = func(offer *models.SaleOffer) error {
		return nil
	}
	mockRepo.getViewByIDFunc = func(id uint) (*views.SaleOfferView, error) {
		return sampleView, nil
	}

	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return true, nil
	}
	mockAccessEvaluator.isOfferLikedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
		return false
	}

	mockImageRetriever.getByOfferIDFunc = func(offerID uint) ([]models.Image, error) {
		return []models.Image{}, nil
	}

	result, err := service.Update(updateDTO, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
}

func TestSaleOfferService_Update_NotOwned(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	updateDTO := createSampleUpdateDTO()
	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 2 // Different user

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	result, err := service.Update(updateDTO, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferNotOwned, err)
}

func TestSaleOfferService_Update_CannotModify(t *testing.T) {
	service, mockRepo, _, _, _, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	updateDTO := createSampleUpdateDTO()
	sampleOffer := createSampleSaleOffer()

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return false, nil
	}

	result, err := service.Update(updateDTO, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferModification, err)
}

func TestSaleOfferService_Publish_Success(t *testing.T) {
	service, mockRepo, _, _, mockImageRetriever, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.Status = enums.READY
	sampleView := createSampleSaleOfferView()

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}
	mockRepo.updateStatusFunc = func(offer *models.SaleOffer, status enums.Status) error {
		return nil
	}
	mockRepo.getViewByIDFunc = func(id uint) (*views.SaleOfferView, error) {
		return sampleView, nil
	}

	mockAccessEvaluator.isOfferLikedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
		return false
	}
	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return true, nil
	}

	mockImageRetriever.getByOfferIDFunc = func(offerID uint) ([]models.Image, error) {
		return []models.Image{}, nil
	}

	result, err := service.Publish(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
}

func TestSaleOfferService_Publish_NotOwned(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 2 // Different user

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	result, err := service.Publish(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferNotOwned, err)
}

func TestSaleOfferService_Publish_NotReady(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.Status = enums.PENDING // Not READY

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	result, err := service.Publish(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferNotReadyToPublish, err)
}

func TestSaleOfferService_Buy_Success(t *testing.T) {
	service, mockRepo, _, _, _, _, _, mockPurchaseCreator := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 2 // Different user
	sampleOffer.Status = enums.PUBLISHED
	sampleOffer.IsAuction = false

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}
	mockRepo.updateStatusFunc = func(offer *models.SaleOffer, status enums.Status) error {
		return nil
	}

	mockPurchaseCreator.createFunc = func(purchase *models.Purchase) error {
		return nil
	}

	result, err := service.Buy(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
}

func TestSaleOfferService_Buy_OwnOffer(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 1 // Same user
	sampleOffer.Status = enums.PUBLISHED

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	result, err := service.Buy(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferOwnedByUser, err)
}

func TestSaleOfferService_Buy_NotPublished(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 2             // Different user
	sampleOffer.Status = enums.PENDING // Not published

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	result, err := service.Buy(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferNotPublished, err)
}

func TestSaleOfferService_Buy_IsAuction(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 2 // Different user
	sampleOffer.Status = enums.PUBLISHED
	sampleOffer.IsAuction = true

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	result, err := service.Buy(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sale_offer.ErrOfferIsAuction, err)
}

func TestSaleOfferService_Delete_Success(t *testing.T) {
	service, mockRepo, _, _, _, mockImageRemover, mockAccessEvaluator, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}
	mockRepo.deleteFunc = func(id uint) error {
		return nil
	}

	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return true, nil
	}

	mockImageRemover.deleteByFolderNameFunc = func(folder string) error {
		return nil
	}

	err := service.Delete(1, 1)

	assert.NoError(t, err)
}

func TestSaleOfferService_Delete_NotOwned(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()
	sampleOffer.UserID = 2 // Different user

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	err := service.Delete(1, 1)

	assert.Error(t, err)
	assert.Equal(t, sale_offer.ErrOfferNotOwned, err)
}

func TestSaleOfferService_Delete_CannotModify(t *testing.T) {
	service, mockRepo, _, _, _, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	sampleOffer := createSampleSaleOffer()

	mockRepo.getByIDFunc = func(id uint) (*models.SaleOffer, error) {
		return sampleOffer, nil
	}

	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return false, nil
	}

	err := service.Delete(1, 1)

	assert.Error(t, err)
	assert.Equal(t, sale_offer.ErrOfferModification, err)
}

func TestSaleOfferService_GetByID_Success(t *testing.T) {
	service, mockRepo, _, _, mockImageRetriever, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	sampleView := createSampleSaleOfferView()
	userID := uint(1)

	mockRepo.getViewByIDFunc = func(id uint) (*views.SaleOfferView, error) {
		return sampleView, nil
	}

	mockAccessEvaluator.isOfferLikedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
		return true
	}
	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return false, nil
	}

	mockImageRetriever.getByOfferIDFunc = func(offerID uint) ([]models.Image, error) {
		return []models.Image{{Url: "http://example.com/image1.jpg"}}, nil
	}

	result, err := service.GetByID(1, &userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.True(t, result.IsLiked)
	assert.False(t, result.CanModify)
	assert.Equal(t, "http://example.com/image1.jpg", result.MainURL)
}

func TestSaleOfferService_GetByID_NotFound(t *testing.T) {
	service, mockRepo, _, _, _, _, _, _ := createMockSaleOfferService()

	userID := uint(1)

	mockRepo.getViewByIDFunc = func(id uint) (*views.SaleOfferView, error) {
		return nil, gorm.ErrRecordNotFound
	}

	result, err := service.GetByID(1, &userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestSaleOfferService_GetDetailedByID_Success(t *testing.T) {
	service, mockRepo, _, _, mockImageRetriever, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	sampleView := createSampleSaleOfferView()
	userID := uint(1)

	mockRepo.getViewByIDFunc = func(id uint) (*views.SaleOfferView, error) {
		return sampleView, nil
	}

	mockAccessEvaluator.isOfferLikedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
		return false
	}
	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return true, nil
	}

	mockImageRetriever.getByOfferIDFunc = func(offerID uint) ([]models.Image, error) {
		return []models.Image{
			{Url: "http://example.com/image1.jpg"},
			{Url: "http://example.com/image2.jpg"},
		}, nil
	}

	result, err := service.GetDetailedByID(1, &userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.False(t, result.IsLiked)
	assert.True(t, result.CanModify)
	assert.Len(t, result.ImagesUrls, 2)
	assert.Equal(t, "http://example.com/image1.jpg", result.ImagesUrls[0])
}

func TestSaleOfferService_GetFiltered_Success(t *testing.T) {
	service, mockRepo, mockManufacturerRetriever, _, mockImageRetriever, _, mockAccessEvaluator, _ := createMockSaleOfferService()

	userID := uint(1)
	filter := &sale_offer.PublishedOffersOnlyFilter{
		BaseOfferFilter: sale_offer.BaseOfferFilter{
			UserID: &userID,
		},
	}
	pagRequest := &pagination.PaginationRequest{Page: 1, PageSize: 10}
	sampleView := createSampleSaleOfferView()
	mockManufacturerRetriever.getAllFunc = func() ([]models.Manufacturer, error) {
		return []models.Manufacturer{{ID: 1, Name: "Toyota"}}, nil
	}
	mockRepo.getFilteredFunc = func(filter sale_offer.OfferFilterIntreface, pagination *pagination.PaginationRequest) ([]views.SaleOfferView, *pagination.PaginationResponse, error) {
		return []views.SaleOfferView{*sampleView}, nil, nil
	}

	mockAccessEvaluator.isOfferLikedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) bool {
		return false
	}
	mockAccessEvaluator.canBeModifiedByUserFunc = func(offer sale_offer.SaleOfferEntityInterface, userID *uint) (bool, error) {
		return true, nil
	}

	mockImageRetriever.getByOfferIDFunc = func(offerID uint) ([]models.Image, error) {
		return []models.Image{}, nil
	}

	result, err := service.GetFiltered(filter, pagRequest)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Offers, 1)
	assert.Equal(t, uint(1), result.Offers[0].ID)
}
