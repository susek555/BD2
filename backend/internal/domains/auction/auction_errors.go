package auction

import "errors"

var ErrBuyNowPriceLessThan1 = errors.New("buy now price must be greater than 0")
var ErrBuyNowPriceLessThanOfferPrice = errors.New("buy now price must be greater than offer price")
var ErrModificationForbidden = errors.New("provided offer does not belong to logged in user")
var ErrAuctionOwnedByUser = errors.New("you cannot buy your own auction")
var ErrAuctionAlreadySold = errors.New("auction is already sold - cannot buy it again")
var ErrNewPriceLessThanOfferPrice = errors.New("new price must be greater than offer price")
