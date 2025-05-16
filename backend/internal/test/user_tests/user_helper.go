package user_tests

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// -----
// Setup
// -----

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&user.User{},
		&user.Company{},
		&user.Person{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getSubtypesRepositories(db *gorm.DB) (generic.CRUDRepository[user.Company], generic.CRUDRepository[user.Person]) {
	comapnyRepo := generic.GetGormRepository[user.Company](db)
	personRepo := generic.GetGormRepository[user.Person](db)
	return comapnyRepo, personRepo
}

func insertUsersToDatabase(users []user.User, repo user.UserRepositoryInterface) error {
	for _, user := range users {
		err := repo.Create(&user)
		if err != nil {
			return err
		}
	}
	return nil
}

func newTestServer(seedUsers []user.User) (*gin.Engine, error) {
	db, err := setupDB()
	if err != nil {
		return nil, err
	}
	verifier := jwt.NewJWTVerifier(JWTSECRET)
	userRepo := user.NewUserRepository(db)
	if err := insertUsersToDatabase(seedUsers, userRepo); err != nil {
		return nil, err
	}
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	r := gin.Default()
	userRoutes := r.Group("/users")
	{
		userRoutes.PUT("/", middleware.Authenticate(verifier), userHandler.UpdateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/id/:id", userHandler.GetUserById)
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)
		userRoutes.DELETE("/:id", middleware.Authenticate(verifier), userHandler.DeleteUser)
	}
	return r, nil
}

// -------------------
// Authorization setup
// -------------------

const JWTSECRET = "secret"

func getValidToken(userId uint, email string) (string, error) {
	secret := []byte("secret")
	return jwt.GenerateToken(email, int64(userId), secret, time.Now().Add(1*time.Hour))
}

// ------------
// Basic models
// ------------

func createPerson(id uint) *user.User {
	user := user.User{
		ID:       id,
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
		Selector: "P",
		Person:   &user.Person{Name: "john person", Surname: "doe person"},
	}
	return &user
}

func createCompany(id uint) *user.User {
	user := user.User{
		ID:       id,
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
		Selector: "C",
		Company:  &user.Company{Name: "john company", NIP: "1234567890"},
	}
	return &user
}
