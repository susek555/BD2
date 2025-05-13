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

func (h *Handler) GetFilteredSaleOffers(c *gin.Context) {
	var filter OfferFilter
	if err := c.ShouldBindJSON(&filter); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	saleOffers, err := h.service.GetFiltered(&filter)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, *saleOffers)
}

func (h *Handler) GetOfferTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"offer_types": OfferTypes})
}

func (h *Handler) GetOrderKeys(c *gin.Context) {
	keys := make([]string, 0, len(OrderMap))
	for k := range OrderMap {
		keys = append(keys, k)
	}
	c.JSON(http.StatusOK, gin.H{"order_keys": keys})
}
