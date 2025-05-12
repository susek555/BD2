package auth_tests

import (
	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(users []user.User) (user.UserRepositoryInterface, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&user.User{},
		&refresh_token.RefreshToken{},
	)
	if err != nil {
		return nil, err
	}
	repo := user.NewUserRepository(db)
	for _, user := range users {
		repo.Create(&user)
	}
	return repo, nil
}

func newTestServer(seedUsers []user.User) (*gin.Engine, error) {
	repo, err := setupDB(seedUsers)
	if err != nil {
		return nil, err
	}

	svc := &auth.AuthService{Repo: repo}
	h := &auth.Handler{Service: svc}

	r := gin.Default()
	r.POST("/register", h.Register)
	return r, nil
}
