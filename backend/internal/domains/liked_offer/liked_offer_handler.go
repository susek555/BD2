package liked_offer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service LikedOfferServiceInterface
	hub     ws.HubInterface
}

func NewHandler(service LikedOfferServiceInterface, hub ws.HubInterface) *Handler {
	return &Handler{
		service: service,
		hub:     hub,
	}
}

// LikeNewOffer godoc
//
//	@Summary		Like new offer
//	@Description	Like new offer by giving it's id. You have to be logged in to perform this operation.
//	@Tags			favourite
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint					true	"Sale offer ID"
//	@Success		200	{object}	nil						"Liked offer"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401	{object}	custom_errors.HTTPError	"Unauthorized - user not logged in"
//	@Failure		404	{object}	custom_errors.HTTPError	"Sale offer not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/favourite/like/{id} [post]
//	@Security		Bearer
func (h *Handler) LikeOffer(c *gin.Context) {
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, _ := c.Get("userID")
	if err := h.service.LikeOffer(uint(offerID), userID.(uint)); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusOK)
	offerIDStr := strconv.FormatUint(offerID, 10)
	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
	h.hub.SubscribeUser(userIDStr, offerIDStr)

}

// DislikeOffer godoc
//
//	@Summary		Dislike offer
//	@Description	Dislike offer by giving it's id. You have to be logged in to perform this operation.
//	@Tags			favourite
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint					true	"Sale offer ID"
//	@Success		204	{object}	nil						"No content"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401	{object}	custom_errors.HTTPError	"Unauthorized - user not logged in"
//	@Failure		404	{object}	custom_errors.HTTPError	"Sale offer not found"
//	@Failure		500	{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/favourite/dislike/{id} [delete]
//	@Security		Bearer
func (h *Handler) DislikeOffer(c *gin.Context) {
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, _ := c.Get("userID")
	if err := h.service.DislikeOffer(uint(offerID), userID.(uint)); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusNoContent)
	offerIDStr := strconv.FormatUint(offerID, 10)
	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
	h.hub.UnsubscribeUser(userIDStr, offerIDStr)
}
