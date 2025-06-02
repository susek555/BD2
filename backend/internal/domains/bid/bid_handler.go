package bid

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	bidService          BidServiceInterface
	redisClient         *redis.Client
	hub                 ws.HubInterface
	notificationService notification.NotificationServiceInterface
	sched               scheduler.SchedulerInterface
}

func NewHandler(service BidServiceInterface, redisClient *redis.Client, hub ws.HubInterface, notificationService notification.NotificationServiceInterface, sched scheduler.SchedulerInterface) *Handler {
	return &Handler{
		bidService:          service,
		redisClient:         redisClient,
		hub:                 hub,
		notificationService: notificationService,
		sched:               sched,
	}
}

// CreateBid godoc
//
//	@Summary		Create a new bid
//	@Description	Create a new bid for an auction
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateBidDTO			true	"Bid details"
//	@Success		201		{object}	RetrieveBidDTO			"Created bid"
//	@Failure		400		{object}	custom_errors.HTTPError	"Bad request"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
//	@Router			/bid [post]
//	@Security		BearerAuth
func (h *Handler) CreateBid(c *gin.Context) {
	var in CreateBidDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	dto, err := h.bidService.Create(&in, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	retrieveDTO := ProcessingToRetrieve(dto)
	c.JSON(http.StatusCreated, retrieveDTO)
	auctionIDStr := strconv.FormatUint(uint64(dto.AuctionID), 10)
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	amountInt64 := int64(dto.Amount)
	notification := &models.Notification{
		OfferID: dto.AuctionID,
	}

	err = h.notificationService.CreateOutbidNotification(notification, amountInt64, dto.Auction)
	if err != nil {
		log.Println("Error creating notification:", err)
		return
	}
	h.hub.SaveNotificationForClients(auctionIDStr, userID, notification)
	if in.Amount >= dto.Auction.BuyNowPrice {
		h.sched.CloseAuction(auctionIDStr)
	}
	go h.hub.SendFourLatestNotificationsToClient(auctionIDStr, userIDStr)

	h.hub.SubscribeUser(userIDStr, auctionIDStr)
}

// GetAllBids godoc
//
//	@Summary		Get all bids
//	@Description	Retrieve all bids
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		RetrieveBidDTO			"List of bids"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad request"
//	@Router			/bid [get]
func (h *Handler) GetAllBids(c *gin.Context) {
	bids, err := h.bidService.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bids)
}

// GetBidByID godoc
//
//	@Summary		Get bid by ID
//	@Description	Retrieve a bid by its ID
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"Bid ID"
//	@Success		200	{object}	RetrieveBidDTO			"Bid details"
//	@Failure		400	{object}	custom_errors.HTTPError	"Bad request"
//	@Router			/bid/{id} [get]
func (h *Handler) GetBidByID(c *gin.Context) {
	bidID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bid, err := h.bidService.GetByID(uint(bidID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}

// GetBidsByBidderID godoc
//
//	@Summary		Get bids by bidder ID
//	@Description	Retrieves all bids placed by a specific bidder
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint					true	"Bidder ID"
//	@Success		200	{array}		RetrieveBidDTO			"List of bids"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid bidder ID or retrieval error"
//	@Router			/bid/bidder/{id} [get]
func (h *Handler) GetBidsByBidderID(c *gin.Context) {
	bidderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bids, err := h.bidService.GetByBidderID(uint(bidderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bids)
}

// GetBidsByAuctionID godoc
//
//	@Summary		Get bids by auction ID
//	@Description	Retrieves all bids placed on a specific auction
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint					true	"Auction ID"
//	@Success		200	{array}		RetrieveBidDTO			"List of bids"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid auction ID or retrieval error"
//	@Router			/bid/auction/{id} [get]
func (h *Handler) GetBidsByAuctionID(c *gin.Context) {
	auctionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bids, err := h.bidService.GetByAuctionID(uint(auctionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bids)
}

// GetHighestBid godoc
//
//	@Summary		Get the highest bid for an auction
//	@Description	Retrieves the highest bid for a specific auction
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint					true	"Auction ID"
//	@Success		200	{object}	RetrieveBidDTO			"Highest bid details"
//	@Failure		400	{object}	custom_errors.HTTPError	"Invalid auction ID or retrieval error"
//	@Router			/bid/highest/{id} [get]
func (h *Handler) GetHighestBid(c *gin.Context) {
	auctionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bid, err := h.bidService.GetHighestBid(uint(auctionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}

// GetHighestBidByUserID godoc
//
//	@Summary		Get the highest bid by a user for a specific auction
//	@Description	Retrieves the highest bid placed by a specific user on a specific auction
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		uint					true	"Auction ID"
//	@Param			bidderID	path		uint					true	"Bidder ID"
//	@Success		200			{object}	RetrieveBidDTO			"Highest bid details"
//	@Failure		400			{object}	custom_errors.HTTPError	"Invalid auction ID, bidder ID, or retrieval error"
//	@Router			/bid/highest/auction/{auctionID}/bidder/{bidderID} [get]
func (h *Handler) GetHighestBidByUserID(c *gin.Context) {
	auctionID, err := strconv.ParseUint(c.Param("auctionID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bidderID, err := strconv.ParseUint(c.Param("bidderID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bid, err := h.bidService.GetHighestBidByUserID(uint(auctionID), uint(bidderID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}
