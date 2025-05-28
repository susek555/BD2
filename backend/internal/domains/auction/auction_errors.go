package auction

import "errors"

var ErrBuyNowPriceLessThan1 = errors.New("buy now price must be greater than 0")
var ErrBuyNowPriceLessThanOfferPrice = errors.New("buy now price must be greater than offer price")