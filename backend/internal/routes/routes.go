package routes

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"log"
	"os"

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
}
