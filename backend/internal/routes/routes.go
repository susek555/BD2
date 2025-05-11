package routes

import (
	"log"
	"os"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auction"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/car/car_params"
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
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
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
	reviewService := review.NewReviewService(initializers.DB)
	reviewHandler := review.NewHandler(reviewService)
	reviewRoutes := router.Group("/review")
	reviewRoutes.GET("/", reviewHandler.GetAllReviews)
	reviewRoutes.GET("/:id", reviewHandler.GetReviewById)
	reviewRoutes.POST("/", reviewHandler.CreateReview)
	reviewRoutes.PUT("/", reviewHandler.UpdateReview)
	reviewRoutes.DELETE("/:id", reviewHandler.DeleteReview)
	reviewRoutes.GET("/reviewer/:id", reviewHandler.GetReviewsByReviewerId)
	reviewRoutes.GET("/reviewee/:id", reviewHandler.GetReviewsByRevieweeId)
	reviewRoutes.GET("/reviewer/reviewee/:reviewerId/:revieweeId", reviewHandler.GetReviewsByReviewerIdAndRevieweeId)
}

func registerCarRoutes(router *gin.Engine) {
	modelRepo := model.NewModelRepository(initializers.DB)
	modelService := model.NewModelService(modelRepo)
	modelHandler := model.NewHandler(modelService)
	manufacturerRepo := manufacturer.NewManufacturerRepository(initializers.DB)
	manufacturerService := manufacturer.NewManufacturerService(manufacturerRepo)
	manufacrerHandler := manufacturer.NewHandler(manufacturerService)
	carHandler := car_params.NewHandler()
	carRoutes := router.Group("/car")
	{
		carRoutes.GET("/colors", carHandler.GetPossibleColors)
		carRoutes.GET("/transmissions", carHandler.GetPossibleTransmissions)
		carRoutes.GET("/fuel-types", carHandler.GetPossibleFuelTypes)
		carRoutes.GET("/drives", carHandler.GetPossibleDrives)
		carRoutes.GET("/manufactures", manufacrerHandler.GetAllManufactures)
		carRoutes.GET("/models/id/:id", modelHandler.GetModelsByManufacturerID)
		carRoutes.GET("/models/name/:name", modelHandler.GetModelsByManufacturerName)
	}
}

func registerSaleOfferRoutes(router *gin.Engine) {
	verifier, _ := initializeVerifier()
	saleOfferRepo := sale_offer.NewSaleOfferRepository(initializers.DB)
	manufacturerRepo := manufacturer.NewManufacturerRepository(initializers.DB)
	manufacturerService := manufacturer.NewManufacturerService(manufacturerRepo)
	saleOfferService := sale_offer.NewSaleOfferService(saleOfferRepo, manufacturerService)
	saleOfferHandler := sale_offer.NewHandler(saleOfferService)
	saleOfferRoutes := router.Group("/sale-offer")
	{
		saleOfferRoutes.POST("/", middleware.Authenticate(verifier), saleOfferHandler.CreateSaleOffer)
		saleOfferRoutes.POST("/filtered", saleOfferHandler.GetFilteredSaleOffers)
		saleOfferRoutes.GET("/offer-types", saleOfferHandler.GetOfferTypes)
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
