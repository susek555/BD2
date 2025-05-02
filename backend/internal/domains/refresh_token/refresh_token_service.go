package refresh_token

import (
	"context"
	"errors"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/generic"
	"gorm.io/gorm"
	"time"
)

type RefreshTokenServiceInterface interface {
	generic.CRUDService[RefreshToken]

	FindByToken(ctx context.Context, token string) (RefreshToken, error)
	FindByUserEmail(ctx context.Context, email string) (RefreshToken, error)
	VerifyExpiration(ctx context.Context, token RefreshToken) (RefreshToken, error)
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

func (s *RefreshTokenService) FindByToken(ctx context.Context, token string) (RefreshToken, error) {
	return s.Repo.FindByToken(token)
}

func (s *RefreshTokenService) FindByUserEmail(ctx context.Context, email string) (RefreshToken, error) {
	return s.Repo.FindByUserEmail(email)
}

func (s *RefreshTokenService) VerifyExpiration(ctx context.Context, token RefreshToken) (RefreshToken, error) {
	if token.ExpiryDate.Before(time.Now()) {
		_ = s.Repo.Delete(token.ID)
		return RefreshToken{}, errors.New("refresh token expired")
	}
	return token, nil
}

func (s *RefreshTokenService) DeleteByUserID(ctx context.Context, userID uint) error {
	return s.Repo.DeleteByUserId(userID)
}
