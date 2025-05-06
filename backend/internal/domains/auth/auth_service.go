package auth

import (
	"context"
	"errors"

	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"gorm.io/gorm"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

var (
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service interface {
	Register(ctx context.Context, in user.CreateUserDTO) error
	Login(ctx context.Context, in LoginInput) (access, refresh string, err error)
	Refresh(ctx context.Context, refreshToken string) (access string, refresh string, err error)
	Logout(ctx context.Context, userID uint, refreshToken string, allDevices bool) error
}

type service struct {
	repo                user.UserRepositoryInterface
	refreshTokenService refresh_token.RefreshTokenServiceInterface
	jwtKey              []byte
}

func NewService(db *gorm.DB, jwtKey []byte) Service {
	userRepo := user.NewUserRepository(db)
	refreshTokenService := refresh_token.NewRefreshTokenService(db)
	return &service{
		repo:                userRepo,
		refreshTokenService: refreshTokenService,
		jwtKey:              jwtKey,
	}
}

func (s *service) Register(ctx context.Context, in user.CreateUserDTO) error {
	userModel, err := in.MapToUser()
	if err != nil {
		return err
	}
	if err := s.repo.Create(userModel); err != nil {
		return err
	}
	return nil
}

func (s *service) Login(ctx context.Context, in LoginInput) (string, string, error) {
	u, err := s.repo.GetByEmail(in.Email)
	if err != nil || u.ID == 0 {
		return "", "", ErrInvalidCredentials
	}

	if !passwords.Match(in.Password, u.Password) {
		return "", "", ErrInvalidCredentials
	}
	access, err := jwt.GenerateToken(u.Email, int64(u.ID), s.jwtKey, time.Now().Add(2*time.Hour))
	if err != nil {
		return "", "", err
	}

	refresh, err := s.newRefreshToken(ctx, u.ID, u.Email)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *service) Refresh(ctx context.Context, provided string) (string, string, error) {
	refresh, err := s.refreshTokenService.FindByToken(ctx, provided)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	if _, err := s.refreshTokenService.VerifyExpiration(ctx, refresh); err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateToken(refresh.User.Email, int64(refresh.User.ID), s.jwtKey, time.Now().Add(2*time.Hour))
	if err != nil {
		return "", "", err
	}
	_ = s.refreshTokenService.Delete(refresh.ID)
	newRefresh, err := s.newRefreshToken(ctx, refresh.UserId, refresh.User.Email)
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}

func (s *service) Logout(ctx context.Context, userID uint, provided string, allDevices bool) error {
	if allDevices {
		return s.refreshTokenService.DeleteByUserID(ctx, userID)
	}
	if provided == "" {
		return errors.New("refresh token required")
	}

	refresh, err := s.refreshTokenService.FindByToken(ctx, provided)
	if err != nil {
		return err
	}
	return s.refreshTokenService.Delete(refresh.ID)
}

func (s *service) newRefreshToken(ctx context.Context, userId uint, userEmail string) (string, error) {
	token, _ := jwt.GenerateToken(userEmail, int64(userId), s.jwtKey, time.Now().Add(30*24*time.Hour))
	refresh := refresh_token.RefreshToken{
		Token:      token,
		UserId:     userId,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
	}
	err := s.refreshTokenService.Create(refresh)
	if err != nil {
		return "", err
	}
	return token, nil
}
