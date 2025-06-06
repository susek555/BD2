package refresh_token

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepositoryInterface interface {
	generic.CRUDRepository[models.RefreshToken]
	FindByToken(token string) (*models.RefreshToken, error)
	FindByUserEmail(email string) ([]models.RefreshToken, error)
	FindByUserID(id uint) ([]models.RefreshToken, error)
	DeleteByUserID(id uint) error
}

type RefreshTokenRepository struct {
	repository *generic.GormRepository[models.RefreshToken]
}

func NewRefreshTokenRepository(dbHandle *gorm.DB) RefreshTokenRepositoryInterface {
	return &RefreshTokenRepository{repository: generic.GetGormRepository[models.RefreshToken](dbHandle)}
}

func (repo *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	return repo.repository.Create(token)
}

func (repo *RefreshTokenRepository) GetAll() ([]models.RefreshToken, error) {
	return repo.repository.GetAll()
}

func (repo *RefreshTokenRepository) GetByID(id uint) (*models.RefreshToken, error) {
	return repo.repository.GetByID(id)
}

func (repo *RefreshTokenRepository) Update(token *models.RefreshToken) error {
	return repo.repository.Update(token)
}

func (repo *RefreshTokenRepository) Delete(id uint) error {
	return repo.repository.Delete(id)
}

func (repo *RefreshTokenRepository) FindByUserEmail(email string) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := repo.repository.
		DB.
		Joins("User").
		Where("users.email = ?", email).
		Find(&tokens).Error
	return tokens, err
}

func (repo *RefreshTokenRepository) FindByUserID(id uint) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := repo.repository.
		DB.
		Where("user_id = ?", id).
		Find(&tokens).Error
	return tokens, err
}

func (repo *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := repo.repository.
		DB.
		Preload("User").
		Where("token = ?", token).
		First(&t).Error
	return &t, err
}

func (repo *RefreshTokenRepository) DeleteByUserID(id uint) error {
	err := repo.repository.
		DB.
		Where("user_id = ?", id).
		Delete(&models.RefreshToken{}).
		Error
	return err
}
