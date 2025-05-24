package image

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, _ := c.Get("userID")
	offer, err := h.saleOfferService.GetByID(uint(offerID), nil)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	if offer.UserID != userID {
		custom_errors.HandleError(c, ErrInvalidOfferID, ErrorMap)
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	files := form.File["images"]
	urls, err := h.imageService.StoreImages(uint(offerID), files)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, urls)
}
