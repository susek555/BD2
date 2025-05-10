package auction

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"net/http"
)

type Handler struct {
	service AuctionServiceInterface
}

func NewHandler(service AuctionServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateAuction(c *gin.Context) {
	var in CreateAuctionDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	err := h.service.Create(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, in)
}

func (h *Handler) GetAllAuctions(c *gin.Context) {
	auctions, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, auctions)
}
