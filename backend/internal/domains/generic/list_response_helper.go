package generic

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleListResponse[T any](c *gin.Context, dtos []T) {
	if len(dtos) == 0 {
		c.JSON(http.StatusOK, []T{})
		return
	}
	c.JSON(http.StatusOK, dtos)
}
