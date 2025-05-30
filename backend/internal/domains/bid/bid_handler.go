package bid

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	bidService          BidServiceInterface
	redisClient         *redis.Client
	hub                 *ws.Hub
	notificationService notification.NotificationServiceInterface
}

func NewHandler(service BidServiceInterface, redisClient *redis.Client, hub *ws.Hub, notificationService notification.NotificationServiceInterface) *Handler {
	return &Handler{
		bidService:          service,
		redisClient:         redisClient,
		hub:                 hub,
		notificationService: notificationService,
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
	userId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	dto, err := h.bidService.Create(&in, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	retrieveDTO := ProcessingToRetrieve(dto)
	c.JSON(http.StatusCreated, retrieveDTO)
	auctionIDStr := strconv.FormatUint(uint64(dto.AuctionID), 10)
	userIDStr := strconv.FormatUint(uint64(userId), 10)
	amountInt64 := int64(dto.Amount)
	notification := &models.Notification{
		OfferID: dto.AuctionID,
	}

	err = h.notificationService.CreateOutbidNotification(notification, amountInt64, dto.Auction)
	if err != nil {
		log.Println("Error creating notification:", err)
		return
	}
	h.hub.SaveNotificationForClients(auctionIDStr, userId, notification)
	// TODO: think about the best way to do this
	// auctionws.PublishAuctionEvent(c, h.redisClient, auctionIDStr, env)
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
	bid, err := h.bidService.GetById(uint(bidID))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}

// GetBidsByBidderId godoc
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
func (h *Handler) GetBidsByBidderId(c *gin.Context) {
	bidderId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bids, err := h.bidService.GetByBidderId(uint(bidderId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bids)
}

// GetBidsByAuctionId godoc
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
func (h *Handler) GetBidsByAuctionId(c *gin.Context) {
	auctionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bids, err := h.bidService.GetByAuctionId(uint(auctionId))
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
	auctionId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bid, err := h.bidService.GetHighestBid(uint(auctionId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}

// GetHighestBidByUserId godoc
//
//	@Summary		Get the highest bid by a user for a specific auction
//	@Description	Retrieves the highest bid placed by a specific user on a specific auction
//	@Tags			bid
//	@Accept			json
//	@Produce		json
//	@Param			auctionId	path		uint					true	"Auction ID"
//	@Param			bidderId	path		uint					true	"Bidder ID"
//	@Success		200			{object}	RetrieveBidDTO			"Highest bid details"
//	@Failure		400			{object}	custom_errors.HTTPError	"Invalid auction ID, bidder ID, or retrieval error"
//	@Router			/bid/highest/auction/{auctionId}/bidder/{bidderId} [get]
func (h *Handler) GetHighestBidByUserId(c *gin.Context) {
	auctionId, err := strconv.ParseUint(c.Param("auctionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bidderId, err := strconv.ParseUint(c.Param("bidderId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	bid, err := h.bidService.GetHighestBidByUserId(uint(auctionId), uint(bidderId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}
