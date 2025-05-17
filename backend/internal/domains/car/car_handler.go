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

// GetManufacturersModelsMap godoc
//
//	@Summary		Get manufacturers and models map
//	@Description	Get manufacturers and models map. Each manufacturer has a list of models (the indices are coresponding).
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
