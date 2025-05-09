package car_params

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetPossibleColors(c *gin.Context) {
	c.JSON(http.StatusOK, Colors)
}

func (h *Handler) GetPossibleDrives(c *gin.Context) {
	c.JSON(http.StatusOK, Drives)
}

func (h *Handler) GetPossibleFuelTypes(c *gin.Context) {
	c.JSON(http.StatusOK, Types)
}

func (h *Handler) GetPossibleTransmissions(c *gin.Context) {
	c.JSON(http.StatusOK, Transmissions)
}
