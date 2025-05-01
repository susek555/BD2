package main

import (
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"github.com/susek555/BD2/car-dealer-api/pkg/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
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

	verifier := utils.NewJWTVerifier(secret)
	router := setupRouter(verifier)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func setupRouter(verifier *utils.JWTVerifier) *gin.Engine {
	r := gin.Default()

	// PUBLIC
	//auth := r.Group("/auth")
	//{
	//	auth.POST("/register", register)
	//	auth.POST("/login", login)
	//}

	// PRIVATE
	api := r.Group("/")
	api.Use(middleware.Authenticate(verifier))
	{
		api.GET("/users", getUsers)
	}

	return r
}

// This code will be in UserController
func getUsers(c *gin.Context) {
	userRepository := user.GetUserRepository(initializers.DB)
	users, err := userRepository.GetAllUsers()
	if err != nil {
		log.Fatal("Something went wrong")
	}
	c.IndentedJSON(http.StatusOK, users)
}
