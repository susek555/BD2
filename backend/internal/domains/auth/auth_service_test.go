package auth

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
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
