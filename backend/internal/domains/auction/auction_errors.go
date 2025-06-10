package auction

import "errors"

var (
	ErrBuyNowPriceLessThan1          = errors.New("buy now price must be greater than 0")
	ErrBuyNowPriceLessThanOfferPrice = errors.New("buy now price must be greater than offer price")
	ErrNewPriceLessThanOfferPrice    = errors.New("new price must be greater than offer price")
	ErrBuyNowNotAvailable            = errors.New("buy now option is not available for this auction")
)
