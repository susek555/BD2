package auction

import "errors"

var (
	ErrBuyNowPriceLessThan1          = errors.New("buy now price must be greater than 0")
	ErrBuyNowPriceLessThanOfferPrice = errors.New("buy now price must be greater than offer price")
	ErrModificationForbidden         = errors.New("provided offer does not belong to logged in user")
	ErrAuctionOwnedByUser            = errors.New("you cannot buy your own auction")
	ErrAuctionNotOwned               = errors.New("provided offer does not belong to logged in user")
	ErrAuctionAlreadySold            = errors.New("auction is already sold - cannot buy it again")
	ErrNewPriceLessThanOfferPrice    = errors.New("new price must be greater than offer price")
	ErrAuctionNotPublished           = errors.New("auction is not published")
)
