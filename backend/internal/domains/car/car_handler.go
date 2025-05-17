package car

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service CarServiceInterface
}

func NewCarHandler(s CarServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetManufacturersModelsMap(c *gin.Context) {
	users, err := h.service.GetManufacturersModelsMap()
	if err != nil {
		c.JSON(http.StatusNotFound, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, users)
}
