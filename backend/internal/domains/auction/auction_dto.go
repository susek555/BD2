package auction

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
)

type CreateAuctionDTO struct {
	sale_offer.CreateSaleOfferDTO
	DateEnd     string `json:"date_end"`
	BuyNowPrice uint   `json:"buy_now_price"`
}

type RetrieveAuctionDTO struct {
	*sale_offer.RetrieveSaleOfferDTO
	DateEnd     string `json:"date_end"`
	BuyNowPrice uint   `json:"buy_now_price"`
}

type UpdateAuctionDTO struct {
	id uint `json:"id"`
	CreateAuctionDTO
}
