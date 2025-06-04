package user_tests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	u "github.com/susek555/BD2/car-dealer-api/internal/test/test_utils"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// -----
// Setup
// -----

func setupDB() (*gorm.DB, error) {
	dsn := "host=localhost user=bd2_user password=bd2_password dbname=bd2_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("TRUNCATE TABLE companies, people, users RESTART IDENTITY CASCADE")
	return db, nil
}

func getRepositories(db *gorm.DB) (user.UserRepositoryInterface, generic.CRUDRepository[models.Company], generic.CRUDRepository[models.Person]) {
	userRepo := user.NewUserRepository(db)
	companyRepo := generic.GetGormRepository[models.Company](db)
	personRepo := generic.GetGormRepository[models.Person](db)
	return userRepo, companyRepo, personRepo
}

func getRepositoryWithUsers(db *gorm.DB, users []models.User) user.UserRepositoryInterface {
	repo := user.NewUserRepository(db)
	for _, user := range users {
		repo.Create(&user)
	}
	return repo
}

func newTestServer(db *gorm.DB, seedUsers []models.User) (*gin.Engine, error) {
	verifier := jwt.NewJWTVerifier(u.JWTSECRET)
	userRepo := getRepositoryWithUsers(db, seedUsers)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)
	r := gin.Default()
	userRoutes := r.Group("/users")
	{
		userRoutes.PUT("/", middleware.Authenticate(verifier), userHandler.UpdateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/id/:id", userHandler.GetUserByID)
		userRoutes.GET("/email/:email", userHandler.GetUserByEmail)
		userRoutes.DELETE("/id/:id", middleware.Authenticate(verifier), userHandler.DeleteUser)
	}
	return r, nil
}

// ------------
// Basic models
// ------------

func createPerson(id uint) *models.User {
	user := models.User{
		ID:       id,
		Username: fmt.Sprintf("john%d", id),
		Email:    fmt.Sprintf("john%d@gmail.com", id),
		Password: "hashed_password",
		Selector: "P",
		Person:   &models.Person{Name: "john person", Surname: "doe person"},
	}
	return &user
}

func createCompany(id uint) *models.User {
	user := models.User{
		ID:       id,
		Username: fmt.Sprintf("john%d", id),
		Email:    fmt.Sprintf("john%d@gmail.com", id),
		Password: "hashed_password",
		Selector: "C",
		Company:  &models.Company{Name: "john company", Nip: fmt.Sprintf("1234567890-%d", id)},
	}
	return &user
}

func withCompanyField(opt u.Option[models.Company]) u.Option[models.User] {
	return func(userObj *models.User) {
		if userObj.Company == nil {
			userObj.Company = &models.Company{}
		}
		opt(userObj.Company)
	}
}

func doUserAndRetrieveUserDTOsMatch(user models.User, dto user.RetrieveUserDTO) bool {
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
