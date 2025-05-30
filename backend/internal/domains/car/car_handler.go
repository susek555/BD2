package car

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/enums"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service CarServiceInterface
}

func NewHandler(s CarServiceInterface) *Handler {
	return &Handler{service: s}
}

// GetManufacturersModelsMap godoc
//
//	@Summary		Get manufacturers and models map
//	@Description	Get manufacturers and models map. Each manufacturer has a list of models (the indices are corresponding).
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ManufacturerModelMap	"map of manufacturers and models"
//	@Failure		404	{object}	custom_errors.HTTPError	"manufacturer of model not found"
//	@Router			/car/manufacturer-model-map [get]
func (h *Handler) GetManufacturersModelsMap(c *gin.Context) {
	users, err := h.service.GetManufacturersModelsMap()
	if err != nil {
		c.JSON(http.StatusNotFound, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, users)
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
	c.JSON(http.StatusOK, gin.H{"colors": enums.Colors})
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
	c.JSON(http.StatusOK, gin.H{"drives": enums.Drives})
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
	c.JSON(http.StatusOK, gin.H{"fuel-types": enums.Types})
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
	c.JSON(http.StatusOK, gin.H{"transmissions": enums.Transmissions})
}
