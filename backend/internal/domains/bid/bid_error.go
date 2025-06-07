package bid

import "errors"

var ErrAuctionNotPublished = errors.New("auction is not published")
var ErrBidderIsAuctionOwner = errors.New("bidder cannot bid on their own auction")
