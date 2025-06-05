package auction

import (
	"errors"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/formats"
)

func (dto *CreateAuctionDTO) MapToAuction() (*models.Auction, error) {
	var auction models.Auction
	endDate, err := parseDate(dto.DateEnd)
	if err != nil {
		return nil, err
	}
	err = validateDateEnd(endDate)
	if err != nil {
		return nil, err
	}
	auction.DateEnd = endDate
	auction.BuyNowPrice = dto.BuyNowPrice
	auction.InitialPrice = dto.Price
	return &auction, nil
}

func (dto *UpdateAuctionDTO) UpdatedAuctionFromDTO(auction *models.Auction) (*models.Auction, error) {
	if dto.DateEnd != nil {
		endDate, err := parseDate(*dto.DateEnd)
		if err != nil {
			return nil, err
		}
		err = validateDateEnd(endDate)
		if err != nil {
			return nil, err
		}
		auction.DateEnd = endDate
	}
	if dto.BuyNowPrice != nil {
		if *dto.BuyNowPrice < 1 {
			return nil, ErrBuyNowPriceLessThan1
		}
		if *dto.BuyNowPrice < auction.Offer.Price {
			return nil, ErrBuyNowPriceLessThanOfferPrice
		}
		auction.BuyNowPrice = *dto.BuyNowPrice
	}
	return auction, nil
}

func parseDate(date string) (time.Time, error) {
	loc, err := time.LoadLocation(formats.DefaultTimezone)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(formats.DateTimeLayout, date, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func validateDateEnd(end time.Time) error {
	if !end.After(time.Now()) {
		return errors.New("date end must be in the future")
	}
	return nil
}
