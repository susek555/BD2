package bid

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service BidServiceInterface
}

func NewHandler(service BidServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateBid(c *gin.Context) {
	var in Bid
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	in.CreatedAt = time.Now()
	err := h.service.Create(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
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
	bidderId, err := strconv.ParseUint(c.Param("bidder"), 10, 32)
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
	auctionId, err := strconv.ParseUint(c.Param("auctionId"), 10, 32)
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
	auctionId, err := strconv.ParseUint(c.Param("auctionId"), 10, 32)
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
