package refresh_token

import (
	"errors"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/models"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
)

//go:generate mockery --name=RefreshTokenServiceInterface --output=../../test/mocks --case=snake --with-expecter
type RefreshTokenServiceInterface interface {
	generic.CRUDService[models.RefreshToken]

	FindByToken(token string) (*models.RefreshToken, error)
	FindByUserEmail(email string) ([]models.RefreshToken, error)
	FindByUserId(id uint) ([]models.RefreshToken, error)
	VerifyExpiration(token *models.RefreshToken) (*models.RefreshToken, error)
	DeleteByUserId(userID uint) error
}

type RefreshTokenService struct {
	generic.GenericService[models.RefreshToken, RefreshTokenRepositoryInterface]
}

func NewRefreshTokenService(repo RefreshTokenRepositoryInterface) RefreshTokenServiceInterface {
	return &RefreshTokenService{
		GenericService: generic.GenericService[models.RefreshToken, RefreshTokenRepositoryInterface]{
			Repo: repo,
		},
	}
}

func (s *RefreshTokenService) FindByToken(token string) (*models.RefreshToken, error) {
	return s.Repo.FindByToken(token)
}

func (s *RefreshTokenService) FindByUserEmail(email string) ([]models.RefreshToken, error) {
	return s.Repo.FindByUserEmail(email)
}

func (s *RefreshTokenService) FindByUserId(id uint) ([]models.RefreshToken, error) {
	return s.Repo.FindByUserId(id)
}

func (s *RefreshTokenService) VerifyExpiration(token *models.RefreshToken) (*models.RefreshToken, error) {
	if token.ExpiryDate.Before(time.Now().UTC()) {
		_ = s.Repo.Delete(token.ID)
		return nil, errors.New("refresh token expired")
	}
	return token, nil
}

func (s *RefreshTokenService) DeleteByUserId(userID uint) error {
	return s.Repo.DeleteByUserId(userID)
}
