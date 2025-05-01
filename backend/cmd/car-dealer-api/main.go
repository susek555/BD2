package main

import (
	"github.com/susek555/BD2/car-dealer-api/middleware"
	"github.com/susek555/BD2/car-dealer-api/utils"
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
	verifier := utils.NewJWTVerifier(os.Getenv("JWT_SECRET"))
	router := gin.Default()
	router.Use(middleware.Authenticate(verifier))
	router.GET("/users/all", getUsers)
	router.Run() // listen and serve on 0.0.0.0:8080
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
