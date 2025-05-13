package refresh_token

import (
	"context"
	"errors"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
)

//go:generate mockery --name=RefreshTokenServiceInterface --output=../../test/mocks --case=snake --with-expecter
type RefreshTokenServiceInterface interface {
	generic.CRUDService[RefreshToken]

	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	FindByUserEmail(ctx context.Context, email string) ([]RefreshToken, error)
	FindByUserId(ctx context.Context, id uint) ([]RefreshToken, error)
	VerifyExpiration(ctx context.Context, token *RefreshToken) (*RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID uint) error
}

type RefreshTokenService struct {
	generic.GenericService[RefreshToken, *RefreshTokenRepository]
}

func NewRefreshTokenService(db *gorm.DB) *RefreshTokenService {
	repo := GetRefreshTokenRepository(db)

	return &RefreshTokenService{
		GenericService: generic.GenericService[RefreshToken, *RefreshTokenRepository]{
			Repo: repo,
		},
	}
}

func (s *RefreshTokenService) FindByToken(ctx context.Context, token string) (*RefreshToken, error) {
	return s.Repo.FindByToken(token)
}

func (s *RefreshTokenService) FindByUserEmail(ctx context.Context, email string) ([]RefreshToken, error) {
	return s.Repo.FindByUserEmail(email)
}

func (s *RefreshTokenService) FindByUserId(ctx context.Context, id uint) ([]RefreshToken, error) {
	return s.Repo.FindByUserId(id)
}

func (s *RefreshTokenService) VerifyExpiration(ctx context.Context, token *RefreshToken) (*RefreshToken, error) {
	if token.ExpiryDate.Before(time.Now()) {
		_ = s.Repo.Delete(token.ID)
		return nil, errors.New("refresh token expired")
	}
	return token, nil
}

func (s *RefreshTokenService) DeleteByUserID(ctx context.Context, userID uint) error {
	return s.Repo.DeleteByUserId(userID)
}
