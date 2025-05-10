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
	Register(ctx context.Context, in user.CreateUserDTO) map[string][]string
	Login(ctx context.Context, in LoginInput) (access, refresh string, err error, user *user.User)
	Refresh(ctx context.Context, refreshToken string) (access string, err error)
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

func (s *service) Register(ctx context.Context, in user.CreateUserDTO) map[string][]string {
	userModel, err := in.MapToUser()
	var errs = make(map[string][]string)
	if err != nil {
		errs["other"] = []string{err.Error()}
	}
	_, noUsername := s.repo.GetByUsername(in.Username)
	if noUsername == nil {
		errs["username"] = []string{"Username already taken"}
	}
	_, noEmail := s.repo.GetByEmail(in.Email)
	if noEmail == nil {
		errs["email"] = []string{"Email already taken"}
	}
	if in.Selector == "C" {
		_, noNip := s.repo.GetByCompanyNip(*in.CompanyNIP)
		if noNip == nil {
			errs["company_nip"] = []string{"NIP already taken"}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	if err := s.repo.Create(&userModel); err != nil {
		errs["other"] = []string{err.Error()}
	}
	return errs
}

func (s *service) Login(ctx context.Context, in LoginInput) (string, string, error, *user.User) {
	u, err := s.repo.GetByEmail(in.Login)
	if err != nil || u.ID == 0 {
		return "", "", ErrInvalidCredentials, &user.User{}
	}

	if !passwords.Match(in.Password, u.Password) {
		return "", "", ErrInvalidCredentials, &user.User{}
	}
	access, err := jwt.GenerateToken(u.Email, int64(u.ID), s.jwtKey, time.Now().Add(2*time.Minute))
	if err != nil {
		return "", "", err, &user.User{}
	}

	refresh, err := s.newRefreshToken(ctx, u.ID, u.Email)
	if err != nil {
		return "", "", err, &user.User{}
	}
	return access, refresh, nil, &u
}

func (s *service) Refresh(ctx context.Context, provided string) (string, error) {
	refresh, err := s.refreshTokenService.FindByToken(ctx, provided)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	if _, err := s.refreshTokenService.VerifyExpiration(ctx, refresh); err != nil {
		return "", err
	}

	access, err := jwt.GenerateToken(refresh.User.Email, int64(refresh.User.ID), s.jwtKey, time.Now().Add(2*time.Hour))
	if err != nil {
		return "", err
	}

	return access, nil
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
	err := s.refreshTokenService.Create(&refresh)
	if err != nil {
		return "", err
	}
	return token, nil
}
