package sale_offer

type OfferType string

const (
	REGULAR_OFFER OfferType = "Regular offer"
	AUCTION       OfferType = "Auction"
	BOTH          OfferType = "Both"
)

var OfferTypes = []OfferType{REGULAR_OFFER, AUCTION, BOTH}
