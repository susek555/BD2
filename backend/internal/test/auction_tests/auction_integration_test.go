package auction_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(manufacturers []models.Manufacturer, models_ []models.Model, cars []models.Car, saleOffers []models.SaleOffer, auctions []models.Auction, users []models.User) (auction.AuctionServiceInterface, error) {
	dsn := "host=localhost user=bd2_user password=bd2_password dbname=bd2_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("TRUNCATE TABLE auctions, sale_offers, cars, models, manufacturers, companies, people, users RESTART IDENTITY CASCADE")
	for _, user := range users {
		err = db.Create(&user).Error
		if err != nil {
			return nil, err
		}
	}
	for _, manufacturer := range manufacturers {
		err = db.Create(&manufacturer).Error
		if err != nil {
			return nil, err
		}
	}
	for _, model := range models_ {
		err = db.Create(&model).Error
		if err != nil {
			return nil, err
		}
	}
	for _, car := range cars {
		err = db.Create(&car).Error
		if err != nil {
			return nil, err
		}
	}
	for _, saleOffer := range saleOffers {
		err = db.Create(&saleOffer).Error
		if err != nil {
			return nil, err
		}
	}
	for _, auction_ := range auctions {
		err = db.Create(&auction_).Error
		if err != nil {
			return nil, err
		}
	}

	repo := auction.NewAuctionRepository(db)
	bidRepo := bid.NewBidRepository(db)
	likedOfferRepo := liked_offer.NewLikedOfferRepository(db)
	saleOfferService := sale_offer.NewSaleOfferService(
		sale_offer.NewSaleOfferRepository(db),
		manufacturer.NewManufacturerRepository(db),
		model.NewModelRepository(db),
		image.NewImageRepository(db),
		sale_offer.NewAccessEvaluator(bidRepo, likedOfferRepo),
	)
	service := auction.NewAuctionService(repo, saleOfferService.(*sale_offer.SaleOfferService))
	return service, nil
}

func newTestServer(seedManufacturers []models.Manufacturer, seedModels []models.Model, seedCars []models.Car, seedSaleOffers []models.SaleOffer, seedAuctions []models.Auction, seedUsers []models.User) (*gin.Engine, auction.AuctionServiceInterface, error) {
	service, err := setupDB(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	if err != nil {
		return nil, nil, err
	}
	verifier := jwt.NewJWTVerifier("secret")
	r := gin.Default()
	ms := new(mocks.SchedulerInterface)
	ms.
		On("AddAuction",
			mock.AnythingOfType("string"),    // ID aukcji
			mock.AnythingOfType("time.Time"), // termin zakończenia
		).
		Return(nil)
	mh := new(mocks.HubInterface)
	mh.
		On("SubscribeUser",
			mock.AnythingOfType("string"), // user ID
			mock.AnythingOfType("string"), // auction ID
		).
		Return(nil)
	mn := new(mocks.NotificationServiceInterface)
	auctionHandler := auction.NewHandler(service, ms, mh, mn)
	auctionRoutes := r.Group("/auction")
	auctionRoutes.GET("/", auctionHandler.GetAllAuctions)
	auctionRoutes.GET("/:id", auctionHandler.GetAuctionById)
	auctionRoutes.POST("/", middleware.Authenticate(verifier), auctionHandler.CreateAuction)
	auctionRoutes.PUT("/", middleware.Authenticate(verifier), auctionHandler.UpdateAuction)
	auctionRoutes.DELETE("/:id", middleware.Authenticate(verifier), auctionHandler.DeleteAuctionById)
	return r, service, nil
}

func getValidToken(userId uint, email string) (string, error) {
	secret := []byte("secret")
	return jwt.GenerateToken(email, int64(userId), secret, time.Now().Add(1*time.Hour))
}

func TestCreateAuctionNoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedManufacturers []models.Manufacturer
	var seedModels []models.Model
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	var seedUsers []models.User
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	wantStatus := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodPost, "/auction/", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", got["message"])
}

func TestCreateAuctionInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedManufacturers []models.Manufacturer
	var seedModels []models.Model
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	var seedUsers []models.User
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	wantStatus := http.StatusForbidden
	req := httptest.NewRequest(http.MethodPost, "/auction/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "forbidden", got["message"])
}
func TestCreateAuctionSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []models.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []models.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	seedUsers := []models.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	token, err := getValidToken(uint(1), seedUsers[0].Email)
	assert.NoError(t, err)
	wantStatus := http.StatusCreated
	auctionInput := auction.CreateAuctionDTO{
		CreateSaleOfferDTO: sale_offer.CreateSaleOfferDTO{
			UserID:             1,
			Description:        "Test auction",
			Price:              10000,
			Margin:             10,
			Vin:                "1HGCM82633A123456",
			ProductionYear:     2020,
			Mileage:            10000,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   "2023-10-01",
			Color:              enums.ORANGE,
			FuelType:           enums.PETROL,
			Transmission:       enums.MANUAL,
			NumberOfGears:      6,
			Drive:              enums.FWD,
			ManufacturerName:   "Toyota",
			ModelName:          "Corolla",
		},
		DateEnd:     "15:04 2026-02-01",
		BuyNowPrice: 12000,
	}
	auctionInputJSON, err := json.Marshal(auctionInput)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(string(auctionInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got auction.RetrieveAuctionDTO
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, auctionInput.BuyNowPrice, got.BuyNowPrice)
	assert.Equal(t, auctionInput.DateEnd, got.DateEnd)
	assert.Equal(t, "herakles", got.Username)
	assert.Equal(t, auctionInput.CreateSaleOfferDTO.ProductionYear, got.ProductionYear)
	assert.Equal(t, auctionInput.CreateSaleOfferDTO.Mileage, got.Mileage)
	assert.Equal(t, auctionInput.CreateSaleOfferDTO.Price, got.Price)
	assert.Equal(t, "Toyota Corolla", got.Name)
}

func TestCreateAuctionInvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []models.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []models.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	seedUsers := []models.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	token, err := getValidToken(uint(1), seedUsers[0].Email)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	auctionInput := auction.CreateAuctionDTO{
		CreateSaleOfferDTO: sale_offer.CreateSaleOfferDTO{
			UserID:             1,
			Description:        "Test auction",
			Price:              10000,
			Margin:             10,
			Vin:                "1HGCM82633A123456",
			ProductionYear:     2020,
			Mileage:            10000,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   "2023-10-01",
			Color:              enums.ORANGE,
			FuelType:           enums.PETROL,
			Transmission:       enums.MANUAL,
			NumberOfGears:      6,
			Drive:              enums.FWD,
			ManufacturerName:   "Toyota",
			ModelName:          "Corolla",
		},
		DateEnd:     "15:04 2023-01-01",
		BuyNowPrice: 12000,
	}
	auctionInputJSON, err := json.Marshal(auctionInput)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(string(auctionInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "date end must be in the future", got["error_description"])
}

func TestCreateAuctionInvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []models.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []models.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	seedUsers := []models.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	token, err := getValidToken(uint(1), seedUsers[0].Email)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	auctionInput := auction.CreateAuctionDTO{
		CreateSaleOfferDTO: sale_offer.CreateSaleOfferDTO{
			UserID:             1,
			Description:        "Test auction",
			Price:              10000,
			Margin:             10,
			Vin:                "1HGCM82633A123456",
			ProductionYear:     2020,
			Mileage:            10000,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   "2023-10-01",
			Color:              enums.ORANGE,
			FuelType:           enums.PETROL,
			Transmission:       enums.MANUAL,
			NumberOfGears:      6,
			Drive:              enums.FWD,
			ManufacturerName:   "Toyota",
			ModelName:          "Corolla",
		},
		DateEnd:     "25:04 02/01/2025",
		BuyNowPrice: 12000,
	}
	auctionInputJSON, err := json.Marshal(auctionInput)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(string(auctionInputJSON)))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "parsing time \"25:04 02/01/2025\": hour out of range", got["error_description"])
}

func TestCreateAuctionZeroBuyNowPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []models.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []models.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	seedUsers := []models.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	token, err := getValidToken(uint(1), seedUsers[0].Email)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	auctionInput := `
	{
		"vin": "1HGCM82633A123456",
		"production_year": 2020,
		"mileage": 10000,
		"number_of_doors": 4,
		"number_of_seats": 5,
		"engine_power": 150,
		"engine_capacity": 2000,
		"registration_number": "ABC123",
		"registration_date": "2023-10-01",
		"color": "Orange",
		"fuel_type": "Petrol",
		"transmission": "Manual",
		"number_of_gears": 6,
		"drive": "FWD",
		"model_id": 1,
		"date_end": "15:04 02/11/2025",
		"buy_now_price": 0
	}
	`
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(auctionInput))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "some fields are missing - ensure that all required fields are present", got["error_description"])
}

func TestCreateAuctionBuyNowPriceLessThanOfferPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []models.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []models.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	seedUsers := []models.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	token, err := getValidToken(uint(1), seedUsers[0].Email)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	auctionInput := `
	{
		"vin": "1HGCM82633A123456",
		"production_year": 2020,
		"mileage": 10000,
		"number_of_doors": 4,
		"number_of_seats": 5,
		"engine_power": 150,
		"price": 10000,
		"engine_capacity": 2000,
		"registration_number": "ABC123",
		"registration_date": "2023-10-01",
		"color": "Orange",
		"fuel_type": "Petrol",
		"transmission": "Manual",
		"number_of_gears": 6,
		"drive": "FWD",
		"model_id": 1,
		"date_end": "15:04 02/11/2025",
		"buy_now_price": 5000
	}
	`
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(auctionInput))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "some fields are missing - ensure that all required fields are present", got["error_description"])
}

func TestCreateAuctionBuyNowPriceNegative(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []models.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []models.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	var seedCars []models.Car
	var seedSaleOffers []models.SaleOffer
	var seedAuctions []models.Auction
	seedUsers := []models.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	server, _, err := newTestServer(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	assert.NoError(t, err)
	token, err := getValidToken(uint(1), seedUsers[0].Email)
	assert.NoError(t, err)
	wantStatus := http.StatusBadRequest
	auctionInput := `
	{
		"vin": "1HGCM82633A123456",
		"production_year": 2020,
		"mileage": 10000,
		"number_of_doors": 4,
		"number_of_seats": 5,
		"engine_power": 150,
		"engine_capacity": 2000,
		"registration_number": "ABC123",
		"registration_date": "2023-10-01",
		"color": "Orange",
		"fuel_type": "Petrol",
		"transmission": "Manual",
		"number_of_gears": 6,
		"drive": "FWD",
		"model_id": 1,
		"date_end": "15:04 02/11/2025",
		"buy_now_price": -100
	}
	`
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(auctionInput))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "json: cannot unmarshal number -100 into Go struct field CreateAuctionDTO.buy_now_price of type uint", got["error_description"])
}
