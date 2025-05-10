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

// GetAllManufactures godoc
//
//	@Summary		Get all manufacturers
//	@Description	Returns a list of all manufacturers stored in the database.
//	@Tags			cars
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		Manufacturer			"List of manufacturers"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/cars/manufacturers [get]
func (h *Handler) GetAllManufactures(c *gin.Context) {
	manufacturers, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, manufacturers)
}
