package main

import (
	"log"
	"os"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/routes"

	"github.com/gin-gonic/gin"

	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.MigrateModels()
}

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}
	jwtKey := []byte(secret)
	verifier := jwt.NewJWTVerifier(secret)

	db := initializers.DB

	authSvc := auth.NewService(db, jwtKey)
	authH := auth.NewHandler(authSvc)

	router := gin.Default()

	authGroup := router.Group("/auth")
	authGroup.POST("/register", authH.Register)
	authGroup.POST("/login", authH.Login)
	authGroup.POST("/refresh", authH.Refresh)

	api := router.Group("/")
	api.Use(middleware.Authenticate(verifier))
	routes.RegisterRoutes(router)
	{
		api.POST("/logout", authH.Logout)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
