package bid

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auctionws"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service     BidServiceInterface
	redisClient *redis.Client
	hub         *auctionws.Hub
}

func NewHandler(service BidServiceInterface, redisClient *redis.Client, hub *auctionws.Hub) *Handler {
	return &Handler{
		service:     service,
		redisClient: redisClient,
		hub:         hub,
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
	dto, err := h.service.Create(&in, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto)
	auctionIDStr := strconv.FormatUint(uint64(dto.AuctionID), 10)
	userIDStr := strconv.FormatUint(uint64(userId), 10)
	amountInt64 := int64(dto.Amount)
	env := auctionws.NewBidEnvelope(auctionIDStr, amountInt64, userIDStr)
	data, err := json.Marshal(env)
	if err != nil {
		log.Println("Error marshalling envelope:", err)
		return
	}
	h.hub.BroadcastLocal(auctionIDStr, data, userIDStr)

	auctionws.PublishAuctionEvent(c, h.redisClient, auctionIDStr, env)

	h.hub.SubscribeUser(userIDStr, auctionIDStr)
}

func (h *Handler) GetAllBids(c *gin.Context) {
	bids, err := h.service.GetAll()
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
	bid, err := h.service.GetById(uint(bidID))
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
	bids, err := h.service.GetByBidderId(uint(bidderId))
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
	bids, err := h.service.GetByAuctionId(uint(auctionId))
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
	bid, err := h.service.GetHighestBid(uint(auctionId))
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
	bid, err := h.service.GetHighestBidByUserId(uint(auctionId), uint(bidderId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
	}
	c.JSON(http.StatusOK, bid)
}
