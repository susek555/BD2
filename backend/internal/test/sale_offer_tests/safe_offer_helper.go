package sale_offer_tests

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ---------
// Constants
// ---------

var MANUFACTURERS = []manufacturer.Manufacturer{
	{ID: 1, Name: "Audi"},
	{ID: 2, Name: "BMW"},
	{ID: 3, Name: "Opel"},
	{ID: 4, Name: "Toyota"},
	{ID: 5, Name: "Skoda"},
}

var MODELS = []model.Model{
	{ID: 1, Name: "A3", ManufacturerID: 1},
	{ID: 2, Name: "M3", ManufacturerID: 2},
	{ID: 3, Name: "Astra", ManufacturerID: 3},
	{ID: 4, Name: "Supra", ManufacturerID: 4},
	{ID: 5, Name: "Octavia", ManufacturerID: 5},
}

var USERS = []user.User{
	{ID: 1, Username: "john", Email: "john@example.com", Selector: "P"},
	{ID: 2, Username: "jane", Email: "jane@example.com", Selector: "P"},
}

// ------
// Setup
// ------

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.Exec("PRAGMA foreign_keys = ON")
	db.AutoMigrate(
		&user.User{},
		&manufacturer.Manufacturer{},
		&model.Model{},
		&car.Car{},
		&sale_offer.Auction{},
		&sale_offer.SaleOffer{},
		&liked_offer.LikedOffer{},
	)
	if err != nil {
		return nil, err
	}
	if err := u.InsertRecordsIntoDB(db, MANUFACTURERS); err != nil {
		return nil, err
	}
	if err := u.InsertRecordsIntoDB(db, MODELS); err != nil {
		return nil, err
	}
	if err := u.InsertRecordsIntoDB(db, USERS); err != nil {
		return nil, err
	}
	return db, nil
}

func getRepositoryWithSaleOffers(db *gorm.DB, offers []sale_offer.SaleOffer) sale_offer.SaleOfferRepositoryInterface {
	repo := sale_offer.NewSaleOfferRepository(db)
	for _, offer := range offers {
		repo.Create(&offer)
	}
	return repo
}

func newTestServer(db *gorm.DB, seedOffers []sale_offer.SaleOffer) (*gin.Engine, sale_offer.SaleOfferServiceInterface, error) {
	verifier := jwt.NewJWTVerifier(u.JWTSECRET)
	saleOfferRepo := getRepositoryWithSaleOffers(db, seedOffers)
	manufacturerRepo := manufacturer.NewManufacturerRepository(db)
	likedOfferRepository := liked_offer.NewLikedOfferRepository(db)
	saleOfferService := sale_offer.NewSaleOfferService(saleOfferRepo, manufacturerRepo, likedOfferRepository)
	saleOfferHandler := sale_offer.NewSaleOfferHandler(saleOfferService)
	r := gin.Default()
	saleOfferRoutes := r.Group("/sale-offer")
	{
		saleOfferRoutes.POST("/", middleware.Authenticate(verifier), saleOfferHandler.CreateSaleOffer)
		saleOfferRoutes.POST("/filtered", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetFilteredSaleOffers)
		saleOfferRoutes.POST("/my-offers", middleware.Authenticate(verifier), saleOfferHandler.GetMySaleOffers)
		saleOfferRoutes.GET("/id/:id", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetSaleOfferByID)
		saleOfferRoutes.GET("/offer-types", saleOfferHandler.GetSaleOfferTypes)
		saleOfferRoutes.GET("/order-keys", saleOfferHandler.GetOrderKeys)
	}
	return r, saleOfferService, nil
}

// ------------
// Basic models
// ------------

func createOffer(id uint) *sale_offer.SaleOffer {
	c := &car.Car{
		OfferID:            id,
		Vin:                "vin",
		ProductionYear:     2025,
		Mileage:            1000,
		NumberOfDoors:      4,
		NumberOfSeats:      5,
		EnginePower:        100,
		EngineCapacity:     2000,
		RegistrationNumber: "default",
		RegistrationDate:   time.Now(),
		Color:              car_params.BLACK,
		FuelType:           car_params.PETROL,
		Transmission:       car_params.MANUAL,
		NumberOfGears:      6,
		Drive:              car_params.FWD,
		ModelID:            1,
		Model: model.Model{
			ID:             1,
			Name:           MODELS[0].Name,
			ManufacturerID: MODELS[0].ManufacturerID,
			Manufacturer:   MANUFACTURERS[0],
		},
	}
	offer := &sale_offer.SaleOffer{
		ID:          id,
		UserID:      1,
		User:        &USERS[0],
		Description: "offer",
		Price:       1000,
		Margin:      sale_offer.LOW_MARGIN,
		DateOfIssue: time.Now(),
		Car:         c,
	}
	return offer
}

func withCarField(opt u.Option[car.Car]) u.Option[sale_offer.SaleOffer] {
	return func(offer *sale_offer.SaleOffer) {
		if offer.Car == nil {
			offer.Car = &car.Car{}
		}
		opt(offer.Car)
	}
}

func withAuctionField(opt u.Option[sale_offer.Auction]) u.Option[sale_offer.SaleOffer] {
	return func(offer *sale_offer.SaleOffer) {
		if offer.Auction == nil {
			offer.Auction = &sale_offer.Auction{}
		}
		opt(offer.Auction)
	}
}

func createAuctionSaleOffer(id uint) *sale_offer.SaleOffer {
	offer := *u.Build(createOffer(id),
		withAuctionField(u.WithField[sale_offer.Auction]("OfferID", id)),
		withAuctionField(u.WithField[sale_offer.Auction]("DateEnd", time.Now())),
		withAuctionField(u.WithField[sale_offer.Auction]("BuyNowPrice", uint(0))))
	return &offer
}

func createSaleOfferDTO() *sale_offer.CreateSaleOfferDTO {
	return &sale_offer.CreateSaleOfferDTO{
		UserID:             1,
		Description:        "offer",
		Price:              1000,
		Margin:             sale_offer.LOW_MARGIN,
		Vin:                "vin",
		ProductionYear:     2025,
		Mileage:            1000,
		NumberOfDoors:      4,
		NumberOfSeats:      5,
		EnginePower:        100,
		EngineCapacity:     2000,
		RegistrationNumber: "default",
		RegistrationDate:   time.Now().Format("2006-01-02"),
		Color:              car_params.BLACK,
		FuelType:           car_params.PETROL,
		Transmission:       car_params.MANUAL,
		NumberOfGears:      6,
		Drive:              car_params.FWD,
		ModelID:            1,
	}
}

func doSaleOfferAndRetrieveSaleOfferDTOsMatch(offer sale_offer.SaleOffer, dto sale_offer.RetrieveSaleOfferDTO, s sale_offer.SaleOfferServiceInterface) bool {
	return offer.ID == dto.ID &&
		offer.User.Username == dto.Username &&
		offer.Car.Model.Manufacturer.Name+" "+offer.Car.Model.Name == dto.Name &&
		offer.Price == dto.Price &&
		offer.Car.Mileage == dto.Mileage &&
		offer.Car.ProductionYear == dto.ProductionYear &&
		offer.Car.Color == dto.Color &&
		(offer.Auction != nil) == dto.IsAuction &&
		(s.IsOfferLikedByUser(offer.ID, offer.User.ID)) == dto.IsLiked
}

func wasEntityAddedToDB[T any](db *gorm.DB, id uint) bool {
	repo := generic.GetGormRepository[T](db)
	_, err := repo.GetById(id)
	return err == nil
}
