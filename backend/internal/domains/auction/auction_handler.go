package auction

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

type RetrieveDetailedSaleOfferDTO = sale_offer.RetrieveDetailedSaleOfferDTO

type Handler struct {
	service             AuctionServiceInterface
	sched               scheduler.SchedulerInterface
	hub                 ws.HubInterface
	notificationService notification.NotificationServiceInterface
}

func NewHandler(service AuctionServiceInterface, sched scheduler.SchedulerInterface, hub ws.HubInterface, notificationService notification.NotificationServiceInterface) *Handler {
	return &Handler{
		service:             service,
		sched:               sched,
		hub:                 hub,
		notificationService: notificationService,
	}
}

// CreateAuction godoc
//
//	@Summary		Create Auction
//	@Description	Creates a new auction with the provided details
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateAuctionDTO				true	"Auction details"
//	@Success		201		{object}	RetrieveDetailedSaleOfferDTO	"Created auction"
//	@Failure		400		{object}	custom_errors.HTTPError			"Bad request"
//	@Failure		401		{object}	custom_errors.HTTPError			"Unauthorized"
//	@Router			/auction [post]
//	@Security		BearerAuth
func (h *Handler) CreateAuction(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	var in CreateAuctionDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	in.UserID = userID
	dto, err := h.service.Create(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	auctionID := strconv.FormatUint(uint64(dto.ID), 10)
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		return
	}
	dateEndLocal, err := time.ParseInLocation(
		formats.DateTimeLayout,
		in.DateEnd,
		loc,
	)
	if err != nil {
		log.Printf("Error parsing dateEnd: %v", err)
		return
	}
	dateEnd := dateEndLocal.UTC()
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	h.hub.SubscribeUser(userIDStr, auctionID)
	h.sched.AddAuction(auctionID, dateEnd)
	log.Printf("scheduler: added %s ends %s", auctionID, dateEnd)
	c.JSON(http.StatusCreated, dto)
}

// UpdateAuction godoc
//
//	@Summary		Update auction
//	@Description	Updates an existing auction with the provided details
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateAuctionDTO				true	"Auction details"
//	@Success		200		{object}	RetrieveDetailedSaleOfferDTO	"Updated auction"
//	@Failure		400		{object}	custom_errors.HTTPError			"Bad request"
//	@Failure		401		{object}	custom_errors.HTTPError			"Unauthorized"
//	@Router			/auction [put]
//	@Security		BearerAuth
func (h *Handler) UpdateAuction(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	var auctionInput UpdateAuctionDTO
	if err := c.ShouldBindJSON(&auctionInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	dto, err := h.service.Update(&auctionInput, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto)
	auctionID := strconv.FormatUint(uint64(dto.ID), 10)
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		return
	}
	dateEndLocal, err := time.ParseInLocation(
		formats.DateTimeLayout,
		*auctionInput.DateEnd,
		loc,
	)
	if err != nil {
		log.Printf("Error parsing dateEnd: %v", err)
		return
	}
	dateEnd := dateEndLocal.UTC()
	h.sched.ModifyAuction(auctionID, dateEnd)
}

// DeleteAuctionByID godoc
//
//	@Summary		Delete auction by ID
//	@Description	Deletes a specific auction by its ID
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Auction ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	custom_errors.HTTPError
//	@Failure		401	{object}	custom_errors.HTTPError
//	@Router			/auction/{id} [delete]
//	@Security		BearerAuth
func (h *Handler) DeleteAuctionByID(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	err = h.service.Delete(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

// BuyNow godoc
//
//	@Summary		Buy an auction at its buy now price
//	@Description	Allows a user to instantly purchase an auction at its buy now price if available
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Auction ID"
//	@Success		200	"Successfully purchased the auction"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid auction ID or buy now operation failed"
//	@Failure		401	{object}	custom_errors.HTTPError	"Unauthorized - user not logged in"
//	@Router			/auction/buy-now/{id} [post]
//	@Security		BearerAuth
func (h *Handler) BuyNow(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	offer, err := h.service.BuyNow(uint(id), uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
	notification := &models.Notification{
		OfferID: uint(id),
	}
	err = h.notificationService.CreateBuyNowNotification(notification, strconv.FormatUint(uint64(userID), 10), offer)
	if err != nil {
		log.Printf("Error creating buy now notification for auction ID %d: %v", id, err)
		return
	}
	h.hub.SaveNotificationForClients(strconv.FormatUint(id, 10), userID, notification)
	h.hub.SendFourLatestNotificationsToClients(strconv.FormatUint(id, 10), strconv.FormatUint(uint64(userID), 10))
	h.sched.ForceCloseAuction(strconv.FormatUint(id, 10), userID, offer.GetPrice())
	h.hub.RemoveRoom(strconv.FormatUint(id, 10))
}
