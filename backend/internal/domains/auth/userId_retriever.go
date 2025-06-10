package auth

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, ErrUserIDNotFound
	}
	return userID.(uint), nil
}
