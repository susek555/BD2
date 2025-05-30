package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/image"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

var AuctionHandler *auction.Handler
var AuthHandler *auth.Handler
var BidHandler *bid.Handler
var CarHandler *car.Handler
var ImageHandler *image.Handler
var ManufacturerHandler *manufacturer.Handler
var ModelHandler *model.Handler
var ReviewHandler *review.Handler
var SaleOfferHandler *sale_offer.Handler
var LikedOfferHandler *liked_offer.Handler
var UserHandler *user.Handler

func InitializeHandlers() {
	AuctionHandler = auction.NewHandler(AuctionService, Sched)
	AuthHandler = auth.NewHandler(AuthService)
	BidHandler = bid.NewHandler(BidService, RedisClient, Hub, NotificationService)
	CarHandler = car.NewHandler(CarService)
	ImageHandler = image.NewHandler(ImageService, SaleOfferService)
	ManufacturerHandler = manufacturer.NewHandler(ManufacturerService)
	ModelHandler = model.NewHandler(ModelService)
	ReviewHandler = review.NewHandler(ReviewService)
	SaleOfferHandler = sale_offer.NewHandler(SaleOfferService)
	LikedOfferHandler = liked_offer.NewHandler(LikedOfferService)
	UserHandler = user.NewHandler(UserService)
}
