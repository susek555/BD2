package sale_offer_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
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

func TestGetFiltered_OneOffer(t *testing.T) {
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
	assert.Equal(t, doSaleOfferAndRetrieveSaleOfferDTOsMatch(seedOffers[0], got.Offers[0]), true)
}
