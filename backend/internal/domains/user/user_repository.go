package user

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user User) error
	GetAllUsers() ([]User, error)
	GetUserById(id uint) (User, error)
	UpdateUser(user User) error
	DeleteUser(id uint)
	GetUserByEmail(email string) (User, error)
}

type UserRepository struct {
	gormRepository *generic.GormRepository[User]
}

func GetUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{gormRepository: generic.GetGormRepository[User](db)}
}

func (userRepository *UserRepository) CreateUser(user User) error {
	return userRepository.gormRepository.Create(user)
}

func (userRepository *UserRepository) GetAllUsers() ([]User, error) {
	return userRepository.gormRepository.GetAll()
}

func (userRepository *UserRepository) GetUserById(id uint) (User, error) {
	return userRepository.gormRepository.GetById(id)
}

func (userRepository *UserRepository) UpdateUser(user User) error {
	return userRepository.gormRepository.Update(user)
}

func (userRepository *UserRepository) DeleteUser(id uint) error {
	return userRepository.gormRepository.Delete(id)
}

func (userRepository *UserRepository) GetUserByEmail(email string) (User, error) {
	var u User
	err := userRepository.gormRepository.DB.Where("email = ?", email).First(&u).Error
	return u, err
}
