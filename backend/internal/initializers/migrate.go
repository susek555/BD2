package initializers

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
)

func MigrateModels() {
	DB.AutoMigrate(&user.User{})
	DB.AutoMigrate(&refresh_token.RefreshToken{})
}
