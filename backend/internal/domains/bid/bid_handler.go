package bid

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"net/http"
	"time"
)

type Handler struct {
	service BidServiceInterface
}

func NewHandler(service BidServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateBid(c *gin.Context) {
	var in Bid
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	in.CreatedAt = time.Now()
	err := h.service.Create(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
