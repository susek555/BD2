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

// UploadImages godoc
//
//	@Summary		Upload images for sale offer
//	@Description	Uploads images for a sale offer. You can upload multiple images at once, but 10 is the limit. Only offers with photos can be published later on (sale-offer/publish).
//	@Tags			image
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int										true	"Sale offer ID"
//	@Param			images	formData	file									true	"Images to upload"
//	@Success		200		{object}	sale_offer.RetrieveDetailedSaleOfferDTO	"Updated sale offer with images"
//	@Failure		400		{object}	custom_errors.HTTPError					"Invalid request"
//	@Failure		401		{object}	custom_errors.HTTPError					"Unauthorized - user must be logged in to upload images for sale offers"
//	@Failure		403		{object}	custom_errors.HTTPError					"Forbidden - user can only upload images for his own offers"
//	@Failure		404		{object}	custom_errors.HTTPError					"Sale offer not found"
//	@Failure		500		{object}	custom_errors.HTTPError					"Internal server error"
//	@Router			/image/{id} [patch]
//	@Security		Bearer
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
	offerDTO, err := h.saleOfferService.GetByID(uint(offerID), &id)
	if err != nil {
		custom_errors.HandleError(c, err, sale_offer.ErrorMap)
		return
	}
	c.JSON(http.StatusOK, offerDTO)
}

// DeleteImage godoc
//
//	@Summary		Delete image
//	@Description	Deletes an image by its URL. The user must be the owner of the image. Removes image from database and cloud storage - both operations must be successful to proceed.
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			url	query	string	true	"Image URL"
//	@Success		204	"No Content - image successfully deleted"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid request"
//	@Failure		401	{object}	custom_errors.HTTPError	"Unauthorized - user must be logged in to delete images"
//	@Failure		403	{object}	custom_errors.HTTPError	"Forbidden - usercan only images that refer to his own offers"
//	@Failure		404	{object}	custom_errors.HTTPError	"Image not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/image [delete]
//	@Security		Bearer
func (h *Handler) DeleteImage(c *gin.Context) {
	id, _ := c.Get("userID")
	userID := id.(uint)
	url := c.Query("url")
	if err := h.imageService.DeleteByURL(url, userID); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
}

// DeleteImages godoc
//
//	@Summary		Delete all images for a sale offer
//	@Description	Deletes all images for a sale offer. The user must be the owner of the given offer. Removes image from database and cloud storage - both operations must be successful to proceed.
//	@Tags			image
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Sale offer ID"
//	@Success		204	"No Content - images successfully deleted"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid request"
//	@Failure		401	{object}	custom_errors.HTTPError	"Unauthorized - user must be logged in to delete images "
//	@Failure		403	{object}	custom_errors.HTTPError	"Forbidden - user can only delete images for his own offers"
//	@Failure		404	{object}	custom_errors.HTTPError	"Sale offer not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/image/offer/{id} [delete]
//	@Security		Bearer
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
