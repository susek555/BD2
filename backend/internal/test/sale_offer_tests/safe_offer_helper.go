package sale_offer_tests

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ---------
// Constants
// ---------

var MANUFACTURERS = []models.Manufacturer{
	{ID: 1, Name: "Audi"},
	{ID: 2, Name: "BMW"},
	{ID: 3, Name: "Opel"},
	{ID: 4, Name: "Toyota"},
	{ID: 5, Name: "Skoda"},
}

var MODELS = []models.Model{
	{ID: 1, Name: "A3", ManufacturerID: 1},
	{ID: 2, Name: "M3", ManufacturerID: 2},
	{ID: 3, Name: "Astra", ManufacturerID: 3},
	{ID: 4, Name: "Supra", ManufacturerID: 4},
	{ID: 5, Name: "Octavia", ManufacturerID: 5},
}

var USERS = []models.User{
	{ID: 1, Username: "john", Email: "john@example.com", Selector: "P"},
	{ID: 2, Username: "jane", Email: "jane@example.com", Selector: "P"},
}

// ------
// Setup
// ------

func setupDB() (*gorm.DB, error) {
	dsn := "host=localhost user=bd2_user password=bd2_password dbname=bd2_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("TRUNCATE TABLE bids, liked_offers, sale_offers, auctions, cars, models, manufacturers, companies, people, users RESTART IDENTITY CASCADE")
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

func getRepositoryWithSaleOffers(db *gorm.DB, offers []models.SaleOffer) sale_offer.SaleOfferRepositoryInterface {
	repo := sale_offer.NewSaleOfferRepository(db)
	for _, offer := range offers {
		repo.Create(&offer)
	}
	return repo
}

func newTestServer(db *gorm.DB, seedOffers []models.SaleOffer) (*gin.Engine, sale_offer.SaleOfferServiceInterface) {
	verifier := jwt.NewJWTVerifier(u.JWTSECRET)
	saleOfferRepo := getRepositoryWithSaleOffers(db, seedOffers)
	manufacturerRepo := manufacturer.NewManufacturerRepository(db)
	modelRepo := model.NewModelRepository(db)
	likedOfferRepository := liked_offer.NewLikedOfferRepository(db)
	bidRepository := bid.NewBidRepository(db)
	imageRepo := image.NewImageRepository(db)
	saleOfferService := sale_offer.NewSaleOfferService(saleOfferRepo, manufacturerRepo, modelRepo, bidRepository, imageRepo, likedOfferRepository)
	likedOfferService := liked_offer.NewLikedOfferService(likedOfferRepository, saleOfferRepo)
	likedOfferHandler := liked_offer.NewHandler(likedOfferService)
	saleOfferHandler := sale_offer.NewHandler(saleOfferService)
	r := gin.Default()
	saleOfferRoutes := r.Group("/sale-offer")
	{
		saleOfferRoutes.POST("/", middleware.Authenticate(verifier), saleOfferHandler.CreateSaleOffer)
		saleOfferRoutes.PUT("/", middleware.Authenticate(verifier), saleOfferHandler.UpdateSaleOffer)
		saleOfferRoutes.POST("/filtered", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetFilteredSaleOffers)
		saleOfferRoutes.POST("/my-offers", middleware.Authenticate(verifier), saleOfferHandler.GetMySaleOffers)
		saleOfferRoutes.GET("/id/:id", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetSaleOfferByID)
		saleOfferRoutes.POST("/like/:id", middleware.Authenticate(verifier), likedOfferHandler.LikeOffer)
		saleOfferRoutes.DELETE("/dislike/:id", middleware.Authenticate(verifier), likedOfferHandler.DislikeOffer)
		saleOfferRoutes.GET("/offer-types", saleOfferHandler.GetSaleOfferTypes)
		saleOfferRoutes.GET("/order-keys", saleOfferHandler.GetOrderKeys)
	}
	return r, saleOfferService
}

// ------------
// Basic models
// ------------

func createOffer(id uint) *models.SaleOffer {
	c := &models.Car{
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
		Model: models.Model{
			ID:             1,
			Name:           MODELS[0].Name,
			ManufacturerID: MODELS[0].ManufacturerID,
			Manufacturer:   MANUFACTURERS[0],
		},
	}
	offer := &models.SaleOffer{
		ID:          id,
		UserID:      1,
		User:        &USERS[0],
		Description: "offer",
		Price:       1000,
		Margin:      models.LOW_MARGIN,
		DateOfIssue: time.Now(),
		Car:         c,
	}
	return offer
}

func withCarField(opt u.Option[models.Car]) u.Option[models.SaleOffer] {
	return func(offer *models.SaleOffer) {
		if offer.Car == nil {
			offer.Car = &models.Car{}
		}
		opt(offer.Car)
	}
}

func withAuctionField(opt u.Option[models.Auction]) u.Option[models.SaleOffer] {
	return func(offer *models.SaleOffer) {
		if offer.Auction == nil {
			offer.Auction = &models.Auction{}
		}
		opt(offer.Auction)
	}
}

func createAuctionSaleOffer(id uint) *models.SaleOffer {
	offer := *u.Build(createOffer(id),
		withAuctionField(u.WithField[models.Auction]("OfferID", id)),
		withAuctionField(u.WithField[models.Auction]("DateEnd", time.Now())),
		withAuctionField(u.WithField[models.Auction]("BuyNowPrice", uint(0))))
	return &offer
}

func createSaleOfferDTO() *sale_offer.CreateSaleOfferDTO {
	return &sale_offer.CreateSaleOfferDTO{
		UserID:             1,
		Description:        "offer",
		Price:              1000,
		Margin:             models.LOW_MARGIN,
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
		Manufacturer:       "Audi",
		Model:              "A3",
	}
}

func doSaleOfferAndRetrieveSaleOfferDTOsMatch(offer models.SaleOffer, dto sale_offer.RetrieveSaleOfferDTO, s sale_offer.SaleOfferServiceInterface, userID *uint) bool {
	var likedCondition bool
	var bidCondition bool
	condition := offer.ID == dto.ID &&
		offer.User.Username == dto.Username &&
		offer.Car.Model.Manufacturer.Name+" "+offer.Car.Model.Name == dto.Name &&
		offer.Price == dto.Price &&
		offer.Car.Mileage == dto.Mileage &&
		offer.Car.ProductionYear == dto.ProductionYear &&
		offer.Car.Color == dto.Color &&
		(offer.Auction != nil) == dto.IsAuction
	if userID == nil {
		likedCondition = false
		bidCondition = false
	} else {
		likedCondition = s.IsOfferLikedByUser(offer.ID, userID)
		bidCondition, _ = s.CanBeModifiedByUser(offer.ID, userID)
	}
	return condition && bidCondition == dto.CanModify && likedCondition == dto.IsLiked
}

func wasEntityAddedToDB[T any](db *gorm.DB, id uint) bool {
	repo := generic.GetGormRepository[T](db)
	_, err := repo.GetById(id)
	return err == nil
}
