package auth

import (
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
)

type AuthServiceInterface interface {
	Register(in user.CreateUserDTO) map[string][]string
	Login(in LoginInput) (access, refresh string, user *models.User, err error)
	Refresh(refreshToken string) (access string, err error)
	Logout(userID uint, refreshToken string, allDevices bool) error
	ChangePassword(userID uint, oldPassword, newPassword string) map[string][]string
}

type AuthService struct {
	Repo                user.UserRepositoryInterface
	RefreshTokenService refresh_token.RefreshTokenServiceInterface
	JwtKey              []byte
}

func NewAuthService(repo user.UserRepositoryInterface, refreshTokenService refresh_token.RefreshTokenServiceInterface, jwtKey []byte) AuthServiceInterface {
	return &AuthService{
		Repo:                repo,
		RefreshTokenService: refreshTokenService,
		JwtKey:              jwtKey,
	}
}

func (s *AuthService) Register(in user.CreateUserDTO) map[string][]string {
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

func (s *AuthService) Login(in LoginInput) (string, string, *models.User, error) {
	u, err := s.Repo.GetByEmail(in.Login)
	if err != nil || u.ID == 0 {
		return "", "", &models.User{}, ErrInvalidCredentials
	}

	if !passwords.Match(in.Password, u.Password) {
		return "", "", &models.User{}, ErrInvalidCredentials
	}
	access, err := jwt.GenerateToken(u.Email, int64(u.ID), s.JwtKey, time.Now().Add(AccessTokenExpirationTime))
	if err != nil {
		return "", "", &models.User{}, err
	}

	refresh, err := s.newRefreshToken(u.ID, u.Email)
	if err != nil {
		return "", "", &models.User{}, err
	}
	return access, refresh, &u, nil
}

func (s *AuthService) Refresh(provided string) (string, error) {
	refresh, err := s.RefreshTokenService.FindByToken(provided)
	if err != nil {
		return "", ErrInvalidRefreshToken
	}

	if _, err := s.RefreshTokenService.VerifyExpiration(refresh); err != nil {
		return "", ErrRefreshTokenExpired
	}

	access, err := jwt.GenerateToken(refresh.User.Email, int64(refresh.User.ID), s.JwtKey, time.Now().Add(2*time.Hour))
	if err != nil {
		return "", err
	}

	return access, nil
}

func (s *AuthService) Logout(userID uint, provided string, allDevices bool) error {
	if provided == "" {
		return ErrRefreshTokenRequired
	}
	refresh, err := s.RefreshTokenService.FindByToken(provided)
	if err != nil {
		return ErrRefreshTokenNotFound
	}
	if allDevices {
		return s.RefreshTokenService.DeleteByUserID(userID)
	}
	return s.RefreshTokenService.Delete(refresh.ID)
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) map[string][]string {
	user_, err := s.Repo.GetByID(userID)
	if err != nil {
		return map[string][]string{"other": {ErrUserNotFound.Error()}}
	}
	if !passwords.Match(oldPassword, user_.Password) {
		return map[string][]string{"old_password": {ErrInvalidOldPassword.Error()}}
	}
	hashedPassword, err := passwords.Hash(newPassword)
	if err != nil {
		return map[string][]string{"other": {err.Error()}}
	}
	if err := s.Repo.UpdatePassword(userID, hashedPassword); err != nil {
		return map[string][]string{"other": {err.Error()}}
	}
	return nil
}

func (s *AuthService) newRefreshToken(userID uint, userEmail string) (string, error) {
	token, _ := jwt.GenerateToken(userEmail, int64(userID), s.JwtKey, time.Now().Add(RefreshTokenExpirationTime))
	refresh := models.RefreshToken{
		Token:      token,
		UserID:     userID,
		ExpiryDate: time.Now().UTC().Add(30 * 24 * time.Hour),
	}
	err := s.RefreshTokenService.Create(&refresh)
	if err != nil {
		return "", ErrCreateRefreshToken
	}
	return token, nil
}
