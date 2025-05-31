package auction

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/scheduler"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

type Handler struct {
	service AuctionServiceInterface
	sched   scheduler.SchedulerInterface
	hub     ws.HubInterface
}

func NewHandler(service AuctionServiceInterface, sched scheduler.SchedulerInterface, hub ws.HubInterface) *Handler {
	return &Handler{
		service: service,
		sched:   sched,
		hub:     hub,
	}
}

// @Summary		Create Auction
// @Description	Creates a new auction with the provided details
// @Tags			auction
// @Accept			json
// @Produce		json
// @Param			body	body		CreateAuctionDTO		true	"Auction details"
// @Success		201		{object}	RetrieveAuctionDTO		"Created auction"
// @Failure		400		{object}	custom_errors.HTTPError	"Bad request"
// @Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
// @Router			/auction [post]
// @Security		BearerAuth
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
		formats.DateTimeLayout,
		in.DateEnd,
		loc,
	)
	if err != nil {
		// TODO: Do sth
	}
	dateEnd := dateEndLocal.UTC()
	userIdStr := strconv.FormatUint(uint64(userId), 10)
	h.hub.SubscribeUser(userIdStr, auctionID)
	h.sched.AddAuction(auctionID, dateEnd)
	log.Printf("scheduler: added %s ends %s", auctionID, dateEnd)
	c.JSON(http.StatusCreated, dto)
}

// GetAllAuctions godoc
//
//	@Summary		Get all auctions
//	@Description	Retrieves all available auctions from the system
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		RetrieveAuctionDTO
//	@Failure		400	{object}	custom_errors.HTTPError
//	@Router			/auction [get]
func (h *Handler) GetAllAuctions(c *gin.Context) {
	auctions, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, auctions)
}

// GetAuctionById godoc
//
//	@Summary		Get auction by ID
//	@Description	Retrieves a specific auction by its ID
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Auction ID"
//	@Success		200	{object}	RetrieveAuctionDTO
//	@Failure		400	{object}	custom_errors.HTTPError
//	@Router			/auction/{id} [get]
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

// DeleteAuctionById godoc
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

// UpdateAuction godoc
//
//	@Summary		Update auction
//	@Description	Updates an existing auction with the provided details
//	@Tags			auction
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateAuctionDTO		true	"Auction details"
//	@Success		200		{object}	RetrieveAuctionDTO		"Updated auction"
//	@Failure		400		{object}	custom_errors.HTTPError	"Bad request"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized"
//	@Router			/auction [put]
//	@Security		BearerAuth
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
	dto, err := h.service.Update(&auctionInput, uint(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto)
}

// BuyNow godoc
// @Summary Buy an auction at its buy now price
// @Description Allows a user to instantly purchase an auction at its buy now price if available
// @Tags auctions
// @Accept json
// @Produce json
// @Param id path int true "Auction ID"
// @Success 200 "Successfully purchased the auction"
// @Failure 400 {object} custom_errors.HTTPError "Invalid auction ID or buy now operation failed"
// @Failure 401 {object} custom_errors.HTTPError "Unauthorized - user not logged in"
// @Router /auctions/buy-now/{id} [delete]
// @Security BearerAuth
func (h *Handler) BuyNow(c *gin.Context) {
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
	err = h.service.BuyNow(uint(id), uint(userId))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}
