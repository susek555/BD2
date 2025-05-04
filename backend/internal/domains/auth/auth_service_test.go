package auth

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"gorm.io/gorm"
	"testing"

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

		in := user.CreateUserDTO{Email: "john@example.com", Password: "sekret", Username: "john", PersonName: &person.Name, PersonSurname: &person.Surname, Selector: "P"}

		uRepo.On("GetByEmail", in.Email).Return(user.User{}, gorm.ErrRecordNotFound)
		uRepo.On("Create", mock.Anything).Return(nil)
		rtSvc.On("Create", mock.Anything).Return(nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		access, refresh, err := svc.Register(ctx, in)

		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		uRepo.AssertExpectations(t)
		rtSvc.AssertExpectations(t)
	})

	t.Run("when email taken it should return ErrEmailTaken", func(t *testing.T) {
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 1, Email: "marta@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		_, _, err := svc.Register(ctx, user.CreateUserDTO{Email: existing.Email})

		assert.ErrorIs(t, err, ErrEmailTaken)
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
			Password:    "sekret",
			Username:    "john",
			CompanyName: &company.Name,
			CompanyNIP:  &company.NIP,
			Selector:    "C",
		}

		uRepo.On("GetByEmail", in.Email).Return(user.User{}, errors.New("not found"))
		uRepo.On("Create", mock.Anything).Return(nil)
		rtSvc.On("Create", mock.Anything).Return(nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		access, refresh, err := svc.Register(ctx, in)

		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		uRepo.AssertExpectations(t)
		rtSvc.AssertExpectations(t)
	})

	t.Run("when email taken it should return ErrEmailTaken", func(t *testing.T) {
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 1, Email: "john@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		_, _, err := svc.Register(ctx, user.CreateUserDTO{Email: existing.Email})

		assert.ErrorIs(t, err, ErrEmailTaken)
		uRepo.AssertExpectations(t)
	})
}

func TestService_Login(t *testing.T) {
	ctx := context.Background()
	validIn := LoginInput{Email: "john@example.com", Password: "sekret"}

	t.Run("valid credentials – returns access & refresh", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 1, Email: validIn.Email, Password: hashPass(t, validIn.Password)}
		uRepo.EXPECT().GetByEmail(validIn.Email).Return(existing, nil)

		rtSvc.EXPECT().Create(mock.AnythingOfType("refresh_token.RefreshToken")).Return(nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		access, refresh, err := svc.Login(ctx, validIn)
		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)

		uid, err := jwt.NewJWTVerifier(string(jwtKey)).VerifyToken(access)
		assert.NoError(t, err)
		assert.Equal(t, int64(existing.ID), uid)
	})

	t.Run("unknown e‑mail – ErrInvalidCredentials", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		uRepo.EXPECT().
			GetByEmail(validIn.Email).
			Return(user.User{}, gorm.ErrRecordNotFound)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		_, _, err := svc.Login(ctx, validIn)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("wrong password – ErrInvalidCredentials", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		badPassUser := user.User{ID: 2, Email: validIn.Email, Password: hashPass(t, "invalidPass")}
		uRepo.EXPECT().GetByEmail(validIn.Email).Return(badPassUser, nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		_, _, err := svc.Login(ctx, validIn)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("refresh‑token save fails – propagates error", func(t *testing.T) {
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)

		existing := user.User{ID: 3, Email: validIn.Email, Password: hashPass(t, validIn.Password)}
		uRepo.EXPECT().GetByEmail(validIn.Email).Return(existing, nil)

		rtSvc.
			EXPECT().
			Create(mock.AnythingOfType("refresh_token.RefreshToken")).
			Return(errors.New("db down"))

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		_, _, err := svc.Login(ctx, validIn)
		assert.EqualError(t, err, "db down")
	})
}
