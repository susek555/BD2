package refresh_token

import (
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

type RefreshTokenRepositoryInterface interface {
	generic.CRUDRepository[RefreshToken]
	FindByToken(token string) (*RefreshToken, error)
	FindByUserEmail(email string) ([]RefreshToken, error)
	FindByUserId(id uint) ([]RefreshToken, error)
	DeleteByUserId(id uint) error
}

type RefreshTokenRepository struct {
	repository *generic.GormRepository[RefreshToken]
}

func GetRefreshTokenRepository(dbHandle *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{repository: generic.GetGormRepository[RefreshToken](dbHandle)}
}

func (repo *RefreshTokenRepository) Create(token *RefreshToken) error {
	return repo.repository.Create(token)
}

func (repo *RefreshTokenRepository) GetAll() ([]RefreshToken, error) {
	return repo.repository.GetAll()
}

func (repo *RefreshTokenRepository) GetById(id uint) (*RefreshToken, error) {
	return repo.repository.GetById(id)
}

func (repo *RefreshTokenRepository) Update(token *RefreshToken) error {
	return repo.repository.Update(token)
}

func (repo *RefreshTokenRepository) Delete(id uint) error {
	return repo.repository.Delete(id)
}

func (repo *RefreshTokenRepository) FindByUserEmail(email string) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := repo.repository.
		DB.
		Joins("User").
		Where("users.email = ?", email).
		Find(&tokens).Error
	return tokens, err
}

func (repo *RefreshTokenRepository) FindByUserId(id uint) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := repo.repository.
		DB.
		Where("user_id = ?", id).
		Find(&tokens).Error
	return tokens, err
}

func (repo *RefreshTokenRepository) FindByToken(token string) (*RefreshToken, error) {
	var t RefreshToken
	err := repo.repository.
		DB.
		Where("token = ?", token).
		First(&t).Error
	return &t, err
}

func (repo *RefreshTokenRepository) DeleteByUserId(id uint) error {
	err := repo.repository.
		DB.
		Where("user_id = ?", id).
		Delete(&RefreshToken{}).
		Error
	return err
}
