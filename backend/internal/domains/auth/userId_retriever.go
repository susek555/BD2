package auth

import "github.com/gin-gonic/gin"

func GetUserId(c *gin.Context) uint {
	userId, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userId.(uint)
}
