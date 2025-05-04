//go:generate swag init --parseDependency --parseInternal
package main

import (
	"log"
	"os"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/routes"

	"github.com/gin-gonic/gin"

	_ "github.com/susek555/BD2/car-dealer-api/docs"
	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.MigrateModels()
}

// @title			Car‑Dealer API
// @version		1.0
// @description	Car-Dealer API
// @contact.name	BD2
// @license.name	MIT
// @host			localhost:8080
// @schemes		http
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
