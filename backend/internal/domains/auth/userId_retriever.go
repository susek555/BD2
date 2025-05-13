package auth

import "github.com/gin-gonic/gin"

func GetUserId(c *gin.Context) (uint, error) {
	userId, exists := c.Get("userID")
	if !exists {
		return 0, ErrUserIdNotFound
	}
	return userId.(uint), nil
}
