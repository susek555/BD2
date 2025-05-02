package initializers

import "github.com/susek555/BD2/car-dealer-api/internal/domains/user"

func MigrateModels() {
	DB.AutoMigrate(
		&user.User{},
		&user.Person{},
		&user.Company{},
	)
}
