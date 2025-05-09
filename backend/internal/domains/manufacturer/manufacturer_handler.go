package manufacturer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ManufacturerServiceInterface
}

func NewHandler(s ManufacturerServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetAllManufactures(c *gin.Context) {
	manufacturers, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, manufacturers)
}
