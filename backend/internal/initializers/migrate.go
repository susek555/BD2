package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/models"
)

func MigrateModels() {
	DB.AutoMigrate(
		&models.User{},
		&models.Person{},
		&models.Company{},
		&models.RefreshToken{},
		&models.Review{},
		&models.Manufacturer{},
		&models.Model{},
		&models.SaleOffer{},
		&models.Car{},
		&models.Auction{},
		&models.Bid{},
		&models.LikedOffer{},
		&models.Notification{},
		&models.ClientNotification{},
		&models.Image{},
	)
}
