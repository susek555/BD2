package middleware

import (
	"github.com/gin-contrib/cors"
	"time"
)

var CorsConfig = cors.Config{
	AllowOrigins:     []string{"http://localhost:3000"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
}
