package auction

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service AuctionServiceInterface
	sched   scheduler.SchedulerInterface
}

func NewHandler(service AuctionServiceInterface, sched scheduler.SchedulerInterface) *Handler {
	return &Handler{
		service: service,
		sched:   sched,
	}
}

// @Summary Create Auction
// @Description Creates a new auction with the provided details
// @Tags auctions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body CreateAuctionDTO true "Auction details"
// @Success 201 {object} Auction "Created auction"
// @Failure 400 {object} custom_errors.HTTPError "Bad request"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized"
// @Router /auctions [post]
// @Security BearerAuth
func (h *Handler) CreateAuction(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	var in CreateAuctionDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	in.UserID = (uint)(userId)
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
		"15:04 02/01/2006",
		in.DateEnd,
		loc,
	)
	if err != nil {
		// TODO: Do sth
	}
	dateEnd := dateEndLocal.UTC()
	h.sched.AddAuction(auctionID, dateEnd)
	log.Printf("scheduler: added %s ends %s", auctionID, dateEnd)
	c.JSON(http.StatusCreated, dto)
}

func (h *Handler) GetAllAuctions(c *gin.Context) {
	auctions, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, auctions)
}

func (h *Handler) GetAuctionById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	auction, err := h.service.GetById(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, auction)
}

func (h *Handler) DeleteAuctionById(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	err = h.service.Delete(uint(id), uint(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateAuction(c *gin.Context) {
	userId, err := auth.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, custom_errors.NewHTTPError(err.Error()))
		return
	}
	var auctionInput UpdateAuctionDTO
	if err := c.ShouldBindJSON(&auctionInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	auctionInput.UserID = (uint)(userId)
	dto, err := h.service.Update(&auctionInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto)
}
