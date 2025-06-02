package sale_offer_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"gorm.io/gorm"
)

// ------------------
// Create offer tests
// ------------------

var db, _ = setupDB()

func TestCreateOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*createSaleOfferDTO())
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*createSaleOfferDTO())
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	token := "invalid_token"
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestCreateOffer_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal("")
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_EmptyStruct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.CreateSaleOfferDTO{})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrMissingFields.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_SettingUserIDManually(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("UserID", uint(3))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, got.Username)
	u.CleanDB(DB)
}

func TestCreateOffer_NonExistentManufacturer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("ManufacturerName", "invalid")))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description, sale_offer.ErrInvalidManufacturer.Error())
}

func TestCreateOffer_NonExistentModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("ModelName", "invalid")))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfDoorsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfDoors", uint(10))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfDoorsLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfDoors", uint(0))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}
func TestCreateOffer_ValidNumberOfDoors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfDoors", uint(5))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfSeatsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfSeats", uint(101))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfSeatsLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfSeats", uint(1))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_ValidNumberOfSeats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfSeats", uint(5))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidEnginePowerGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EnginePower", uint(10000))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}
func TestCreateOffer_InvalidEnginePowerLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EnginePower", uint(0))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_ValidEnginePower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EnginePower", uint(100))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidEngineCapacityGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EngineCapacity", uint(9001))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidEngineCapacityLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EngineCapacity", uint(0))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_ValidEngineCapacity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EngineCapacity", uint(1000))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfGearsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfGears", uint(11))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfGearsLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfGears", uint(0))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestCreateOffer_ValidNumberOfGears(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfGears", uint(5))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidColor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Color", enums.Color("invalid_color"))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidColor.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_ValidColor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Color", enums.BLACK)))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidFuelType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("FuelType", enums.FuelType("invalid_fuel_type"))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidFuelType.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_ValidFuelType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("FuelType", enums.PETROL)))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidTransmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Transmission", enums.Transmission("invalid_transmission"))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidTransmission.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_ValidTransmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Transmission", enums.MANUAL)))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidDrive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Drive", enums.Drive("invalid_drive"))))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidDrive.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_ValidDrive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, svc, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Drive", enums.FWD)))
	setOffersStatusToPublished(db)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
}

// -----------------------
// Update sale offer tests
// ----------------------

func TestUpdateOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{})
	assert.NoError(t, err)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestUpdateOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{})
	assert.NoError(t, err)
	token := "invalid_token"
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
}

func TestUpdateOffer_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal("")
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
}

func TestUpdateOffer_NotUsersOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1})
	assert.NoError(t, err)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferNotOwned.Error())
}

func TestUpdateOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestUpdateOffer_Description(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _ := newTestServer(db, seedOffers)
	description := "Updated description"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Description: &description})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Description, description)
	u.CleanDB(DB)
}

// ------------------------------
// Get filtered sale offers tests
// For more get filtered tests see: sale_offer_filter_test.go
// ------------------------------

func TestGetFiltered_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _ := newTestServer(db, seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	u.CleanDB(DB)
}

func TestGetFiltered_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, a := newTestServer(db, seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0], a, nil))
	u.CleanDB(DB)
}

func TestGetFiltered_OneAuction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	server, _, _ := newTestServer(db, seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.Equal(t, *sale_offer.MapToDTO(&seedOffers[0]), got.Offers[0])
	u.CleanDB(DB)
}

func TestGetFiltered_AuctionsAndOffersCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}

	server, _, _ := newTestServer(db, seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		assert.Equal(t, *sale_offer.MapToDTO(&seedOffers[i]), got.Offers[i])
	}
	u.CleanDB(DB)
}

func TestGetFiltered_AuthorizedOtherUserOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, a := newTestServer(db, seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[i], got.Offers[i], a, &user.ID))
	}
	u.CleanDB(DB)
}

func TestGetFiltered_AuthorizedMyOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}

	server, _, _ := newTestServer(db, seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got.Offers)
	u.CleanDB(DB)
}

// ----------------------
// Get offer by id tests
// ----------------------

func TestGetSaleOfferById_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestGetSaleOfferById_NonExistentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/2", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestGetSaleOfferById_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/-1", nil, nil)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)

	u.CleanDB(DB)
}

func TestGetSaleOfferById_StringID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/abc", nil, nil)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestGetById_RegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveSaleOfferDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, seedOffers[0].ID, got.ID)
	u.CleanDB(DB)
}

func TestGetById_AuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveSaleOfferDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, seedOffers[0].ID, got.ID)
	u.CleanDB(DB)
}

// ----------------------
// Get my offers tests
// ----------------------

func TestGetMyOffers_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestGetMyOffers_NoOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got.Offers)
	u.CleanDB(DB)
}

func TestGetMyOffers_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, a := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0], a, &user.ID))
	u.CleanDB(DB)
}

func TestGetMyOffers_OneAuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}

	server, _, a := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0], a, &user.ID))
	u.CleanDB(DB)
}

func TestGetMyOffers_AuctionsAndOffersCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}

	server, _, a := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[i], got.Offers[i], a, &user.ID))
	}
	u.CleanDB(DB)
}

// ----------
// Like tests
// ----------

func TestLikeOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestLikeOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/2", nil, &token)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestLikeOffer_OwnOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, liked_offer.ErrLikeOwnOffer.Error())
}

func TestLikeOffer_Successful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	likedOfferRepo := liked_offer.NewLikedOfferRepository(db)
	likedOffers, err := likedOfferRepo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.True(t, len(likedOffers) == 1)
	u.CleanDB(DB)
}

func TestLikeOffer_AlreadyLiked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, liked_offer.ErrLikeAlreadyLikedOffer.Error())
}

// -------------
// Dislike tests
// -------------

func TestDislikeOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/favourite/dislike/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestDislikeOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/favourite/dislike/2", nil, &token)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestDislikeOffer_NotLikedOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/favourite/dislike/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, liked_offer.ErrDislikeNotLikedOffer.Error())
}

func TestDislikeOffer_Successful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	_, receivedStatus = u.PerformRequest(server, http.MethodDelete, "/favourite/dislike/1", nil, &token)
	assert.Equal(t, http.StatusNoContent, receivedStatus)
	likedOfferRepo := liked_offer.NewLikedOfferRepository(db)
	likedOffers, err := likedOfferRepo.GetByUserID(user.ID)
	assert.NoError(t, err)
	assert.True(t, len(likedOffers) == 0)
	u.CleanDB(DB)
}

// ----------
// Test enums
// ----------
func TestGetOfferTypes_Successful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/offer-types", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got map[string][]sale_offer.OfferType
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got["offer_types"], sale_offer.OfferTypes)
	u.CleanDB(DB)
}

func TestGetOrderKeys_Successful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/order-keys", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got map[string][]string
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.ElementsMatch(t, got["order_keys"], sale_offer.GetKeysFromMap(sale_offer.OrderKeysMap))
	u.CleanDB(DB)
}
