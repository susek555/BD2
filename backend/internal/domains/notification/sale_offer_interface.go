package notification

import "github.com/susek555/BD2/car-dealer-api/internal/enums"

type SaleOfferInterface interface {
	GetBrand() string
	GetModel() string
	GetPrice() uint
	HasBuyNowPrice() bool
	GetStatus() enums.Status
	BelongsToUser(userID uint) bool
	GetID() uint
}
