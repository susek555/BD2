package notification

import "github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"

type Notification struct {
	ID                    uint `json:"id" gorm:"primaryKey"`
	OfferID               uint `json:"offer_id"`
	*sale_offer.SaleOffer `json:"sale_offer" gorm:"foreignKey:OfferID;references:ID"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	Date                  string `json:"date"`
}
