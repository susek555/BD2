package routes

import (
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
}

func registerUserRoutes(router *gin.Engine) {
	userRepo := user.NewUserRepository(initializers.DB)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.PUT("/", userHandler.UpdateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/id/:id", userHandler.GetUserById)
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}
}

func registerAuthRoutes(router *gin.Engine) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	jwtKey := []byte(secret)
	verifier := jwt.NewJWTVerifier(secret)
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
