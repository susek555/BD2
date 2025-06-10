package sale_offer

import "github.com/susek555/BD2/car-dealer-api/internal/enums"

type SaleOfferEntityInterface interface {
	GetID() uint
	BelongsToUser(userID uint) bool
	IsAuctionOffer() bool
	GetStatus() enums.Status
}
