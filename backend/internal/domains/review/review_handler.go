package review

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"net/http"
	"strconv"
)

type Handler struct {
	service ReviewServiceInterface
}

func NewHandler(service ReviewServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllReviews(c *gin.Context) {
	reviews, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func (h *Handler) GetReviewById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	review, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, review)
}

func (h *Handler) CreateReview(c *gin.Context) {
	var review Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := h.service.Create(&review); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusCreated, review)
}
