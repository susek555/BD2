package user

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	generic.CRUDRepository[User]
	GetByEmail(email string) (User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func GetUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user User) error {
	if err := r.DB.Create(&user).Error; err != nil {
		return err
	}
	if subtype := user.GetSubtype(); subtype != nil {
		subtype.SetUserID(user.ID)
		return subtype.SaveSubtype(r.DB)
	}
	return nil
}

func (r *UserRepository) GetAll() ([]User, error) {
	var users []User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetById(id uint) (User, error) {
	var user User
	err := r.DB.Find(&user, id).Error
	return user, err
}

func (r *UserRepository) Update(user User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return err
	}
	if subtype := user.GetSubtype(); subtype != nil {
		return subtype.SaveSubtype(r.DB)
	}
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&User{}, id).Error
}

func (r *UserRepository) GetByEmail(email string) (User, error) {
	var u User
	err := r.DB.Where("email = ?", email).First(&u).Error
	return u, err
}
