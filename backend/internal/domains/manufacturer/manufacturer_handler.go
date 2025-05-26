package manufacturer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
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
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		RetrieveManufacturerDTO		"List of manufacturers"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/car/manufacturers [get]
func (h *Handler) GetAllManufactures(c *gin.Context) {
	manufacturers, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, manufacturers)
}
