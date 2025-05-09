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
	service ModelServiceInterace
}

func NewHandler(s ModelServiceInterace) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetModelsByManufacturerID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	models, err := h.service.GetByManufactuerID(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, models)
}

func (h *Handler) GetModelsByManufacturerName(c *gin.Context) {
	name := c.Param("email")
	models, err := h.service.GetByManufacuterName(name)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, models)
}
