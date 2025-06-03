package sale_offer

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/ws"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"github.com/susek555/BD2/car-dealer-api/pkg/pagination"
)

type Handler struct {
	service             SaleOfferServiceInterface
	hub                 ws.HubInterface
	notificationService notification.NotificationServiceInterface
}

func NewHandler(s SaleOfferServiceInterface, hub ws.HubInterface, notificationService notification.NotificationServiceInterface) *Handler {
	return &Handler{
		service:             s,
		hub:                 hub,
		notificationService: notificationService,
	}
}

// CreateSaleOffer godoc
//
//	@Summary		Create a new sale offer
//	@Description	Creates a new sale offer in the database. To create a sale offer, the user must be logged in. There are several constraints on the offer fields, such as:
//	@Description	- Color must be one of the predefined colors (endpoint: /car/colors)
//	@Description	- Fuel type must be one of the predefined fuel types (endpoint: /car/fuel_types)
//	@Description	- Transmission must be one of the predefined transmission types (endpoint: /car/transmissions)
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
//	@Param			offer	body		CreateSaleOfferDTO				true	"Sale offer form"
//	@Success		201		{object}	RetrieveDetailedSaleOfferDTO	"Created - returns the created sale offer"
//	@Failure		400		{object}	custom_errors.HTTPError			"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError			"Unauthorized - user not logged in"
//	@Failure		500		{object}	custom_errors.HTTPError			"Internal server error"
//	@Router			/sale-offer [post]
//	@Security		Bearer
func (h *Handler) CreateSaleOffer(c *gin.Context) {
	userID, _ := c.Get("userID")
	var offerDTO CreateSaleOfferDTO
	if err := c.ShouldBindJSON(&offerDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	offerDTO.UserID = userID.(uint)
	retrieveDTO, err := h.service.Create(&offerDTO)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusCreated, retrieveDTO)
	offerIDstr := strconv.FormatUint(uint64(retrieveDTO.ID), 10)
	userIDstr := strconv.FormatUint(uint64(userID.(uint)), 10)
	h.hub.SubscribeUser(userIDstr, offerIDstr)
}

// UpdateSaleOffer godoc
//
//	@Summary		Update a sale offer
//	@Description	Updates an existing sale offer in the database. To update a sale offer, the user must be logged in and must be the owner of the offer. Constraints are the same as when creating a sale offer.
//	@Tags			sale-offer
//	@Accept			json
//	@Produce		json
//	@Param			offer	body		UpdateSaleOfferDTO				true	"Sale offer form"
//	@Success		200		{object}	RetrieveDetailedSaleOfferDTO	"Updated - returns the updated sale offer"
//	@Failure		400		{object}	custom_errors.HTTPError			"Invalid input data"
//	@Failure		401		{object}	custom_errors.HTTPError			"Unauthorized - user must be logged in to update his offer"
//	@Failure		403		{object}	custom_errors.HTTPError			"Forbidden - user can only update his own offer"
//	@Failure		404		{object}	custom_errors.HTTPError			"Sale offer not found"
//	@Failure		500		{object}	custom_errors.HTTPError			"Internal server error"
//	@Router			/sale-offer [put]
//	@Security		Bearer
func (h *Handler) UpdateSaleOffer(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := userID.(uint)
	var offerDTO UpdateSaleOfferDTO
	if err := c.ShouldBindJSON(&offerDTO); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}

	retrieveDTO, err := h.service.Update(&offerDTO, id)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, retrieveDTO)
}

func (h *Handler) PublishSaleOffer(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	retrieveDTO, err := h.service.Publish(uint(id), userID.(uint))
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, retrieveDTO)
}

// GetFilteredSaleOffers godoc
//
//	@Summary		Get filtered sale offers
//	@Description	Returns a list of sale offers in paginated form. If the user is logged in, the results contain he offers created by the user. The results are filtered based on request's body. There are several constraints on the filter fields, such as:
//	@Description	- Auction type must be one of the predefined offer types (endpoint: /sale-offer/offer-types)
//	@Description	- Order key must be one of the predefined order keys (endpoint: /sale-offer/order-keys)
//	@Description	- List of manufacturers must contain only predefined manufacturers (endpoint: /car/manufacturers)
//	@Description	- List of colors must contain only predefined colors (endpoint: /car/colors)
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
	filterRequest := NewOfferFilterRequest()
	if err := c.ShouldBindJSON(filterRequest); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	filterRequest.Filter.UserID = getOptionalUserID(c)
	saleOffers, err := h.service.GetFiltered(&filterRequest.Filter, &filterRequest.PagRequest)
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
func (h *Handler) GetSaleOfferByIDDetailed(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	offerDTO, err := h.service.GetByIDDetailed(uint(id), getOptionalUserID(c))
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
//	@Failure		401		{object}	custom_errors.HTTPError			"Unauthorized - user must be logged in to retrieve his offers"
//	@Failure		500		{object}	custom_errors.HTTPError			"Internal server error"
//	@Router			/sale-offer/my-offers [post]
//	@Security		Bearer
func (h *Handler) GetMySaleOffers(c *gin.Context) {
	userID, _ := c.Get("userID")
	var pagRequest pagination.PaginationRequest
	if err := c.ShouldBindJSON(&pagRequest); err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	saleOffers, err := h.service.GetByUserID(userID.(uint), &pagRequest)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.JSON(http.StatusOK, saleOffers)
}

// GetSaleOfferTypes godoc
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

// @Summary		Buy a sale offer
// @Description	Allows a user to buy an item from a sale offer
// @Tags			SaleOffers
// @Accept			json
// @Produce		json
// @Param			id	path	uint	true	"Sale Offer ID"
// @Security		ApiKeyAuth
// @Success		200	"Successfully purchased offer"
// @Failure		403	"Forbidden - user cannot buy his own offer"
// @Failure		404	"Not Found - sale offer not found"
// @Failure		500	"Internal Server Error"
// @Failure		401	"Unauthorized - user must be logged in to buy an offer"
// @Router			/sale-offer/buy/{id} [delete]
// @Security		BearerAuth
func (h *Handler) Buy(c *gin.Context) {
	offerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	userID, err := auth.GetUserId(c)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	offer, err := h.service.Buy(uint(offerID), userID)
	if err != nil {
		custom_errors.HandleError(c, err, ErrorMap)
		return
	}
	c.Status(http.StatusOK)
	notification := &models.Notification{
		OfferID: uint(offerID),
	}
	err = h.notificationService.CreateBuyNotication(notification, strconv.FormatUint(uint64(userID), 10), offer)
	if err != nil {
		log.Printf("Error creating buy notification for offer ID %d: %v", offerID, err)
		return
	}
	h.hub.SaveNotificationForClients(strconv.FormatUint(offerID, 10), userID, notification)
	go h.hub.SendFourLatestNotificationsToClient(strconv.FormatUint(offerID, 10), strconv.FormatUint(uint64(userID), 10))
}

func getOptionalUserID(c *gin.Context) *uint {
	var id *uint
	userID, ok := c.Get("userID")
	if !ok {
		id = nil
	} else {
		uid := userID.(uint)
		id = &uid
	}
	return id
}
