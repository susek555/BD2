package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/bid"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/manufacturer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/model"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/review"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/sale_offer"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

func MigrateModels() {
	DB.AutoMigrate(
		&user.User{},
		&user.Person{},
		&user.Company{},
		&refresh_token.RefreshToken{},
		&review.Review{},
		&manufacturer.Manufacturer{},
		&model.Model{},
		&sale_offer.SaleOffer{},
		&sale_offer.Car{},
		&sale_offer.Auction{},
		&bid.Bid{},
	)
}
