//go:generate swag init --parseDependency --parseInternal
package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auctionws"
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
}

// @title			Car‑Dealer API
// @version		1.0
// @description	Car-Dealer API
// @contact.name	BD2
// @license.name	MIT
// @host			localhost:8080
// @schemes		http
func main() {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	hub := auctionws.NewHub()
	go hub.Run()
	hub.StartRedisFanIn(ctx, redisClient)

	router := gin.Default()
	router.Use(cors.New(middleware.CorsConfig))
	routes.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}

	// Start the WebSocket server
	verifier := jwt.NewJWTVerifier(os.Getenv("JWT_SECRET"))
	wsHandler := gin.WrapH(auctionws.ServeWS(hub))
	router.GET("/ws", middleware.Authenticate(verifier), wsHandler)

}
