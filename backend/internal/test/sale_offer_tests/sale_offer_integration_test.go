package sale_offer_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/custom_errors"
	"gorm.io/gorm"
)

// ------------------------------
// Get filtered sale offers tests
// For more get filtred tests see: sale_offer_filter_test.go
// ------------------------------

func TestGetFiltered_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{}
	server, _, _ := newTestServer(seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
}

func TestGetFiltered_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0]))
}

func TestGetFiltered_OneAuction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createAuctionSaleOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0]))
}

func TestGetFiltered_AuctionsAndOffersCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, _ := newTestServer(seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[i], got.Offers[i]))
	}
}

func TestGetFiltered_AuthorizedOtherUserOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, _ := newTestServer(seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	user := USERS[1]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, &token)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[i], got.Offers[i]))
	}
}

func TestGetFiltered_AuthorizedMyOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, _ := newTestServer(seedOffers)
	filter := sale_offer.NewOfferFilter()
	filter.Pagination = *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(filter)
	assert.NoError(t, err)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/filtered", body, &token)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got.Offers)
}

// ----------------------
// Get offfer by id tests
// ----------------------

func TestGetSaleOfferById_EmptyDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusNotFound, recievedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestGetSaleOfferById_NonExstientID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/2", nil, nil)
	assert.Equal(t, http.StatusNotFound, recievedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got.Description, gorm.ErrRecordNotFound.Error())
}

func TestGetSaleOfferById_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/-1", nil, nil)
	assert.Equal(t, http.StatusInternalServerError, recievedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
}

func TestGetSaleOfferById_StringID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/abc", nil, nil)
	assert.Equal(t, http.StatusInternalServerError, recievedStatus)
	var got custom_errors.HTTPError
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.Description)
}

func TestGetById_RegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveSaleOfferDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, seedOffers[0].ID, got.ID)
}

func TestGetById_AuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createAuctionSaleOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/id/1", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveSaleOfferDTO
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, seedOffers[0].ID, got.ID)
}

// ----------------------
// Get my offers tests
// ----------------------

func TestGetMyOffers_NotAuthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	_, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/my-offers", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, recievedStatus)
}

func TestGetMyOffers_NoOffers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{}
	server, _, _ := newTestServer(seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Empty(t, got.Offers)
}

func TestGetMyOffers_OneRegularOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0]))
}

func TestGetMyOffers_OneAuctionOffer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createAuctionSaleOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0]))
}

func TestGetMyOffers_AuctionsAndOffersCombined(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{
		*createAuctionSaleOffer(1),
		*createOffer(2),
		*createAuctionSaleOffer(3),
		*createOffer(4),
	}
	server, _, _ := newTestServer(seedOffers)
	user := USERS[0]
	token, _ := u.GetValidToken(user.ID, user.Email)
	pagRequest := *u.GetDefaultPaginationRequest()
	body, err := json.Marshal(pagRequest)
	assert.NoError(t, err)
	response, recievedStatus := u.PerformRequest(server, http.MethodPost, "/sale-offer/my-offers", body, &token)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got sale_offer.RetrieveOffersWithPagination
	err = json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, len(seedOffers), len(got.Offers))
	for i := range len(seedOffers) {
		assert.True(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[i], got.Offers[i]))
	}
}

func TestGetOfferTypes_Successful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/offer-types", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got map[string][]sale_offer.OfferType
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got["offer_types"], sale_offer.OfferTypes)
}

func TestGetOrderKeys_Successful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedOffers := []sale_offer.SaleOffer{*createOffer(1)}
	server, _, _ := newTestServer(seedOffers)
	response, recievedStatus := u.PerformRequest(server, http.MethodGet, "/sale-offer/order-keys", nil, nil)
	assert.Equal(t, http.StatusOK, recievedStatus)
	var got map[string][]string
	err := json.Unmarshal(response, &got)
	assert.NoError(t, err)
	assert.Equal(t, got["order_keys"], sale_offer.GetKeysFromMap(sale_offer.OrderKeysMap))
}
