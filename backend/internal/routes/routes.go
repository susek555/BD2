package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/initializers"
)

func RegisterRoutes(router *gin.Engine) {
	registerUserRoutes(router)
}

func registerUserRoutes(router *gin.Engine) {
	userRepo := user.GetUserRepository(initializers.DB)
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
