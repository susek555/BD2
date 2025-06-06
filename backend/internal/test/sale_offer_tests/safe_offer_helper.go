package sale_offer_tests

import (
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/purchase"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
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
		offer.Status = enums.PUBLISHED
		repo.Create(&offer)
	}
	return repo
}

func newTestServer(db *gorm.DB, seedOffers []models.SaleOffer) (*gin.Engine, sale_offer.SaleOfferServiceInterface, sale_offer.OfferAccessEvaluatorInterface, image.ImageServiceInterface) {
	os.Setenv("CLOUDINARY_URL", "cloudinary://1234567890:abcdefghijklmnopqrstuvwxyz@my‐cloud")

	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, _ := cloudinary.NewFromURL(cloudinaryURL)
	verifier := jwt.NewJWTVerifier(u.JWTSECRET)
	saleOfferRepo := getRepositoryWithSaleOffers(db, seedOffers)
	manufacturerRepo := manufacturer.NewManufacturerRepository(db)
	modelRepo := model.NewModelRepository(db)
	likedOfferRepository := liked_offer.NewLikedOfferRepository(db)
	bidRepository := bid.NewBidRepository(db)
	imageRepo := image.NewImageRepository(db)
	accessEvaluator := sale_offer.NewAccessEvaluator(bidRepository, likedOfferRepository)
	purchaseCreator := purchase.NewPurchaseRepository(db)
	saleOfferService := sale_offer.NewSaleOfferService(saleOfferRepo, manufacturerRepo, modelRepo, imageRepo, accessEvaluator, purchaseCreator)
	likedOfferService := liked_offer.NewLikedOfferService(likedOfferRepository, saleOfferRepo)
	imageService := image.NewImageService(imageRepo, &image.ImageBucket{CloudinaryClient: cld}, saleOfferRepo)
	imageHandler := image.NewHandler(imageService, saleOfferService)
	mh := new(mocks.HubInterface)
	mh.On("SubscribeUser", mock.Anything, mock.Anything).Return()
	mh.On("UnsubscribeUser", mock.Anything, mock.Anything).Return()
	likedOfferHandler := liked_offer.NewHandler(likedOfferService, mh)
	mn := new(mocks.NotificationServiceInterface)
	saleOfferHandler := sale_offer.NewHandler(saleOfferService, mh, mn)
	r := gin.Default()
	saleOfferRoutes := r.Group("/sale-offer")
	{
		saleOfferRoutes.POST("/", middleware.Authenticate(verifier), saleOfferHandler.CreateSaleOffer)
		saleOfferRoutes.PUT("/", middleware.Authenticate(verifier), saleOfferHandler.UpdateSaleOffer)
		saleOfferRoutes.POST("/filtered", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetFilteredSaleOffers)
		saleOfferRoutes.POST("/my-offers", middleware.Authenticate(verifier), saleOfferHandler.GetMySaleOffers)
		saleOfferRoutes.GET("/id/:id", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetDetailedSaleOfferByID)
		saleOfferRoutes.GET("/offer-types", saleOfferHandler.GetSaleOfferTypes)
		saleOfferRoutes.GET("/order-keys", saleOfferHandler.GetOrderKeys)
		saleOfferRoutes.PUT("/publish/:id", middleware.Authenticate(verifier), saleOfferHandler.PublishSaleOffer)
	}
	favourtieRoutes := r.Group("/favourite")
	{
		favourtieRoutes.POST("/like/:id", middleware.Authenticate(verifier), likedOfferHandler.LikeOffer)
		favourtieRoutes.DELETE("/dislike/:id", middleware.Authenticate(verifier), likedOfferHandler.DislikeOffer)
	}
	imageRoutes := r.Group("/image")
	{
		imageRoutes.PATCH("/:id", middleware.Authenticate(verifier), imageHandler.UploadImages)
		imageRoutes.DELETE("/", middleware.Authenticate(verifier), imageHandler.DeleteImage)
	}
	return r, saleOfferService, accessEvaluator, imageService
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
		Color:              enums.BLACK,
		FuelType:           enums.PETROL,
		Transmission:       enums.MANUAL,
		NumberOfGears:      6,
		Drive:              enums.FWD,
		ModelID:            1,
		Model: &models.Model{
			ID:             1,
			Name:           MODELS[0].Name,
			ManufacturerID: MODELS[0].ManufacturerID,
			Manufacturer:   &MANUFACTURERS[0],
		},
	}
	offer := &models.SaleOffer{
		ID:          id,
		UserID:      1,
		User:        &USERS[0],
		Description: "offer",
		Price:       1000,
		Margin:      enums.LOW_MARGIN,
		DateOfIssue: time.Now(),
		Car:         c,
		Status:      enums.PUBLISHED,
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
		Margin:             enums.LOW_MARGIN,
		Vin:                "vin",
		ProductionYear:     2025,
		Mileage:            1000,
		NumberOfDoors:      4,
		NumberOfSeats:      5,
		EnginePower:        100,
		EngineCapacity:     2000,
		RegistrationNumber: "default",
		RegistrationDate:   time.Now().Format("2006-01-02"),
		Color:              enums.BLACK,
		FuelType:           enums.PETROL,
		Transmission:       enums.MANUAL,
		NumberOfGears:      6,
		Drive:              enums.FWD,
		ManufacturerName:   "Audi",
		ModelName:          "A3",
	}
}

func setOffersStatusToPublished(db *gorm.DB) {
	db.Model(&models.SaleOffer{}).Update("status", enums.PUBLISHED)
}

func wasEntityAddedToDB[T any](db *gorm.DB, id uint) bool {
	repo := generic.GetGormRepository[T](db)
	_, err := repo.GetByID(id)
	return err == nil
}
