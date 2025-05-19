package sale_offer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type Handler struct {
	service SaleOfferServiceInterface
}

func NewSaleOfferHandler(s SaleOfferServiceInterface) *Handler {
	return &Handler{service: s}
}

// CreateSaleOffer godoc
//
//	@Summary		Create a new sale offer
//	@Description	Creates a new sale offer in the database. To create a sale offer, the user must be logged in. There are several constraints on the offer fields, such as:
//	@Description	- Color must be one of the predefined colors (endpoint: /car/colors)
//	@Description	- Fuel type must be one of the predefined fuel types (endpoint: /car/fuel_types)
//	@Description	- Transmission must be one of the predefined transmission types (endpotin: /car/transmissions)
//	@Description	- Drive must be one of the predefined drive types (endpoint: /car/drives)
//	@Description	- Model must be one of the predefined models (endpoint: /car/models or /car/models/:id)
//	@Description	- Number of doors must be between 1 and 6
//	@Description	- Number of seats must be between 2 and 100
//	@Description	- Engine power must be less than or equal to 9999 (in horsepower)
//	@Description	- Engine capacity must be less than or equal to 9000 (in cm3)
//	@Description	- Number of gears must be between 1 and 10
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Param			offer	body		CreateSaleOfferDTO		true	"Sale offer form"
//	@Success		201		{object}	CreateSaleOfferDTO		"Created - returns the created sale offer"
//	@Failure		400		{object}	custom_errors.HTTPError	"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError	"Unauthorized - user not logged in"
//	@Failure		500		{object}	custom_errors.HTTPError	"Internal server error"
//	@Router			/sale-offer [post]
//	@Security		Bearer
func (h *Handler) CreateSaleOffer(c *gin.Context) {
	var offerDTO CreateSaleOfferDTO
	if err := c.ShouldBindJSON(&offerDTO); err != nil {
		c.JSON(http.StatusBadRequest, custom_errors.NewHTTPError(err.Error()))
		return
	}
	userID, _ := c.Get("userID")
	offerDTO.UserID = userID.(uint)
	dto, err := h.service.Create(offerDTO)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusCreated, dto)
}

// GetFilteredSaleOffers godoc
//
//	@Summary		Get filtered sale offers
//	@Description	Returns a list of sale offers in paginated form. If the user is logged in, the results are not containing the offers created by the user. The results are filtered based on request's body. There are several constraints on the filter fields, such as:
//	@Description	- Offer type must be one of the predefined offer types (endpoint: /sale-offer/offer-types)
//	@Description	- Order key must be one of the predefined order keys (endpoint: /sale-offer/order-keys)
//	@Description	- List of manufacturers must contain only predefined manufacturers (endpoint: /car/manufacturers)
//	@Description	- List of colors must containonly predefined colors (endpoint: /car/colors)
//	@Description	- List of drives must contain only predefined drives (endpoint: /car/drives)
//	@Description	- List of fuel types must contain only predefined fuel types (endpoint: /car/fuel_types)
//	@Description	- List of transmissions must contain only predefined transmission types (endpoint: /car/transmissions)
//	@Description	- Whenever you use a range, the min value must be less than or equal to the max value, you can provide only one of them, and the other will be ignored.
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Param			filter	body		OfferFilter						true	"Sale offer filter"
//	@Success		200		{object}	RetrieveOffersWithPagination	"List of sale offers"
//	@Failure		400		{object}	custom_errors.HTTPError			"Invalid input data"
//	@Failure		500		{object}	custom_errors.HTTPError			"Internal server error"
//	@Router			/sale-offer/filtered [post]
func (h *Handler) GetFilteredSaleOffers(c *gin.Context) {
	filter := NewOfferFilter()
	if err := c.ShouldBindJSON(filter); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, ok := c.Get("userID")
	if !ok {
		filter.UserID = nil
	} else {
		uid := userID.(uint)
		filter.UserID = &uid
	}
	saleOffers, err := h.service.GetFiltered(filter)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, saleOffers)
}

// GetSaleOfferByID godoc
//
//	@Summary		Get sale offer by ID
//	@Description	Returns a sale offer by its ID. Can be used to retrieve detailed information about sale offer.
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint							true	"Sale offer ID"
//	@Success		200	{object}	RetrieveDetailedSaleOfferDTO	"Sale offer details"
//	@Failure		400	{object}	custom_errors.HTTPError			"Invalid input data"
//	@Failure		404	{object}	custom_errors.HTTPError			"Sale offer not found"
//	@Failure		500	{object}	custom_errors.HTTPError			"Internal server error"
//	@Router			/sale-offer/id/{id} [get]
func (h *Handler) GetSaleOfferByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	offerDTO, err := h.service.GetByID(uint(id))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, offerDTO)
}

// GetMySaleOffers godoc
//
//	@Summary		Get my sale offers
//	@Description	Returns a list of all sale offers created by the logged-in user.
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Param			filter	body		pagination.PaginationRequest	true	"Pagination request"
//	@Success		200		{object}	RetrieveOffersWithPagination	"List of sale offers"
//	@Failure		401		{object}	custom_errors.HTTPError			"Unauthorized - user not logged in"
//	@Failure		500		{object}	custom_errors.HTTPError			"Internal server error"
//	@Router			/sale-offer/my-offers [post]
//	@Security		Bearer
func (h *Handler) GetMySaleOffers(c *gin.Context) {
	var pagRequest pagination.PaginationRequest
	if err := c.ShouldBindJSON(&pagRequest); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, _ := c.Get("userID")
	saleOffers, err := h.service.GetByUserID(userID.(uint), &pagRequest)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, saleOffers)
}

// GetOfferTypes godoc
//
//	@Summary		Get offer types
//	@Description	Returns a list of all possible offer types that are accepted when using filtering. If you choose both the auctions and regular offers will be found.
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string	"List of offer types"
//	@Router			/sale-offer/offer-types [get]
func (h *Handler) GetSaleOfferTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"offer_types": OfferTypes})
}

// GetOrderKeys godoc
//
//	@Summary		Get order keys
//	@Description	Returns a list of all possible order keys that are accepted when using filtering.
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string][]string	"List of order keys"
//	@Router			/sale-offer/order-keys [get]
func (h *Handler) GetOrderKeys(c *gin.Context) {
	keys := make([]string, 0, len(OrderKeysMap))
	for k := range OrderKeysMap {
		keys = append(keys, k)
	}
	c.JSON(http.StatusOK, gin.H{"order_keys": keys})
}
