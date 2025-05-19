package user_tests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
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
	db.Exec("PRAGMA foreign_keys = ON")
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

func getRepositories(db *gorm.DB) (user.UserRepositoryInterface, generic.CRUDRepository[user.Company], generic.CRUDRepository[user.Person]) {
	userRepo := user.NewUserRepository(db)
	comapnyRepo := generic.GetGormRepository[user.Company](db)
	personRepo := generic.GetGormRepository[user.Person](db)
	return userRepo, comapnyRepo, personRepo
}

func getRepositoryWithUsers(db *gorm.DB, users []user.User) user.UserRepositoryInterface {
	repo := user.NewUserRepository(db)
	for _, user := range users {
		repo.Create(&user)
	}
	return repo
}

func newTestServer(db *gorm.DB, seedUsers []user.User) (*gin.Engine, error) {
	verifier := jwt.NewJWTVerifier(u.JWTSECRET)
	userRepo := getRepositoryWithUsers(db, seedUsers)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)
	r := gin.Default()
	userRoutes := r.Group("/users")
	{
		userRoutes.PUT("/", middleware.Authenticate(verifier), userHandler.UpdateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/id/:id", userHandler.GetUserById)
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)
		userRoutes.DELETE("/id/:id", middleware.Authenticate(verifier), userHandler.DeleteUser)
	}
	return r, nil
}

// ------------
// Basic models
// ------------

func createPerson(id uint) *user.User {
	user := user.User{
		ID:       id,
		Username: fmt.Sprintf("john%d", id),
		Email:    fmt.Sprintf("john%d@gmail.com", id),
		Password: "hashed_password",
		Selector: "P",
		Person:   &user.Person{Name: "john person", Surname: "doe person"},
	}
	return &user
}

func createCompany(id uint) *user.User {
	user := user.User{
		ID:       id,
		Username: fmt.Sprintf("john%d", id),
		Email:    fmt.Sprintf("john%d@gmail.com", id),
		Password: "hashed_password",
		Selector: "C",
		Company:  &user.Company{Name: "john company", NIP: fmt.Sprintf("1234567890-%d", id)},
	}
	return &user
}

func withCompanyField(opt u.Option[user.Company]) u.Option[user.User] {
	return func(userObj *user.User) {
		if userObj.Company == nil {
			userObj.Company = &user.Company{}
		}
		opt(userObj.Company)
	}
}

func doUserAndRetrieveUserDTOsMatch(user user.User, dto user.RetrieveUserDTO) bool {
	if user.ID != dto.ID || user.Username != dto.Username || user.Email != dto.Email {
		return false
	}
	if (user.Company == nil) != (dto.CompanyName == nil) {
		return false
	}
	if user.Company != nil && user.Company.Name != *dto.CompanyName {
		return false
	}

	if (user.Person == nil) != (dto.PersonName == nil) {
		return false
	}
	if user.Person != nil && user.Person.Name != *dto.PersonName {
		return false
	}
	return true
}
