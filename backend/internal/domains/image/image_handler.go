package image

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	imageService     ImageServiceInterface
	saleOfferService sale_offer.SaleOfferServiceInterface
}

func NewHandler(imgSvc ImageServiceInterface, offerSvc sale_offer.SaleOfferServiceInterface) *Handler {
	return &Handler{imageService: imgSvc, saleOfferService: offerSvc}
}

func (h *Handler) UploadImages(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := userID.(uint)
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	files := form.File["images"]
	if err = h.imageService.Store(uint(offerID), files, id); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	offer, err := h.saleOfferService.Update(&sale_offer.UpdateSaleOfferDTO{ID: uint(offerID), Status: &models.READY}, id)
	if err != nil {
		custom_errors.HandleError(c, err, sale_offer.ErrorMap)
		return
	}
	c.JSON(http.StatusOK, offer)
}

func (h *Handler) DeleteImage(c *gin.Context) {
	id, _ := c.Get("userID")
	userID := id.(uint)
	url := c.Param("url")
	if err := h.imageService.DeleteByURL(url, userID); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteImages(c *gin.Context) {
	id, _ := c.Get("userID")
	userID := id.(uint)
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	if err := h.imageService.DeleteByOfferID(uint(offerID), userID); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
}
