package sale_offer_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/purchase"
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*createSaleOfferDTO())
	assert.NoError(t, err)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*createSaleOfferDTO())
	assert.NoError(t, err)
	token := "invalid_token"
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestCreateOffer_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("ManufacturerName", "invalid")))
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("ModelName", "invalid")))
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

func TestCreateOffer_ManufacturerAndModelDontMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(),
		u.WithField[sale_offer.CreateSaleOfferDTO]("ManufacturerName", "Toyota"),
		u.WithField[sale_offer.CreateSaleOfferDTO]("ModelName", "Civic")))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidManufacturerModelPair.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_DescriptionTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	s := strings.Repeat("a", 2001)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Description", s)))
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

func TestCreateOffer_VinTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	s := strings.Repeat("a", 18)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Vin", s)))
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

func TestCreateOffer_RegistrationNumberTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	s := string(make([]byte, 21))
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("RegistrationNumber", s)))
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

func TestCreateOffer_InvalidMargin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Margin", enums.MarginValue(100))))
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

func TestCreateOffer_ValidMargin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Margin", enums.HIGH_MARGIN)))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfDoorsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfDoors", uint(10))))
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfDoors", uint(0))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfDoors", uint(5))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfSeatsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfSeats", uint(101))))
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfSeats", uint(1))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfSeats", uint(5))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidEnginePowerGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EnginePower", uint(10000))))
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EnginePower", uint(0))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EnginePower", uint(100))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidEngineCapacityGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EngineCapacity", uint(9001))))
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EngineCapacity", uint(0))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("EngineCapacity", uint(1000))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidNumberOfGearsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfGears", uint(11))))
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
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfGears", uint(0))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("NumberOfGears", uint(5))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidColor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Color", enums.Color("invalid_color"))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Color", enums.BLACK)))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidFuelType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("FuelType", enums.FuelType("invalid_fuel_type"))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("FuelType", enums.PETROL)))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidTransmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Transmission", enums.Transmission("invalid_transmission"))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Transmission", enums.MANUAL)))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
	u.CleanDB(DB)
}

func TestCreateOffer_InvalidDrive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Drive", enums.Drive("invalid_drive"))))
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
	server, svc, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("Drive", enums.FWD)))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusCreated, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.True(t, wasEntityAddedToDB[models.Car](db, uint(1)))
	assert.False(t, wasEntityAddedToDB[models.Auction](db, uint(1)))
}

func TestCreateOffer_ProductionYearTooOld(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("ProductionYear", uint(1315))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidProductionYear.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_ProductionYearInTheFuture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("ProductionYear", uint(3000))))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidProductionYear.Error())
	u.CleanDB(DB)
}

func TestCreateOffer_RegistrationDateInTheFuture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(*u.Build(createSaleOfferDTO(), u.WithField[sale_offer.CreateSaleOfferDTO]("RegistrationDate", "2030-01-01")))
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description, sale_offer.ErrInvalidRegistrationDate.Error())
	u.CleanDB(DB)
}

// -----------------------
// Update sale offer tests
// ----------------------

func TestUpdateOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{})
	assert.NoError(t, err)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestUpdateOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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

func TestUpdateOffer_OfferAlreadySold(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Status", enums.SOLD))}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusConflict, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferAlreadySold.Error())
}

func TestUpdateOffer_ManufacturerAndModelDontMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	invalidManufacturer := "InvalidManufacturer"
	invalidModel := "InvalidModel"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, ManufacturerName: &invalidManufacturer, ModelName: &invalidModel})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidManufacturerModelPair.Error())
}

func TestUpdateOffer_ModelOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*u.Build(createOffer(1),
		withCarField(u.WithField[models.Car]("ModelID", uint(4))))} //  Toyota Supra
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	model := "Corolla"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, ModelName: &model})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Model, model)
	u.CleanDB(DB)
}

func TestUpdateOffer_ModelAndManufacturer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*u.Build(createOffer(1),
		withCarField(u.WithField[models.Car]("ModelID", uint(4))))} //  Toyota Supra
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	manufacturer := "Audi"
	model := "A3"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, ModelName: &model, ManufacturerName: &manufacturer})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Model, model)
	assert.Equal(t, offer.Brand, manufacturer)
	u.CleanDB(DB)
}

func TestUpdateOffer_Description(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
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
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Description, description)
	u.CleanDB(DB)
}

func TestUpdateOffer_DescriptionTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	s := strings.Repeat("a", 2001)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Description: &s})
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

func TestUpdateOffer_Vin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	vn := "1HGCM82633A123456"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Vin: &vn})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Vin, vn)
	u.CleanDB(DB)
}

func TestUpdateOffer_VinTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	vn := strings.Repeat("a", 18)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Vin: &vn})
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

func TestUpdateOffer_RegistrationNumberTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	regNum := strings.Repeat("a", 21)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, RegistrationNumber: &regNum})
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

func TestUpdateOffer_RegistrationNumber(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	regNum := "ABC1234"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, RegistrationNumber: &regNum})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.RegistrationNumber, regNum)
	u.CleanDB(DB)
}

func TestUpdateOffer_Price(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	price := uint(10000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Price: &price})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Price, price)
	u.CleanDB(DB)
}

func TestUpdateOffer_Mileage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	mileage := uint(10000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Mileage: &mileage})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Mileage, mileage)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidMargin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	marginValue := enums.MarginValue(1000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Margin: &marginValue})
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

func TestUpdateOffer_ValidMargin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	marginValue := enums.MarginValue(10)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Margin: &marginValue})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Margin, marginValue)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidNumberOfDoorsLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	numberOfDoors := uint(0)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfDoors: &numberOfDoors})
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

func TestUpdateOffer_InvalidNumberOfDoorsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	numberOfDoors := uint(1000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfDoors: &numberOfDoors})
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

func TestUpdateOffer_ValidNumberOfDoors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	numberOfDoors := uint(4)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfDoors: &numberOfDoors})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.NumberOfDoors, numberOfDoors)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidNumberOfSeatsLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	numberOfSeats := uint(0)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfSeats: &numberOfSeats})
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

func TestUpdateOffer_InvalidNumberOfSeatsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	numberOfSeats := uint(1000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfSeats: &numberOfSeats})
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

func TestUpdateOffer_ValidNumberOfSeats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	numberOfSeats := uint(5)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfSeats: &numberOfSeats})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.NumberOfSeats, numberOfSeats)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidEnginePowerLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	enginePower := uint(0)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, EnginePower: &enginePower})
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

func TestUpdateOffer_InvalidEnginePowerGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	enginePower := uint(10000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, EnginePower: &enginePower})
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

func TestUpdateOffer_ValidEnginePower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	enginePower := uint(1500)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, EnginePower: &enginePower})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.EnginePower, enginePower)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidEngineCapacityLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	engineCapacity := uint(0)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, EngineCapacity: &engineCapacity})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidEngineCapacityGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	engineCapacity := uint(10000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, EngineCapacity: &engineCapacity})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_ValidEngineCapacity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	engineCapacity := uint(2000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, EngineCapacity: &engineCapacity})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.EngineCapacity, engineCapacity)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidNumberOfGearsLower(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	numberOfGears := uint(0)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfGears: &numberOfGears})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidNumberOfGearsGreater(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	numberOfGears := uint(1000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfGears: &numberOfGears})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_ValidNumberOfGears(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	numberOfGears := uint(6)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, NumberOfGears: &numberOfGears})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.NumberOfGears, numberOfGears)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidColor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	color := enums.Color("InvalidColor")
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Color: &color})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_ValidColor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	color := enums.BLACK
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Color: &color})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Color, color)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidFuelType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	fuelType := enums.FuelType("InvalidFuel")
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, FuelType: &fuelType})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_ValidFuelType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	fuelType := enums.DIESEL
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, FuelType: &fuelType})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.FuelType, fuelType)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidDrive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	driveType := enums.Drive("InvalidDrive")
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Drive: &driveType})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_ValidDrive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	driveType := enums.AWD
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Drive: &driveType})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Drive, driveType)
	u.CleanDB(DB)
}

func TestUpdateOffer_InvalidTransmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	transmissionType := enums.Transmission("InvalidTransmission")
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Transmission: &transmissionType})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestUpdateOffer_ValidTransmission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	transmissionType := enums.AUTOMATIC
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, Transmission: &transmissionType})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.Transmission, transmissionType)
	u.CleanDB(DB)
}

func TestUpdateOffer_ProductionYearTooOld(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	year := uint(1315)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, ProductionYear: &year})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidProductionYear.Error())
	u.CleanDB(DB)
}

func TestUpdateOffer_ProductionYearInTheFuture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	year := uint(3000)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, ProductionYear: &year})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrInvalidProductionYear.Error())
	u.CleanDB(DB)
}

func TestUpdateOffer_ProductionYear(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	year := uint(2020)
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, ProductionYear: &year})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.ProductionYear, year)
	u.CleanDB(DB)
}

func TestUpdateOffer_RegistrationDateInTheFuture(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	registrationDate := "2030-01-01"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, RegistrationDate: &registrationDate})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description, sale_offer.ErrInvalidRegistrationDate.Error())
	u.CleanDB(DB)
}

func TestUpdateOffer_RegistrationDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	registrationDate := "2020-01-01"
	body, err := json.Marshal(sale_offer.UpdateSaleOfferDTO{ID: 1, RegistrationDate: &registrationDate})
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, offer.RegistrationDate, registrationDate)
	u.CleanDB(DB)
}

func TestPubishSaleOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/publish/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
}

func TestPublishSaleOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	token := "invalid_token"
	_, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/publish/1", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
}

func TestPublishSaleOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/publish/1", nil, &token)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestPublishSaleOffer_NotUsersOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/publish/1", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferNotOwned.Error())
}

func TestPublishSaleOffer_NotReadyToPublish(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPut, "/sale-offer/publish/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestPublishSaleOffer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	seedOffers[0].Status = enums.READY

	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)

	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)

	dto := createSaleOfferDTO()
	svc.Create(dto)

	imagesDir := filepath.Join("images")

	filenames := []string{"test1.png", "test2.png", "test3.png"}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, fname := range filenames {
		filePath := filepath.Join(imagesDir, fname)
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatalf("cannot open file %s: %v", filePath, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("images", fname)
		if err != nil {
			t.Fatalf("CreateFormFile err: %v", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			t.Fatalf("io.Copy err %s: %v", fname, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Err %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/image/1", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "200 status code expected after uploading images")

	publishPath := "/sale-offer/publish/1"
	publishReq := httptest.NewRequest(http.MethodPut, publishPath, nil)
	publishReq.Header.Set("Authorization", "Bearer "+token)

	publishW := httptest.NewRecorder()
	server.ServeHTTP(publishW, publishReq)

	assert.Equal(t, http.StatusOK, publishW.Code, "200 status code expected after publishing offer")

	var got sale_offer.RetrieveDetailedSaleOfferDTO
	if err := json.Unmarshal(publishW.Body.Bytes(), &got); err != nil {
		t.Fatalf("Marshall err: %v", err)
	}

	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, enums.PUBLISHED, offer.Status)
	u.CleanDB(db)
}

// --------------------
// Buy sale offer Tests
// --------------------

func TestBuySaleOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestBuySaleOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	token := "invalid"
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestBuySaleOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestBuySaleOffer_OwnOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferOwnedByUser.Error())
	u.CleanDB(DB)
}

func TestBuySaleOffer_NotPublishedPendingStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Status", enums.PENDING))}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferNotPublished.Error())
	u.CleanDB(DB)
}

func TestBuySaleOffer_NotPublishedReadyStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Status", enums.READY))}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferNotPublished.Error())
	u.CleanDB(DB)
}

func TestBuySaleOffer_NotPublishedSoldStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Status", enums.SOLD))}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferNotPublished.Error())
	u.CleanDB(DB)
}

func TestBuySaleOffer_OfferIsAuction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createAuctionSaleOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferIsAuction.Error())
	u.CleanDB(DB)
}

func TestBuySaleOffer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Status", enums.PUBLISHED))}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/buy/1", nil, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveDetailedSaleOfferDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	offer, err := svc.GetDetailedByID(1, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, *offer, got)
	assert.Equal(t, enums.SOLD, offer.Status)
	purchaseRepo := purchase.NewPurchaseRepository(DB)
	purchase, err := purchaseRepo.GetByID(offer.ID)
	assert.NoError(t, err)
	assert.Equal(t, purchase.OfferID, offer.ID)
	u.CleanDB(DB)
}

// -----------------------
// Delete sale offer tests
// -----------------------

func TestDeleteOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/sale-offer/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestDeleteOffer_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	token := "invalid"
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/sale-offer/1", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestDeleteOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/sale-offer/1", nil, &token)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestDeleteOffer_NotUsersOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/sale-offer/1", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferNotOwned.Error())
	u.CleanDB(DB)
}

func TestDeleteOffer_OfferAlreadySold(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{
		*u.Build(createOffer(1), u.WithField[models.SaleOffer]("Status", enums.SOLD))}
	db, _ := setupDB()
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/sale-offer/1", nil, &token)
	assert.Equal(t, http.StatusConflict, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, sale_offer.ErrOfferAlreadySold.Error())
	u.CleanDB(DB)
}

func TestDeleteOffer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer = []models.SaleOffer{*createOffer(1)}
	db, _ := setupDB()
	server, svc, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/sale-offer/1", nil, &token)
	assert.Equal(t, http.StatusNoContent, receivedStatus)
	_, err := svc.GetDetailedByID(1, &user.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	u.CleanDB(DB)
}

// ------------------------------
// Get filtered sale offers tests
// For more get filtered tests see: sale_offer_filter_test.go
// ------------------------------

func TestGetFiltered_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.Equal(t, seedOffers[0].ID, got.Offers[0].ID)
	u.CleanDB(DB)
}

func TestGetFiltered_OneAuction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	server, s, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	offer, _ := s.GetByID(1, nil)
	assert.Equal(t, *offer, got.Offers[0])
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

	server, s, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		offer, _ := s.GetByID(uint(i+1), nil)
		assert.Equal(t, *offer, got.Offers[i])
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
	server, s, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
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
		offer, _ := s.GetByID(uint(i+1), nil)
		assert.Equal(t, *offer, got.Offers[i])
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
	server, _, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(got.Offers))
	u.CleanDB(DB)
}

// ----------------------
// Get users offers tests
// ----------------------
func TestGetUsersOffers_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestGetUsersOffers_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	token := "invalid"
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestGetUsersOffers_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(got.Offers))
	u.CleanDB(DB)
}

func TestGetUsersOffers_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(got.Offers))
	assert.Equal(t, seedOffers[0].ID, got.Offers[0].ID)
	u.CleanDB(DB)
}

func TestGetUsersOffers_OneAuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(got.Offers))
	assert.Equal(t, seedOffers[0].ID, got.Offers[0].ID)
	u.CleanDB(DB)
}

func TestGetUsersOffers_AuctionsAndOffersCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, _, _ := newTestServer(db, seedOffers)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(got.Offers))
	u.CleanDB(DB)
}

// ----------------------
// Get offer by id tests
// ----------------------

func TestGetSaleOfferByID_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestGetSaleOfferByID_NonExistentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/2", nil, nil)
	assert.Equal(t, http.StatusNotFound, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
	u.CleanDB(DB)
}

func TestGetSaleOfferByID_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/-1", nil, nil)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)

	u.CleanDB(DB)
}

func TestGetSaleOfferByID_StringID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/abc", nil, nil)
	assert.Equal(t, http.StatusBadRequest, receivedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
	u.CleanDB(DB)
}

func TestGetByID_RegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}

	server, _, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveSaleOfferDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, seedOffers[0].ID, got.ID)
	u.CleanDB(DB)
}

func TestGetByID_AuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}

	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestGetMyOffers_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	token := "invalid"
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestGetMyOffers_NoOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
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
	server, s, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	offer, _ := s.GetByID(uint(1), &user.ID)
	assert.Equal(t, *offer, got.Offers[0])
	u.CleanDB(DB)
}

func TestGetMyOffers_OneAuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	server, s, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	offer, _ := s.GetByID(uint(1), &user.ID)
	assert.Equal(t, *offer, got.Offers[0])
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
	server, s, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		offer, _ := s.GetByID(uint(i+1), &user.ID)
		assert.Equal(t, *offer, got.Offers[i])
	}
	u.CleanDB(DB)
}

func TestGetMyOffers_OtherUserOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(got.Offers))
	u.CleanDB(DB)
}

// ----------------------
// Get liked offers tests
// ----------------------

func TestGetLikedOffers_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/liked-offers", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestGetLikedOffers_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	token := "invalid"
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/liked-offers", nil, &token)
	assert.Equal(t, http.StatusForbidden, receivedStatus)
	u.CleanDB(DB)
}

func TestGetLikedOffers_NoOffersLiked(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedOffers []models.SaleOffer
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/liked-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got.Offers)
	u.CleanDB(DB)
}

func TestGetLikedOffers_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, s, _, _ := newTestServer(db, seedOffers)
	user := USERS[1] // Assuming user[1] is the one who likes the offer
	token, _ := u.GetValidToken(user.ID, user.Email)
	likedOfferRepo := liked_offer.NewLikedOfferRepository(db)
	likedOfferRepo.Create(&models.LikedOffer{
		UserID:  user.ID,
		OfferID: seedOffers[0].ID,
	})
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/liked-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(got.Offers))
	offer, _ := s.GetByID(uint(1), &user.ID)
	assert.Equal(t, *offer, got.Offers[0])
	u.CleanDB(DB)
}

func TestGetLikedOffers_OneAuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createAuctionSaleOffer(1)}
	server, s, _, _ := newTestServer(db, seedOffers)
	user := USERS[1] // Assuming user[1] is the one who likes the offer
	token, _ := u.GetValidToken(user.ID, user.Email)
	likedOfferRepo := liked_offer.NewLikedOfferRepository(db)
	likedOfferRepo.Create(&models.LikedOffer{
		UserID:  user.ID,
		OfferID: seedOffers[0].ID,
	})
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/liked-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(got.Offers))
	offer, _ := s.GetByID(uint(1), &user.ID)
	assert.Equal(t, *offer, got.Offers[0])
	u.CleanDB(DB)
}

func TestGetLikedOffers_AuctionsAndOffersCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, s, _, _ := newTestServer(db, seedOffers)
	user := USERS[1] // Assuming user[1] is the one who likes the offers
	token, _ := u.GetValidToken(user.ID, user.Email)
	likedOfferRepo := liked_offer.NewLikedOfferRepository(db)
	for _, offer := range seedOffers {
		likedOfferRepo.Create(&models.LikedOffer{
			UserID:  user.ID,
			OfferID: offer.ID,
		})
	}
	filterRequest := sale_offer.NewOfferFilterRequest()
	filterRequest.PagRequest = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filterRequest)
	assert.NoError(t, err)
	response, receivedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/liked-offers", body, &token)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		offer, _ := s.GetByID(uint(i+1), &user.ID)
		assert.Equal(t, *offer, got.Offers[i])
	}
	u.CleanDB(DB)
}

// ----------
// Like tests
// ----------

func TestLikeOffer_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestLikeOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	_, receivedStatus := u.PerformRequest(server, http.MethodPost, "/favourite/like/1", nil, &token)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
	_, receivedStatus := u.PerformRequest(server, http.MethodDelete, "/favourite/dislike/1", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, receivedStatus)
	u.CleanDB(DB)
}

func TestDislikeOffer_OfferNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []models.SaleOffer{*createOffer(1)}
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
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
	server, _, _, _ := newTestServer(db, seedOffers)
	response, receivedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/order-keys", nil, nil)
	assert.Equal(t, http.StatusOK, receivedStatus)
	var got map[string][]string
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.ElementsMatch(t, got["order_keys"], sale_offer.GetKeysFromMap(sale_offer.OrderKeysMap))
	u.CleanDB(DB)
}
