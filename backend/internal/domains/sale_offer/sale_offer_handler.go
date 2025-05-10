package sale_offer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service SaleOfferServiceInterface
}

func NewHandler(s SaleOfferServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateSaleOffer(c *gin.Context) {
	var offerDTO CreateSaleOfferDTO
	if err := c.ShouldBindJSON(&offerDTO); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	if err := h.service.Create(offerDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusCreated, offerDTO)
}
