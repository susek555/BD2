package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/utils"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

type ctxKey string

const userIDKey ctxKey = "userID"

func Authenticate(verify *utils.JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawHeader := c.GetHeader(authorizationHeader)
		if rawHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		if !strings.HasPrefix(rawHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(rawHeader, bearerPrefix)

		userID, err := verify.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}

		c.Set(string(userIDKey), userID)
		c.Next()
	}
}
