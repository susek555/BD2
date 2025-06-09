package bid

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/notification"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
)

type SaleOfferAdapter struct {
	Svc sale_offer.SaleOfferServiceInterface
}

func (a SaleOfferAdapter) GetDetailedByID(id uint, userID *uint) (notification.SaleOfferInterface, error) {
	return a.Svc.GetDetailedByID(id, userID)
}
