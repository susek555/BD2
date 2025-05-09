//go:build unit
// +build unit

package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	u "github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"gorm.io/gorm"
)

func createUser() *u.User {
	name := "john person"
	surname := "doe person"
	return &u.User{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Password: "hashed_password",
		Selector: "P",
		Person:   &u.Person{Name: name, Surname: surname},
	}
}

func TestGetAll_EmptyDatabase(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetAll").Return([]u.User{}, nil)
	uService := u.NewService(uRepo)
	users, err := uService.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestGetAll_RecordsFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetAll").Return([]u.User{*createUser()}, nil)
	uService := u.NewService(uRepo)
	users, err := uService.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestGetById_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(u.User{}, gorm.ErrRecordNotFound)
	uService := u.NewService(uRepo)
	user, err := uService.GetById(1)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Empty(t, user)
}

func TestGetById_UserFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(*createUser(), nil)
	uService := u.NewService(uRepo)
	user, err := uService.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Username)
}

func TestGetByEmail_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetByEmail", "john@example.com").Return(u.User{}, gorm.ErrRecordNotFound)
	uService := u.NewService(uRepo)
	user, err := uService.GetByEmail("john@example.com")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Empty(t, user)
}

func TestGetByEmail_UserFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetByEmail", "john@example.com").Return(*createUser(), nil)
	uService := u.NewService(uRepo)
	user, err := uService.GetByEmail("john@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestUpdate_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(u.User{}, gorm.ErrRecordNotFound)
	uService := u.NewService(uRepo)
	err := uService.Update(u.UpdateUserDTO{ID: 1})
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUpdate_UserFoundNoChange(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(*createUser(), nil)
	uRepo.On("Update", createUser()).Return(nil)
	uService := u.NewService(uRepo)
	err := uService.Update(u.UpdateUserDTO{ID: 1})
	assert.NoError(t, err)
}

func TestUpdate_UserFoundChange(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("GetById", uint(1)).Return(*createUser(), nil)
	uService := u.NewService(uRepo)
	uRepo.On("Update", mock.AnythingOfType("*user.User")).Return(nil)
	newUsername := "new_username"
	newEmail := "john_updated@example.com"
	dto := u.UpdateUserDTO{ID: 1, Username: &newUsername, Email: &newEmail}
	err := uService.Update(dto)
	assert.NoError(t, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("Delete", uint(1)).Return(gorm.ErrRecordNotFound)
	uService := u.NewService(uRepo)
	err := uService.Delete(1)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDelete_UserFound(t *testing.T) {
	uRepo := mocks.NewUserRepositoryInterface(t)
	uRepo.On("Delete", uint(1)).Return(nil)
	uService := u.NewService(uRepo)
	err := uService.Delete(1)
	assert.NoError(t, err)
}
