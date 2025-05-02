package user

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	generic.CRUDRepository[User]
	GetUserByEmail(email string) (User, error)
}

type UserRepository struct {
	*generic.GormRepository[User]
}

func GetUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{GormRepository: generic.GetGormRepository[User](db)}
}

func (userRepository *UserRepository) GetUserByEmail(email string) (User, error) {
	var u User
	err := userRepository.DB.Where("email = ?", email).First(&u).Error
	return u, err
}
