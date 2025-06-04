package auction

import (
	"errors"
	"log"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
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
	if dto.UpdateSaleOfferDTO != nil {
		var err error
		auction.Offer, err = dto.UpdateSaleOfferDTO.UpdatedOfferFromDTO(auction.Offer)
		if err != nil {
			return nil, err
		}
	}
	return auction, nil
}

func MapToDTO(auction *models.Auction) *RetrieveAuctionDTO {
	loc, err := time.LoadLocation(formats.DefaultTimezone)
	if err != nil {
		log.Println("error loading location:", err)
	}

	dateEnd := auction.DateEnd.In(loc).Format(formats.DateTimeLayout)
	offerDTO := sale_offer.MapToDTO(auction.Offer)
	return &RetrieveAuctionDTO{
		offerDTO,
		dateEnd,
		auction.BuyNowPrice,
	}
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
