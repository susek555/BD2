package user_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"gorm.io/gorm"
)

func createUser() *models.User {
	name := "john person"
	surname := "doe person"
	return &models.User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
		Selector: "P",
		Person:   &models.Person{Name: name, Surname: surname},
	}
}

func TestGetAll_EmptyDatabase(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetAll").Return([]models.User{}, nil)
	uService := user.NewUserService(uRepo)
	users, err := uService.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestGetAll_RecordsFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetAll").Return([]models.User{*createUser()}, nil)
	uService := user.NewUserService(uRepo)
	users, err := uService.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestGetById_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(nil, gorm.ErrRecordNotFound)
	uService := user.NewUserService(uRepo)
	user, err := uService.GetById(1)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Empty(t, user)
}

func TestGetById_UserFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(createUser(), nil)
	uService := user.NewUserService(uRepo)
	user, err := uService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
}

func TestGetByEmail_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetByEmail", "john@example.com").Return(models.User{}, gorm.ErrRecordNotFound)
	uService := user.NewUserService(uRepo)
	user, err := uService.GetByEmail("john@example.com")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Empty(t, user)
}

func TestGetByEmail_UserFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetByEmail", "john@example.com").Return(*createUser(), nil)
	uService := user.NewUserService(uRepo)
	user, err := uService.GetByEmail("john@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestUpdate_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(nil, gorm.ErrRecordNotFound)
	uService := user.NewUserService(uRepo)
	err := uService.Update(&user.UpdateUserDTO{ID: 1})
	assert.Equal(t, []string{gorm.ErrRecordNotFound.Error()}, err["id"])
}

func TestUpdate_UserFoundNoChange(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	existingUser := createUser()
	uRepo.On("GetById", uint(1)).Return(existingUser, nil)
	uRepo.On("GetByUsername", "john").Return(*existingUser, nil)
	uRepo.On("GetByEmail", "john@example.com").Return(*existingUser, nil)
	uRepo.On("Update", existingUser).Return(nil)
	uService := user.NewUserService(uRepo)
	err := uService.Update(&user.UpdateUserDTO{ID: 1})
	assert.Empty(t, err)
}

func TestUpdate_UserFoundChange(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(createUser(), nil)
	uService := user.NewUserService(uRepo)
	uRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
	uRepo.On("GetByUsername", "new_username").Return(models.User{}, gorm.ErrRecordNotFound)
	uRepo.On("GetByEmail", "john_updated@example.com").Return(models.User{}, gorm.ErrRecordNotFound)
	newUsername := "new_username"
	newEmail := "john_updated@example.com"
	dto := user.UpdateUserDTO{ID: 1, Username: &newUsername, Email: &newEmail}
	err := uService.Update(&dto)
	assert.Empty(t, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("Delete", uint(1)).Return(gorm.ErrRecordNotFound)
	uService := user.NewUserService(uRepo)
	err := uService.Delete(1)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDelete_UserFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("Delete", uint(1)).Return(nil)
	uService := user.NewUserService(uRepo)
	err := uService.Delete(1)
	assert.NoError(t, err)
}
