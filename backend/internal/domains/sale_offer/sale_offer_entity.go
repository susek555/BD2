package sale_offer

// Can be either view or model
type SaleOfferEntityInterface interface {
	GetID() uint
	BelongsToUser(userID uint) bool
	IsAuctionOffer() bool
}
