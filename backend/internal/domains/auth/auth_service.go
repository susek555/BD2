package auth

import (
	"context"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"gorm.io/gorm"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, in user.CreateUserDTO) map[string][]string
	Login(ctx context.Context, in LoginInput) (access, refresh string, err error, user *user.User)
	Refresh(ctx context.Context, refreshToken string) (access string, err error)
	Logout(ctx context.Context, userID uint, refreshToken string, allDevices bool) error
}

type AuthService struct {
	Repo                user.UserRepositoryInterface
	RefreshTokenService refresh_token.RefreshTokenServiceInterface
	JwtKey              []byte
}

func NewService(db *gorm.DB, jwtKey []byte) AuthServiceInterface {
	userRepo := user.NewUserRepository(db)
	refreshTokenService := refresh_token.NewRefreshTokenService(db)
	return &AuthService{
		Repo:                userRepo,
		RefreshTokenService: refreshTokenService,
		JwtKey:              jwtKey,
	}
}

func (s *AuthService) Register(ctx context.Context, in user.CreateUserDTO) map[string][]string {
	userModel, err := in.MapToUser()
	var errs = make(map[string][]string)
	if err != nil {
		errs["other"] = []string{err.Error()}
	}
	_, noUsername := s.Repo.GetByUsername(in.Username)
	if noUsername == nil {
		errs["username"] = []string{ErrUsernameTaken.Error()}
	}
	_, noEmail := s.Repo.GetByEmail(in.Email)
	if noEmail == nil {
		errs["email"] = []string{ErrEmailTaken.Error()}
	}
	if in.Selector == "C" {
		_, noNip := s.Repo.GetByCompanyNip(*in.CompanyNIP)
		if noNip == nil {
			errs["company_nip"] = []string{ErrNipAlreadyTaken.Error()}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	if err := s.Repo.Create(userModel); err != nil {
		errs["other"] = []string{err.Error()}
	}
	return errs
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (string, string, error, *user.User) {
	u, err := s.Repo.GetByEmail(in.Login)
	if err != nil || u.ID == 0 {
		return "", "", ErrInvalidCredentials, &user.User{}
	}

	if !passwords.Match(in.Password, u.Password) {
		return "", "", ErrInvalidCredentials, &user.User{}
	}
	access, err := jwt.GenerateToken(u.Email, int64(u.ID), s.JwtKey, time.Now().Add(2*time.Minute))
	if err != nil {
		return "", "", err, &user.User{}
	}

	refresh, err := s.newRefreshToken(ctx, u.ID, u.Email)
	if err != nil {
		return "", "", err, &user.User{}
	}
	return access, refresh, nil, &u
}

func (s *AuthService) Refresh(ctx context.Context, provided string) (string, error) {
	refresh, err := s.RefreshTokenService.FindByToken(ctx, provided)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	if _, err := s.RefreshTokenService.VerifyExpiration(ctx, refresh); err != nil {
		return "", ErrRefreshTokenExpired
	}

	access, err := jwt.GenerateToken(refresh.User.Email, int64(refresh.User.ID), s.JwtKey, time.Now().Add(2*time.Hour))
	if err != nil {
		return "", err
	}

	return access, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint, provided string, allDevices bool) error {
	if provided == "" {
		return ErrRefreshTokenRequired
	}
	refresh, err := s.RefreshTokenService.FindByToken(ctx, provided)
	if err != nil {
		return ErrRefreshTokenNotFound
	}
	if allDevices {
		return s.RefreshTokenService.DeleteByUserID(ctx, userID)
	}
	return s.RefreshTokenService.Delete(refresh.ID)
}

func (s *AuthService) newRefreshToken(ctx context.Context, userId uint, userEmail string) (string, error) {
	token, _ := jwt.GenerateToken(userEmail, int64(userId), s.JwtKey, time.Now().Add(30*24*time.Hour))
	refresh := refresh_token.RefreshToken{
		Token:      token,
		UserId:     userId,
		ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
	}
	err := s.RefreshTokenService.Create(&refresh)
	if err != nil {
		return "", ErrCreateRefreshToken
	}
	return token, nil
}
