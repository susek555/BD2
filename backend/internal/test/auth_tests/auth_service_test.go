package auth_tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"gorm.io/gorm"

	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
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
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)
		person := user.Person{Name: "John", Surname: "Doe"}

		in := user.CreateUserDTO{Email: "john@example.com", Password: "secret", Username: "john", PersonName: &person.Name, PersonSurname: &person.Surname, Selector: "P"}

		uRepo.On("GetByEmail", in.Email).Return(user.User{}, gorm.ErrRecordNotFound)
		uRepo.On("GetByUsername", in.Username).Return(user.User{}, gorm.ErrRecordNotFound)
		uRepo.On("Create", mock.Anything).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(ctx, in)

		assert.Empty(t, err)

		uRepo.AssertExpectations(t)
		rtSvc.AssertExpectations(t)
	})

	t.Run("when email taken it should return ErrEmailTaken", func(t *testing.T) {
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 1, Email: "marta@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)
		uRepo.On("GetByUsername", existing.Username).Return(existing, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(ctx, user.CreateUserDTO{Email: existing.Email})

		assert.NotEmpty(t, err, auth.ErrEmailTaken)
		uRepo.AssertExpectations(t)
	})
}

func TestService_Register_Company(t *testing.T) {
	t.Run("email not taken - should return access and refresh tokens", func(t *testing.T) {
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)
		company := user.Company{Name: "Awesome Name", NIP: "123233234234"}

		in := user.CreateUserDTO{
			Email:       "john@example.com",
			Password:    "secret",
			Username:    "john",
			CompanyName: &company.Name,
			CompanyNIP:  &company.NIP,
			Selector:    "C",
		}

		uRepo.On("GetByEmail", in.Email).Return(user.User{}, errors.New("not found"))
		uRepo.On("GetByUsername", in.Username).Return(user.User{}, errors.New("not found"))
		uRepo.On("GetByCompanyNip", *in.CompanyNIP).Return(user.User{}, errors.New("not found"))
		uRepo.On("Create", mock.Anything).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(ctx, in)

		assert.Empty(t, err)

		uRepo.AssertExpectations(t)
		rtSvc.AssertExpectations(t)
	})

	t.Run("when email taken it should return ErrEmailTaken", func(t *testing.T) {
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 1, Email: "john@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)
		uRepo.On("GetByUsername", existing.Username).Return(existing, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Register(ctx, user.CreateUserDTO{Email: existing.Email})

		assert.NotEmpty(t, err, auth.ErrEmailTaken)
		uRepo.AssertExpectations(t)
	})
}

func TestService_Login(t *testing.T) {
	ctx := context.Background()
	validIn := auth.LoginInput{Login: "john@example.com", Password: "secret"}

	t.Run("valid credentials – returns access & refresh", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 1, Email: validIn.Login, Password: hashPass(t, validIn.Password)}
		uRepo.EXPECT().GetByEmail(validIn.Login).Return(existing, nil)

		rtSvc.EXPECT().Create(mock.AnythingOfType("*refresh_token.RefreshToken")).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		access, refresh, err, user_ := svc.Login(ctx, validIn)
		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		uid, err := jwt.NewJWTVerifier(string(jwtKey)).VerifyToken(access)
		assert.NoError(t, err)
		assert.Equal(t, *user_, existing)
		assert.Equal(t, int64(existing.ID), uid)
	})

	t.Run("unknown e‑mail – ErrInvalidCredentials", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		uRepo.EXPECT().
			GetByEmail(validIn.Login).
			Return(user.User{}, gorm.ErrRecordNotFound)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, _, err, _ := svc.Login(ctx, validIn)
		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("wrong password – ErrInvalidCredentials", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		badPassUser := user.User{ID: 2, Email: validIn.Login, Password: hashPass(t, "invalidPass")}
		uRepo.EXPECT().GetByEmail(validIn.Login).Return(badPassUser, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, _, err, _ := svc.Login(ctx, validIn)
		assert.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("refresh‑token save fails – propagates error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 3, Email: validIn.Login, Password: hashPass(t, validIn.Password)}
		uRepo.EXPECT().GetByEmail(validIn.Login).Return(existing, nil)

		rtSvc.
			EXPECT().
			Create(mock.AnythingOfType("*refresh_token.RefreshToken")).
			Return(errors.New("db down"))

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, _, err, _ := svc.Login(ctx, validIn)
		assert.EqualError(t, err, "error - create refresh token")
	})
}

func TestService_Refresh(t *testing.T) {
	ctx := context.Background()
	oldToken := "old_refresh"

	baseRT := refresh_token.RefreshToken{
		ID:         101,
		Token:      oldToken,
		UserId:     1,
		ExpiryDate: time.Now().Add(24 * time.Hour),
		User:       user.User{ID: 1, Email: "john@example.com"},
	}

	t.Run("happy‑path – returns new access", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().FindByToken(ctx, oldToken).Return(&baseRT, nil)
		rtSvc.EXPECT().VerifyExpiration(ctx, &baseRT).Return(&baseRT, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		access, err := svc.Refresh(ctx, oldToken)

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
			FindByToken(ctx, oldToken).
			Return(nil, errors.New("sql: no rows"))

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, err := svc.Refresh(ctx, oldToken)
		assert.EqualError(t, err, "invalid refresh token")
	})

	t.Run("token expired", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		expired := errors.New("refresh token expired")

		rtSvc.EXPECT().FindByToken(ctx, oldToken).Return(&baseRT, nil)
		rtSvc.EXPECT().VerifyExpiration(ctx, &baseRT).Return(nil, expired)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		_, err := svc.Refresh(ctx, oldToken)
		assert.EqualError(t, err, expired.Error())
	})
}

func TestService_Logout(t *testing.T) {
	ctx := context.Background()
	userID := uint(8)
	rt := refresh_token.RefreshToken{ID: 42, Token: "r1"}

	t.Run("allDevices=true – deleteByUserID", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			DeleteByUserID(ctx, userID).
			Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}
		rtSvc.EXPECT().FindByToken(ctx, "gy").Return(&rt, nil)
		err := svc.Logout(ctx, userID, "gy", true)
		assert.NoError(t, err)
	})

	t.Run("single device – valid token", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().FindByToken(ctx, rt.Token).Return(&rt, nil)
		rtSvc.EXPECT().Delete(rt.ID).Return(nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(ctx, userID, rt.Token, false)
		assert.NoError(t, err)
	})

	t.Run("single device – no token found", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(ctx, userID, "", false)
		assert.EqualError(t, err, "refresh token required")
	})

	t.Run("single device – FindByToken returns error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			FindByToken(ctx, rt.Token).
			Return(nil, errors.New("refresh token not found"))

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(ctx, userID, rt.Token, false)
		assert.EqualError(t, err, "refresh token not found")
	})

	t.Run("allDevices – DeleteByUserID returns error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		rtSvc.EXPECT().
			DeleteByUserID(ctx, userID).
			Return(errors.New("db down"))
		rtSvc.EXPECT().FindByToken(ctx, rt.Token).Return(&rt, nil)

		svc := &auth.AuthService{Repo: uRepo, RefreshTokenService: rtSvc, JwtKey: jwtKey}

		err := svc.Logout(ctx, userID, rt.Token, true)
		assert.EqualError(t, err, "db down")
	})
}
