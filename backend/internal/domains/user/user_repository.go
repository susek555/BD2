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

func GetUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.DB.Create(&user).Error; err != nil {
			return err
		}
		subtype := user.GetSubtype()
		subtype.SetUserID(user.ID)
		return subtype.SaveSubtype(r.DB)
	})
}

func (r *UserRepository) GetAll() ([]User, error) {
	var users []User
	err := r.DB.Preload("Company").Preload("Person").Find(&users).Error
	return users, err
}

func (r *UserRepository) GetById(id uint) (User, error) {
	var user User
	err := r.DB.Preload("Company").Preload("Person").Find(&user, id).Error
	return user, err
}

func (r *UserRepository) Update(user User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := r.DB.Save(user).Error; err != nil {
			return err
		}
		subtype := user.GetSubtype()
		return subtype.SaveSubtype(r.DB)
	})
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		return r.DB.Delete(&User{}, id).Error
	})

}

func (r *UserRepository) GetByEmail(email string) (User, error) {
	var u User
	err := r.DB.Preload("Comapny").Preload("Person").Where("email = ?", email).First(&u).Error
	return u, err
}
