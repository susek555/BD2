package auth_tests

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
	"gorm.io/gorm"
)

var jwtKey = []byte("test-secret")

func hashPass(t *testing.T, raw string) string {
	hashed, err := passwords.Hash(raw)
	if err != nil {
		t.Fatalf("cannot hash password: %v", err)
	}
	return hashed
}

func TestService_Register_Person(t *testing.T) {
	t.Run("email not taken - should return access and refresh tokens", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)
		person := models.Person{Name: "John", Surname: "Doe"}

		in := user.CreateUserDTO{Email: "john@example.com", Password: "secret", Username: "john", PersonName: &person.Name, PersonSurname: &person.Surname, Selector: "P"}

		uRepo.On("GetByEmail", in.Email).Return(models.User{}, gorm.ErrRecordNotFound)
		uRepo.On("GetByUsername", in.Username).Return(models.User{}, gorm.ErrRecordNotFound)
		uRepo.On("Create", mock.Anything).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(in)

		assert.Empty(t, err)

		uRepo.AssertExpectations(t)
		rtSvc.AssertExpectations(t)
	})

	t.Run("when email taken it should return ErrEmailTaken", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := models.User{ID: 1, Email: "marta@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)
		uRepo.On("GetByUsername", existing.Username).Return(existing, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(user.CreateUserDTO{Email: existing.Email})

		assert.NotEmpty(t, err, auth.ErrEmailTaken)
		uRepo.AssertExpectations(t)
	})
}

func TestService_Register_Company(t *testing.T) {
	t.Run("email not taken - should return access and refresh tokens", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)
		company := models.Company{Name: "Awesome Name", Nip: "123233234234"}

		in := user.CreateUserDTO{
			Email:       "john@example.com",
			Password:    "secret",
			Username:    "john",
			CompanyName: &company.Name,
			CompanyNIP:  &company.Nip,
			Selector:    "C",
		}

		uRepo.On("GetByEmail", in.Email).Return(models.User{}, errors.New("not found"))
		uRepo.On("GetByUsername", in.Username).Return(models.User{}, errors.New("not found"))
		uRepo.On("GetByCompanyNip", *in.CompanyNIP).Return(models.User{}, errors.New("not found"))
		uRepo.On("Create", mock.Anything).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(in)

		assert.Empty(t, err)

		uRepo.AssertExpectations(t)
		rtSvc.AssertExpectations(t)
	})

	t.Run("when email taken it should return ErrEmailTaken", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := models.User{ID: 1, Email: "john@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)
		uRepo.On("GetByUsername", existing.Username).Return(existing, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(user.CreateUserDTO{Email: existing.Email})

		assert.NotEmpty(t, err, auth.ErrEmailTaken)
		uRepo.AssertExpectations(t)
	})
}

func TestService_Login(t *testing.T) {
	validIn := auth.LoginInput{Login: "john@example.com", Password: "secret"}

	t.Run("valid credentials - returns access & refresh", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := models.User{ID: 1, Email: validIn.Login, Password: hashPass(t, validIn.Password)}
		uRepo.EXPECT().GetByEmail(validIn.Login).Return(existing, nil)

		rtSvc.EXPECT().Create(mock.AnythingOfType("*models.RefreshToken")).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		access, refresh, user_, err := svc.Login(validIn)
		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		uid, err := jwt.NewJWTVerifier(string(jwtKey)).VerifyToken(access)
		assert.NoError(t, err)
		assert.Equal(t, *user_, existing)
		assert.Equal(t, int64(existing.ID), uid)
	})

	t.Run("unknown email - ErrInvalidCredentials", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		uRepo.EXPECT().
			GetByEmail(validIn.Login).
			Return(models.User{}, gorm.ErrRecordNotFound)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, _, _, err := svc.Login(validIn)
		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("wrong password - ErrInvalidCredentials", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		badPassUser := models.User{ID: 2, Email: validIn.Login, Password: hashPass(t, "invalidPass")}
		uRepo.EXPECT().GetByEmail(validIn.Login).Return(badPassUser, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, _, _, err := svc.Login(validIn)
		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("refresh-token save fails - propagates error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := models.User{ID: 3, Email: validIn.Login, Password: hashPass(t, validIn.Password)}
		uRepo.EXPECT().GetByEmail(validIn.Login).Return(existing, nil)

		rtSvc.
			EXPECT().
			Create(mock.AnythingOfType("*models.RefreshToken")).
			Return(errors.New("db down"))

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, _, _, err := svc.Login(validIn)
		assert.EqualError(t, err, "error - create refresh token")
	})
}

func TestService_Refresh(t *testing.T) {
	oldToken := "old_refresh"

	baseRT := models.RefreshToken{
		ID:         101,
		Token:      oldToken,
		UserID:     1,
		ExpiryDate: time.Now().Add(24 * time.Hour),
		User:       &models.User{ID: 1, Email: "john@example.com"},
	}

	t.Run("happy‑path – returns new access", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().FindByToken(oldToken).Return(&baseRT, nil)
		rtSvc.EXPECT().VerifyExpiration(&baseRT).Return(&baseRT, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		access, err := svc.Refresh(oldToken)

		assert.NoError(t, err)
		assert.NotEmpty(t, access)

		uid, err := jwt.NewJWTVerifier(string(jwtKey)).VerifyToken(access)
		assert.NoError(t, err)
		assert.Equal(t, int64(baseRT.User.ID), uid)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			FindByToken(oldToken).
			Return(nil, errors.New("sql: no rows"))

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, err := svc.Refresh(oldToken)
		assert.EqualError(t, err, "invalid refresh token")
	})

	t.Run("token expired", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		expired := errors.New("refresh token expired")

		rtSvc.EXPECT().FindByToken(oldToken).Return(&baseRT, nil)
		rtSvc.EXPECT().VerifyExpiration(&baseRT).Return(nil, expired)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, err := svc.Refresh(oldToken)
		assert.EqualError(t, err, expired.Error())
	})
}

func TestService_Logout(t *testing.T) {
	userID := uint(8)
	rt := models.RefreshToken{ID: 42, Token: "r1"}

	t.Run("allDevices=true - deleteByUserID", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			DeleteByUserID(userID).
			Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}
		rtSvc.EXPECT().FindByToken("gy").Return(&rt, nil)
		err := svc.Logout(userID, "gy", true)
		assert.NoError(t, err)
	})

	t.Run("single device – valid token", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().FindByToken(rt.Token).Return(&rt, nil)
		rtSvc.EXPECT().Delete(rt.ID).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(userID, rt.Token, false)
		assert.NoError(t, err)
	})

	t.Run("single device – no token found", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(userID, "", false)
		assert.EqualError(t, err, "refresh token required")
	})

	t.Run("single device – FindByToken returns error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			FindByToken(rt.Token).
			Return(nil, errors.New("refresh token not found"))

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(userID, rt.Token, false)
		assert.EqualError(t, err, "refresh token not found")
	})

	t.Run("allDevices – DeleteByUserID returns error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			DeleteByUserID(userID).
			Return(errors.New("db down"))
		rtSvc.EXPECT().FindByToken(rt.Token).Return(&rt, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(userID, rt.Token, true)
		assert.EqualError(t, err, "db down")
	})
}
