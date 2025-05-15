package auction

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
)

type Handler struct {
	service AuctionServiceInterface
}

func NewHandler(service AuctionServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateAuction(c *gin.Context) {
	var in CreateAuctionDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	auction, err := in.MapToAuction()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	err = h.service.Create(auction)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	dto := MapToDTO(auction)
	c.JSON(http.StatusOK, dto)
}

func (h *Handler) GetAllAuctions(c *gin.Context) {
	auctions, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	var auctionsDTO []RetrieveAuctionDTO
	for _, auction := range auctions {
		dto := MapToDTO(&auction)
		auctionsDTO = append(auctionsDTO, *dto)
	}
	c.JSON(http.StatusOK, auctionsDTO)
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
	dto := MapToDTO(auction)
	c.JSON(http.StatusOK, dto)
}

func (h *Handler) DeleteAuctionById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	err = h.service.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateAuction(c *gin.Context) {
	var auctionInput UpdateAuctionDTO
	if err := c.ShouldBindJSON(&auctionInput); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	auction, err := auctionInput.MapToAuction()
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	err = h.service.Update(auction)
	if err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	dto := MapToDTO(auction)
	c.JSON(http.StatusOK, dto)
}
