package user

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepositoryInterface --output=../../test/mocks --case=snake --with-expecter
type UserRepositoryInterface interface {
	generic.CRUDRepository[models.User]
	GetByEmail(email string) (models.User, error)
	GetByCompanyNip(nip string) (models.User, error)
	GetByUsername(username string) (models.User, error)
	UpdatePassword(userID uint, newPassword string) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		subtype := user.GetSubtype()
		subtype.SetUserID(user.ID)
		return subtype.SaveSubtype(tx)
	})
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.buildBaseQuery().Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, err
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.buildBaseQuery().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var u models.User
	err := r.buildBaseQuery().Where("email = ?", email).First(&u).Error
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetByUsername(username string) (models.User, error) {
	var u models.User
	err := r.buildBaseQuery().Where("username = ?", username).First(&u).Error
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetByCompanyNip(nip string) (models.User, error) {
	var u models.User
	err := r.buildBaseQuery().
		Joins("JOIN companies ON companies.user_id = users.id").
		Where("companies.nip = ?", nip).
		First(&u).Error

	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		subtype := user.GetSubtype()
		return subtype.SaveSubtype(tx)
	})
}

func (r *UserRepository) UpdatePassword(userID uint, newPassword string) error {
	return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) buildBaseQuery() *gorm.DB {
	return r.DB.Preload("Person").Preload("Company")
}
