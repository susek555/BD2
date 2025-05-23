package model

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"gorm.io/gorm"
)

var ErrorMap = map[error]int{
	strconv.ErrSyntax:      http.StatusBadRequest,
	gorm.ErrRecordNotFound: http.StatusNotFound,
}

type Handler struct {
	service ModelServiceInterface
}

func NewHandler(s ModelServiceInterface) *Handler {
	return &Handler{service: s}
}

// GetModelsByManufacturerID godoc
//
//	@Summary		Get all models by manufacturer id
//	@Description	Returns a list of all models stored in the database for a given manufacturer id.
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"Manufacturer ID"
//	@Success		200	{array}		RetrieveModelDTO		"List of models"
//	@Failure		400	{object}	custom_errors.HTTPError	"Id is not a number"
//	@Failure		404	{object}	custom_errors.HTTPError	"Models not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/car/models/id/{id} [get]
func (h *Handler) GetModelsByManufacturerID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	models, err := h.service.GetByManufacturerID(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, models)
}

// GetModelsByManufacturerName godoc
//
//	@Summary		Get all models by manufacturer name
//	@Description	Returns a list of all models stored in the database for a given manufacturer name.
//	@Tags			car
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string					true	"Manufacturer name"
//	@Success		200		{array}		RetrieveModelDTO		"List of models"
//	@Failure		404		{object}	custom_errors.HTTPError	"Models not found"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/car/models/name/{name} [get]
func (h *Handler) GetModelsByManufacturerName(c *gin.Context) {
	name := c.Param("name")
	models, err := h.service.GetByManufacturerName(name)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, models)
}
