package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	SecWebSocketHeader  = "Sec-WebSocket-Protocol"
)

type ctxKey string

const userIDKey ctxKey = "userID"
const ctxTokenKey = "wsToken"

func Authenticate(verify *jwt.JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawHeader := c.GetHeader(AuthorizationHeader)
		if rawHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		if !strings.HasPrefix(rawHeader, BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(rawHeader, BearerPrefix)

		userID, err := verify.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}

		c.Set(string(userIDKey), uint(userID))
		c.Next()
	}
}

func AuthenticateWebSocket(verify *jwt.JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(SecWebSocketHeader)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		userID, err := verify.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			return
		}

		c.Set(string(userIDKey), uint(userID))
		c.Set(ctxTokenKey, token)
		c.Next()
	}
}

func OptionalAuthenticate(verify *jwt.JWTVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawHeader := c.GetHeader(AuthorizationHeader)
		if rawHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(rawHeader, BearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(rawHeader, BearerPrefix)

		userID, err := verify.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		c.Set(string(userIDKey), uint(userID))
		c.Next()
	}
}
