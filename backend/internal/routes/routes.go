package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
)

func RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is healthy",
		})
	})

	registerAuthRoutes(router)
	registerUserRoutes(router)
	registerReviewRoutes(router)
	registerCarRoutes(router)
	registerSaleOfferRoutes(router)
	registerAuctionRoutes(router)
	registerBidRoutes(router)
	registerWebsocket(router)
	registerImageRoutes(router)
	registerFavouriteRoutes(router)
}

func registerWebsocket(router *gin.Engine) {
	// Start the WebSocket server
	wsHandler := ws.ServeWS(initializers.Hub.(*ws.Hub))
	router.GET("/ws", middleware.AuthenticateWebSocket(initializers.Verifier), wsHandler)
}

func registerUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.PUT("/", middleware.Authenticate(initializers.Verifier), initializers.UserHandler.UpdateUser)
		userRoutes.GET("/", initializers.UserHandler.GetAllUsers)
		userRoutes.GET("/id/:id", initializers.UserHandler.GetUserByID)
		userRoutes.GET("/email/:email", initializers.UserHandler.GetUserByEmail)
		userRoutes.DELETE("/:id", middleware.Authenticate(initializers.Verifier), initializers.UserHandler.DeleteUser)
	}
}

func registerAuthRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", initializers.AuthHandler.Register)
		authRoutes.POST("/login", initializers.AuthHandler.Login)
		authRoutes.POST("/refresh", initializers.AuthHandler.Refresh)
		authRoutes.PUT("/change-password", middleware.Authenticate(initializers.Verifier), initializers.AuthHandler.ChangePassword)
	}
	router.POST("/logout", middleware.Authenticate(initializers.Verifier), initializers.AuthHandler.Logout)
}

func registerReviewRoutes(router *gin.Engine) {
	reviewRoutes := router.Group("/review")
	reviewRoutes.GET("/", initializers.ReviewHandler.GetAllReviews)
	reviewRoutes.GET("/:id", initializers.ReviewHandler.GetReviewByID)
	reviewRoutes.POST("/", middleware.Authenticate(initializers.Verifier), initializers.ReviewHandler.CreateReview)
	reviewRoutes.PUT("/", middleware.Authenticate(initializers.Verifier), initializers.ReviewHandler.UpdateReview)
	reviewRoutes.DELETE("/:id", middleware.Authenticate(initializers.Verifier), initializers.ReviewHandler.DeleteReview)
	reviewRoutes.POST("/reviewer/:id", initializers.ReviewHandler.GetReviewsByReviewerID)
	reviewRoutes.POST("/reviewee/:id", initializers.ReviewHandler.GetReviewsByRevieweeID)
	reviewRoutes.GET("/reviewer/reviewee/:reviewerID/:revieweeID", initializers.ReviewHandler.GetReviewsByReviewerIDAndRevieweeID)
	reviewRoutes.POST("/filtered", initializers.ReviewHandler.GetFilteredReviews)
	reviewRoutes.GET("/average-rating/:id", initializers.ReviewHandler.GetAverageRatingByRevieweeID)
	reviewRoutes.GET("/frequency/:id", initializers.ReviewHandler.GetFrequencyOfRatingByRevieweeID)
}

func registerCarRoutes(router *gin.Engine) {
	carRoutes := router.Group("/car")
	{
		carRoutes.GET("/colors", initializers.CarHandler.GetPossibleColors)
		carRoutes.GET("/transmissions", initializers.CarHandler.GetPossibleTransmissions)
		carRoutes.GET("/fuel-types", initializers.CarHandler.GetPossibleFuelTypes)
		carRoutes.GET("/drives", initializers.CarHandler.GetPossibleDrives)
		carRoutes.GET("/manufactures", initializers.ManufacturerHandler.GetAllManufactures)
		carRoutes.GET("/models/id/:id", initializers.ModelHandler.GetModelsByManufacturerID)
		carRoutes.GET("/models/name/:name", initializers.ModelHandler.GetModelsByManufacturerName)
		carRoutes.GET("/manufacturer-model-map", initializers.CarHandler.GetManufacturersModelsMap)
	}
}

func registerSaleOfferRoutes(router *gin.Engine) {
	saleOfferRoutes := router.Group("/sale-offer")
	{
		saleOfferRoutes.POST("/", middleware.Authenticate(initializers.Verifier), initializers.SaleOfferHandler.CreateSaleOffer)
		saleOfferRoutes.PUT("/", middleware.Authenticate(initializers.Verifier), initializers.SaleOfferHandler.UpdateSaleOffer)
		saleOfferRoutes.PUT("/publish/:id", middleware.Authenticate(initializers.Verifier), initializers.SaleOfferHandler.PublishSaleOffer)
		saleOfferRoutes.POST("/my-offers", middleware.Authenticate(initializers.Verifier), initializers.SaleOfferHandler.GetMySaleOffers)
		saleOfferRoutes.POST("/filtered", middleware.OptionalAuthenticate(initializers.Verifier), initializers.SaleOfferHandler.GetFilteredSaleOffers)
		saleOfferRoutes.GET("/id/:id", middleware.OptionalAuthenticate(initializers.Verifier), initializers.SaleOfferHandler.GetDetailedSaleOfferByID)
		saleOfferRoutes.GET("/offer-types", initializers.SaleOfferHandler.GetSaleOfferTypes)
		saleOfferRoutes.GET("/order-keys", initializers.SaleOfferHandler.GetOrderKeys)
		saleOfferRoutes.POST("/buy/:id", middleware.Authenticate(initializers.Verifier), initializers.SaleOfferHandler.Buy)
	}
}

func registerAuctionRoutes(router *gin.Engine) {
	auctionRoutes := router.Group("/auction")
	auctionRoutes.POST("/", middleware.Authenticate(initializers.Verifier), initializers.AuctionHandler.CreateAuction)
	auctionRoutes.PUT("/", middleware.Authenticate(initializers.Verifier), initializers.AuctionHandler.UpdateAuction)
	auctionRoutes.DELETE("/:id", middleware.Authenticate(initializers.Verifier), initializers.AuctionHandler.DeleteAuctionByID)
	auctionRoutes.POST("/buy-now/:id", middleware.Authenticate(initializers.Verifier), initializers.AuctionHandler.BuyNow)
}

func registerBidRoutes(router *gin.Engine) {
	bidRoutes := router.Group("/bid")
	bidRoutes.POST("/", middleware.Authenticate(initializers.Verifier), initializers.BidHandler.CreateBid)
	bidRoutes.GET("/", initializers.BidHandler.GetAllBids)
	bidRoutes.GET("/:id", initializers.BidHandler.GetBidByID)
	bidRoutes.GET("/bidder/:id", initializers.BidHandler.GetBidsByBidderID)
	bidRoutes.GET("/auction/:id", initializers.BidHandler.GetBidsByAuctionID)
	bidRoutes.GET("/highest/:id", initializers.BidHandler.GetHighestBid)
	bidRoutes.GET("/highest/auction/:auctionID/bidder/:bidderID", initializers.BidHandler.GetHighestBidByUserID)
}

func registerFavouriteRoutes(router *gin.Engine) {
	favourtiesRoutes := router.Group("/favourite")
	{
		favourtiesRoutes.POST("/like/:id", middleware.Authenticate(initializers.Verifier), initializers.LikedOfferHandler.LikeOffer)
		favourtiesRoutes.DELETE("/dislike/:id", middleware.Authenticate(initializers.Verifier), initializers.LikedOfferHandler.DislikeOffer)
	}
}

func registerImageRoutes(router *gin.Engine) {
	imageRoutes := router.Group("/image")
	{
		imageRoutes.PATCH("/:id", middleware.Authenticate(initializers.Verifier), initializers.ImageHandler.UploadImages)
		imageRoutes.DELETE("/", middleware.Authenticate(initializers.Verifier), initializers.ImageHandler.DeleteImage)
		imageRoutes.DELETE("/offer/:id", middleware.Authenticate(initializers.Verifier), initializers.ImageHandler.DeleteImages)
	}
}
