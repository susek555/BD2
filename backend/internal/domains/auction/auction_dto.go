package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"time"
)

type CreateAuctionDTO struct {
	sale_offer.CreateSaleOfferDTO
	DateEnd     time.Time `json:"date_end"`
	BuyNowPrice uint      `json:"buy_now_price"`
}
