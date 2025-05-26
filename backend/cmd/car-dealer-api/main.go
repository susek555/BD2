//go:generate swag init --parseDependency --parseInternal
package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"

	"github.com/susek555/BD2/car-dealer-api/internal/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/susek555/BD2/car-dealer-api/docs"
	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.MigrateModels()
	initializers.ConnectToRedis(context.Background())
	initializers.ConnectToCloudinary()
	initializers.InitializeVerifier()
	initializers.InitializeHub()
	initializers.InitializeRepos()
	initializers.InitializeServices()
	initializers.InitializeScheduler()
	initializers.InitializeHandlers()
}

// @title			Car‑Dealer API
// @version		1.0
// @description	Car-Dealer API
// @contact.name	BD2
// @license.name	MIT
// @host			localhost:8080
// @schemes		http
func main() {
	router := gin.Default()
	router.Use(cors.New(middleware.CorsConfig))
	routes.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
