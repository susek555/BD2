package user

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type UserRepositoryInterface interface {
	generic.CRUDRepository[User]
	GetByEmail(email string) (User, error)
	GetByCompanyNip(nip string) (User, error)
	GetByUsername(username string) (User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		subtype := user.GetSubtype()
		subtype.SetUserID(user.ID)
		return subtype.SaveSubtype(tx)
	})
}

func (r *UserRepository) GetAll() ([]User, error) {
	var users []User
	err := r.DB.Preload("Company").Preload("Person").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, err
}

func (r *UserRepository) GetById(id uint) (User, error) {
	var user User
	err := r.DB.Preload("Company").Preload("Person").First(&user, id).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (User, error) {
	var u User
	err := r.DB.Preload("Company").Preload("Person").Where("email = ?", email).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetByUsername(username string) (User, error) {
	var u User
	err := r.DB.Preload("Company").Preload("Person").Where("username = ?", username).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetByCompanyNip(nip string) (User, error) {
	var u User
	err := r.DB.
		Joins("JOIN companies ON companies.user_id = users.id").
		Where("companies.nip = ?", nip).
		Preload("Company").
		Preload("Person").
		First(&u).Error

	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *UserRepository) Update(user *User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		subtype := user.GetSubtype()
		return subtype.SaveSubtype(tx)
	})
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		return r.DB.Delete(&User{}, id).Error
	})
}
