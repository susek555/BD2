package liked_offer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service LikeOfferServiceInterface
}

func NewLikedOfferHandler(s LikeOfferServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) LikeNewOffer(c *gin.Context) {
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userID, _ := c.Get("userID")
	likedOffer := LikedOffer{OfferID: uint(offerID), UserID: userID.(uint)}
	if err := h.service.Create(&likedOffer); err != nil {
		c.JSON(http.StatusInternalServerError, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, likedOffer)
}

func (h *Handler) DislikeOffer(c *gin.Context) {
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userID, _ := c.Get("userID")
	if err := h.service.Delete(uint(offerID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
