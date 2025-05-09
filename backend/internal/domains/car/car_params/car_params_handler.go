package car_params

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetPossibleColors godoc
//
//	@Summary		Get all possible colors
//	@Description	Returns a list of all possible colors that are accepted when creating a new offer. If color of your car is not in the list, chose 'other'.
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	string	"List of colors"
//	@Router			/cars/colors [get]
func (h *Handler) GetPossibleColors(c *gin.Context) {
	c.JSON(http.StatusOK, Colors)
}

// GetPossibleDrives godoc
//
//	@Summary		Get all possible drives
//	@Description	Returns a list of all possible drives that are accepted when creating a new offer.
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	string	"List of drives"
//	@Router			/cars/drives [get]
func (h *Handler) GetPossibleDrives(c *gin.Context) {
	c.JSON(http.StatusOK, Drives)
}

// GetPossibleFuelTypes godoc
//
//	@Summary		Get all possible fuel types
//	@Description	Returns a list of all possible fuel types that are accepted when creating a new offer.
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	string	"List of fuel types"
//	@Router			/cars/fuel-types [get]
func (h *Handler) GetPossibleFuelTypes(c *gin.Context) {
	c.JSON(http.StatusOK, Types)
}

// GetPossibleTransmissions godoc
//
//	@Summary		Get all possible transmissions
//	@Description	Returns a list of all possible transmissions that are accepted when creating a new offer.
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	string	"List of transmissions"
//	@Router			/cars/transmissions [get]
func (h *Handler) GetPossibleTransmissions(c *gin.Context) {
	c.JSON(http.StatusOK, Transmissions)
}
