package sale_offer

type SaleOfferEntityInterface interface {
	GetID() uint
	BelongsToUser(userID uint) bool
	IsAuctionOffer() bool
}
