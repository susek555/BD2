package auction_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(manufacturers []manufacturer.Manufacturer, models []model.Model, cars []sale_offer.Car, saleOffers []sale_offer.SaleOffer, auctions []sale_offer.Auction, users []user.User) (auction.AuctionServiceInterface, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&user.User{},
		&user.Person{},
		user.Company{},
		&manufacturer.Manufacturer{},
		&model.Model{},
		&sale_offer.Car{},
		&sale_offer.SaleOffer{},
		&sale_offer.Auction{},
	)
	if err != nil {
		return nil, err
	}
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
	for _, model := range models {
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
	for _, auction := range auctions {
		err = db.Create(&auction).Error
		if err != nil {
			return nil, err
		}
	}

	repo := auction.NewAuctionRepository(db)
	service := auction.NewAuctionService(repo)
	return service, nil
}

func newTestServer(seedManufacturers []manufacturer.Manufacturer, seedModels []model.Model, seedCars []sale_offer.Car, seedSaleOffers []sale_offer.SaleOffer, seedAuctions []sale_offer.Auction, seedUsers []user.User) (*gin.Engine, auction.AuctionServiceInterface, error) {
	service, err := setupDB(seedManufacturers, seedModels, seedCars, seedSaleOffers, seedAuctions, seedUsers)
	if err != nil {
		return nil, nil, err
	}
	verifier := jwt.NewJWTVerifier("secret")
	r := gin.Default()
	auctionHandler := auction.NewHandler(service)
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
	seedManufacturers := []manufacturer.Manufacturer{}
	seedModels := []model.Model{}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{}
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
	seedManufacturers := []manufacturer.Manufacturer{}
	seedModels := []model.Model{}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{}
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
	seedManufacturers := []manufacturer.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []model.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
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
			Margin:             1000,
			DateOfIssue:        time.Now(),
			Vin:                "1HGCM82633A123456",
			ProductionYear:     2020,
			Mileage:            10000,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   "2023-10-01",
			Color:              car_params.ORANGE,
			FuelType:           car_params.PETROL,
			Transmission:       car_params.MANUAL,
			NumberOfGears:      6,
			Drive:              car_params.FWD,
			ModelID:            1,
		},
		DateEnd:     "15:04 02/01/2026",
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
	seedManufacturers := []manufacturer.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []model.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
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
			Margin:             1000,
			DateOfIssue:        time.Now(),
			Vin:                "1HGCM82633A123456",
			ProductionYear:     2020,
			Mileage:            10000,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   "2023-10-01",
			Color:              car_params.ORANGE,
			FuelType:           car_params.PETROL,
			Transmission:       car_params.MANUAL,
			NumberOfGears:      6,
			Drive:              car_params.FWD,
			ModelID:            1,
		},
		DateEnd:     "15:04 02/01/2025",
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
	seedManufacturers := []manufacturer.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []model.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
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
			Margin:             1000,
			DateOfIssue:        time.Now(),
			Vin:                "1HGCM82633A123456",
			ProductionYear:     2020,
			Mileage:            10000,
			NumberOfDoors:      4,
			NumberOfSeats:      5,
			EnginePower:        150,
			EngineCapacity:     2000,
			RegistrationNumber: "ABC123",
			RegistrationDate:   "2023-10-01",
			Color:              car_params.ORANGE,
			FuelType:           car_params.PETROL,
			Transmission:       car_params.MANUAL,
			NumberOfGears:      6,
			Drive:              car_params.FWD,
			ModelID:            1,
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
	seedManufacturers := []manufacturer.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []model.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
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
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(string(auctionInput)))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "buy now price must be greater than 0", got["error_description"])
}

func TestCreateAuctionBuyNowPriceLessThanOfferPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []manufacturer.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []model.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
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
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(string(auctionInput)))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, "buy now price must be greater than offer price", got["error_description"])
}

func TestCreateAuctionBuyNowPriceNegative(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedManufacturers := []manufacturer.Manufacturer{
		{
			ID:   1,
			Name: "Toyota",
		},
	}
	seedModels := []model.Model{
		{
			ID:             1,
			Name:           "Corolla",
			ManufacturerID: 1,
		},
	}
	seedCars := []sale_offer.Car{}
	seedSaleOffers := []sale_offer.SaleOffer{}
	seedAuctions := []sale_offer.Auction{}
	seedUsers := []user.User{
		{
			Username: "herakles",
			Email:    "herakles@gmail.com",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
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
	req := httptest.NewRequest(http.MethodPost, "/auction/", strings.NewReader(string(auctionInput)))
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
