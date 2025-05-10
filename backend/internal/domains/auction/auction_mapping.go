package auction

import (
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"time"
)

func (dto *CreateAuctionDTO) MapToAuction() (*sale_offer.Auction, error) {
	var auction sale_offer.Auction
	offer, err := dto.MapToSaleOffer()
	if err != nil {
		return nil, err
	}
	endDate, err := parseDate(dto.DateEnd)
	if err != nil {
		return nil, err
	}
	err = validateDateEnd(endDate)
	if err != nil {
		return nil, err
	}
	auction.Offer = offer
	auction.DateEnd = endDate
	auction.BuyNowPrice = dto.BuyNowPrice
	return &auction, nil
}

func parseDate(date string) (time.Time, error) {
	layout := "11:20 02/01/2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func validateDateEnd(end time.Time) error {
	if !end.After(time.Now()) {
		return errors.New("date end must be in the future")
	}
	return nil
}
