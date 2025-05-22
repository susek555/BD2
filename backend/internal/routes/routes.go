package routes

import (
	"log"
	"os"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/liked_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
)

func RegisterRoutes(router *gin.Engine) {
	registerAuthRoutes(router)
	registerUserRoutes(router)
	registerReviewRoutes(router)
	registerCarRoutes(router)
	registerSaleOfferRoutes(router)
	registerAuctionRoutes(router)
	registerBidRoutes(router)
}

func initializeVerifier() (*jwt.JWTVerifier, []byte) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	verifier := jwt.NewJWTVerifier(secret)
	jwtKey := []byte(secret)
	return verifier, jwtKey
}

func registerUserRoutes(router *gin.Engine) {
	userRepo := user.NewUserRepository(initializers.DB)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	verifier, _ := initializeVerifier()
	userRoutes := router.Group("/users")
	{
		userRoutes.PUT("/", middleware.Authenticate(verifier), userHandler.UpdateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/id/:id", userHandler.GetUserById)
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)
		userRoutes.DELETE("/:id", middleware.Authenticate(verifier), userHandler.DeleteUser)
	}
}

func registerAuthRoutes(router *gin.Engine) {
	verifier, jwtKey := initializeVerifier()
	authService := auth.NewService(initializers.DB, jwtKey)
	authHandler := auth.NewHandler(authService)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
	}
	router.POST("/logout", middleware.Authenticate(verifier), authHandler.Logout)
}

func registerReviewRoutes(router *gin.Engine) {
	verifier, _ := initializeVerifier()
	reviewRepo := review.NewReviewRepository(initializers.DB)
	reviewService := review.NewReviewService(reviewRepo)
	reviewHandler := review.NewHandler(reviewService)
	reviewRoutes := router.Group("/review")
	reviewRoutes.GET("/", reviewHandler.GetAllReviews)
	reviewRoutes.GET("/:id", reviewHandler.GetReviewById)
	reviewRoutes.POST("/", middleware.Authenticate(verifier), reviewHandler.CreateReview)
	reviewRoutes.PUT("/", middleware.Authenticate(verifier), reviewHandler.UpdateReview)
	reviewRoutes.DELETE("/:id", middleware.Authenticate(verifier), reviewHandler.DeleteReview)
	reviewRoutes.GET("/reviewer/:id", reviewHandler.GetReviewsByReviewerId)
	reviewRoutes.GET("/reviewee/:id", reviewHandler.GetReviewsByRevieweeId)
	reviewRoutes.GET("/reviewer/reviewee/:reviewerId/:revieweeId", reviewHandler.GetReviewsByReviewerIdAndRevieweeId)
}

func registerCarRoutes(router *gin.Engine) {
	modelRepo := model.NewModelRepository(initializers.DB)
	modelService := model.NewModelService(modelRepo)
	modelHandler := model.NewModelHandler(modelService)
	manufacturerRepo := manufacturer.NewManufacturerRepository(initializers.DB)
	manufacturerService := manufacturer.NewManufacturerService(manufacturerRepo)
	manufacturerHandler := manufacturer.NewManufacturerHandler(manufacturerService)
	carParamHandler := car_params.NewHandler()
	carService := car.NewCarService(manufacturerRepo, modelRepo)
	carHandler := car.NewCarHandler(carService)
	carRoutes := router.Group("/car")
	{
		carRoutes.GET("/colors", carParamHandler.GetPossibleColors)
		carRoutes.GET("/transmissions", carParamHandler.GetPossibleTransmissions)
		carRoutes.GET("/fuel-types", carParamHandler.GetPossibleFuelTypes)
		carRoutes.GET("/drives", carParamHandler.GetPossibleDrives)
		carRoutes.GET("/manufactures", manufacturerHandler.GetAllManufactures)
		carRoutes.GET("/models/id/:id", modelHandler.GetModelsByManufacturerID)
		carRoutes.GET("/models/name/:name", modelHandler.GetModelsByManufacturerName)
		carRoutes.GET("/manufacturer-model-map", carHandler.GetManufacturersModelsMap)
	}
}

func registerSaleOfferRoutes(router *gin.Engine) {
	verifier, _ := initializeVerifier()
	saleOfferRepo := sale_offer.NewSaleOfferRepository(initializers.DB)
	manufacturerRepo := manufacturer.NewManufacturerRepository(initializers.DB)
	likedOfferRepository := liked_offer.NewLikedOfferRepository(initializers.DB)
	bidRepository := bid.NewBidRepository(initializers.DB)
	saleOfferService := sale_offer.NewSaleOfferService(saleOfferRepo, manufacturerRepo, likedOfferRepository, bidRepository)
	saleOfferHandler := sale_offer.NewSaleOfferHandler(saleOfferService)
	saleOfferRoutes := router.Group("/sale-offer")
	{
		saleOfferRoutes.POST("/", middleware.Authenticate(verifier), saleOfferHandler.CreateSaleOffer)
		saleOfferRoutes.POST("/my-offers", middleware.Authenticate(verifier), saleOfferHandler.GetMySaleOffers)
		saleOfferRoutes.POST("/filtered", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetFilteredSaleOffers)
		saleOfferRoutes.GET("/id/:id", middleware.OptionalAuthenticate(verifier), saleOfferHandler.GetSaleOfferByID)
		saleOfferRoutes.POST("/like/:id", middleware.Authenticate(verifier), saleOfferHandler.LikeOffer)
		saleOfferRoutes.DELETE("/dislike/:id", middleware.Authenticate(verifier), saleOfferHandler.DislikeOffer)
		saleOfferRoutes.GET("/offer-types", saleOfferHandler.GetSaleOfferTypes)
		saleOfferRoutes.GET("/order-keys", saleOfferHandler.GetOrderKeys)
	}
}

func registerAuctionRoutes(router *gin.Engine) {
	auctionRepo := auction.NewAuctionRepository(initializers.DB)
	auctionService := auction.NewAuctionService(auctionRepo)
	auctionHandler := auction.NewHandler(auctionService)
	auctionRoutes := router.Group("/auction")
	auctionRoutes.GET("/", auctionHandler.GetAllAuctions)
	auctionRoutes.GET("/:id", auctionHandler.GetAuctionById)
	auctionRoutes.POST("/", auctionHandler.CreateAuction)
	auctionRoutes.PUT("/", auctionHandler.UpdateAuction)
	auctionRoutes.DELETE("/:id", auctionHandler.DeleteAuctionById)

}

func registerBidRoutes(router *gin.Engine) {
	bidRepo := bid.NewBidRepository(initializers.DB)
	bidService := bid.NewBidService(bidRepo)
	bidHandler := bid.NewHandler(bidService)
	bidRoutes := router.Group("/bid")
	bidRoutes.POST("/", bidHandler.CreateBid)
	bidRoutes.GET("/", bidHandler.GetAllBids)
	bidRoutes.GET("/:id", bidHandler.GetBidByID)
	bidRoutes.GET("/bidder/:id", bidHandler.GetBidsByBidderId)
	bidRoutes.GET("/auction/:id", bidHandler.GetBidsByAuctionId)
	bidRoutes.GET("/highest/:id", bidHandler.GetHighestBid)
	bidRoutes.GET("/highest/auction/:auctionId/bidder/:bidderId", bidHandler.GetHighestBidByUserId)
}
