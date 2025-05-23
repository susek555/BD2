package bid

import (
	"encoding/json"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auctionws"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	bidService          BidServiceInterface
	redisClient         *redis.Client
	hub                 *auctionws.Hub
	notificationService notification.NotificationServiceInterface
}

func NewHandler(service BidServiceInterface, redisClient *redis.Client, hub *auctionws.Hub, notificationService notification.NotificationServiceInterface) *Handler {
	return &Handler{
		bidService:          service,
		redisClient:         redisClient,
		hub:                 hub,
		notificationService: notificationService,
	}
}

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
	c.JSON(http.StatusCreated, dto)
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
	env := auctionws.NewNotificationEnvelope(notification)
	data, err := json.Marshal(env)
	if err != nil {
		log.Println("Error marshalling envelope:", err)
		return
	}
	h.hub.BroadcastLocal(auctionIDStr, data, userIDStr)
	// TODO: think about the best way to do this
	// auctionws.PublishAuctionEvent(c, h.redisClient, auctionIDStr, env)

	h.hub.SubscribeUser(userIDStr, auctionIDStr)
}

func (h *Handler) GetAllBids(c *gin.Context) {
	bids, err := h.bidService.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bids)
}

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
