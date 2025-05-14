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
//	@Description	Returns a list of all possible colors that are accepted when creating a new offer. If the color of your car is not in the list, choose 'other'.
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string	"List of colors"
//	@Router			/car/colors [get]
func (h *Handler) GetPossibleColors(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"colors": Colors})
}

// GetPossibleDrives godoc
//
//	@Summary		Get all possible drives
//	@Description	Returns a list of all possible drives that are accepted when creating a new offer.
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string	"List of drives"
//	@Router			/car/drives [get]
func (h *Handler) GetPossibleDrives(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"drives": Drives})
}

// GetPossibleFuelTypes godoc
//
//	@Summary		Get all possible fuel types
//	@Description	Returns a list of all possible fuel types that are accepted when creating a new offer.
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string	"List of fuel types"
//	@Router			/car/fuel-types [get]
func (h *Handler) GetPossibleFuelTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"fuel-types": Types})
}

// GetPossibleTransmissions godoc
//
//	@Summary		Get all possible transmissions
//	@Description	Returns a list of all possible transmissions that are accepted when creating a new offer.
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string	"List of transmissions"
//	@Router			/car/transmissions [get]
func (h *Handler) GetPossibleTransmissions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"transmissions": Transmissions})
}
