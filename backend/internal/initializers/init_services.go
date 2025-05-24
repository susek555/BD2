package initializers

import (
	"log"
	"os"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

var AuctionService auction.AuctionServiceInterface
var AuthService auth.AuthServiceInterface
var BidService bid.BidServiceInterface
var CarService car.CarServiceInterface
var ImageService image.ImageServiceInterface
var ManufacturerService manufacturer.ManufacturerServiceInterface
var ModelService model.ModelServiceInterface
var NotificationService notification.NotificationServiceInterface
var RefreshTokenService refresh_token.RefreshTokenServiceInterface
var ReviewService review.ReviewServiceInterface
var SaleOfferService sale_offer.SaleOfferServiceInterface
var LikedOfferService liked_offer.LikedOfferServiceInterface
var UserService user.UserServiceInterface

func InitializeServices() {
	AuctionService = auction.NewAuctionService(AuctionRepo)
	RefreshTokenService = refresh_token.NewRefreshTokenService(RefreshTokenRepo)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	AuthService = auth.NewAuthService(UserRepo, RefreshTokenService, []byte(secret))
	BidService = bid.NewBidService(BidRepo)
	CarService = car.NewCarService(ManufacturerRepo, ModelRepo)
	ImageService = image.NewImageService(ImageRepo, ImageBucket)
	ManufacturerService = manufacturer.NewManufacturerService(ManufacturerRepo)
	ModelService = model.NewModelService(ModelRepo)
	NotificationService = notification.NewNotificationService(NotificationRepo)
	ReviewService = review.NewReviewService(ReviewRepo)
	SaleOfferService = sale_offer.NewSaleOfferService(SaleOfferRepo, ManufacturerRepo, BidRepo, ImageRepo, LikedOfferRepo)
	LikedOfferService = liked_offer.NewLikedOfferService(LikedOfferRepo, SaleOfferRepo)
	UserService = user.NewUserService(UserRepo)
}
