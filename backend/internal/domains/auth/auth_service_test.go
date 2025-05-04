package auth

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/susek555/BD2/car-dealer-api/internal/test/mocks"
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

/*
	---  TESTY  --------------------------------------------------------------
*/

func TestService_Register_Person(t *testing.T) {
	t.Run("email not taken - should return access and refresh tokens", func(t *testing.T) {
		ctx := context.Background()
		uRepo := mocks.NewUserRepositoryInterface(t)
		rtSvc := mocks.NewRefreshTokenServiceInterface(t)
		person := user.Person{Name: "John", Surname: "Doe"}

		in := user.CreateUserDTO{Email: "john@example.com", Password: "sekret", Username: "john", PersonName: &person.Name, PersonSurname: &person.Surname, Selector: "P"}

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

		existing := user.User{ID: 1, Email: "marta@example.com"}
		uRepo.On("GetByEmail", existing.Email).Return(existing, nil)

		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}

		_, _, err := svc.Register(ctx, user.CreateUserDTO{Email: existing.Email})

		assert.ErrorIs(t, err, ErrEmailTaken)
		uRepo.AssertExpectations(t)
	})
}

//func TestService_Login(t *testing.T) {
//	ctx := context.Background()
//	input := LoginInput{Email: "marta@example.com", Password: "sekret"}
//
//	t.Run("poprawne dane – daje access + refresh", func(t *testing.T) {
//		uRepo := new(userRepoMock)
//		rtSvc := new(refreshServiceMock)
//
//		hashed := hashPass(t, input.Password)
//		existing := &user.User{ID: 1, Email: input.Email, Password: hashed}
//
//		uRepo.On("GetByEmail", input.Email).Return(existing, nil)
//		rtSvc.On("Create", mock.AnythingOfType("refresh_token.RefreshToken")).Return(nil)
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		access, refresh, err := svc.Login(ctx, input)
//
//		assert.NoError(t, err)
//		assert.NotEmpty(t, access)
//		assert.NotEmpty(t, refresh)
//	})
//
//	t.Run("złe hasło – ErrInvalidCredentials", func(t *testing.T) {
//		uRepo := new(userRepoMock)
//		rtSvc := new(refreshServiceMock)
//
//		hashed := hashPass(t, "innehaslo")
//		existing := &user.User{ID: 1, Email: input.Email, Password: hashed}
//
//		uRepo.On("GetByEmail", input.Email).Return(existing, nil)
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		_, _, err := svc.Login(ctx, input)
//		assert.ErrorIs(t, err, ErrInvalidCredentials)
//	})
//}
//
//func TestService_Refresh(t *testing.T) {
//	ctx := context.Background()
//	oldRefresh := "old_token"
//
//	t.Run("poprawny refresh – nowe access + refresh", func(t *testing.T) {
//		uRepo := new(userRepoMock) // nie będzie używane, ale struct wymaga
//		rtSvc := new(refreshServiceMock)
//
//		userID := uint(1)
//		userEmail := "marta@example.com"
//
//		rt := &refresh_token.RefreshToken{
//			ID:         123,
//			Token:      oldRefresh,
//			UserId:     userID,
//			ExpiryDate: time.Now().Add(24 * time.Hour),
//			User:       user.User{ID: userID, Email: userEmail},
//		}
//
//		rtSvc.On("FindByToken", oldRefresh).Return(rt, nil)
//		rtSvc.On("VerifyExpiration", rt).Return(rt, nil)
//		rtSvc.On("Delete", rt.ID).Return(nil)
//		rtSvc.On("Create", mock.AnythingOfType("refresh_token.RefreshToken")).Return(nil)
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		access, newRefresh, err := svc.Refresh(ctx, oldRefresh)
//
//		assert.NoError(t, err)
//		assert.NotEmpty(t, access)
//		assert.NotEqual(t, oldRefresh, newRefresh)
//	})
//
//	t.Run("nieznany refresh – błąd", func(t *testing.T) {
//		uRepo := new(userRepoMock)
//		rtSvc := new(refreshServiceMock)
//
//		rtSvc.On("FindByToken", oldRefresh).Return(nil, errors.New("not found"))
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		_, _, err := svc.Refresh(ctx, oldRefresh)
//		assert.Error(t, err)
//	})
//}
//
//func TestService_Logout(t *testing.T) {
//	ctx := context.Background()
//	token := "rt"
//	userID := uint(1)
//
//	t.Run("logout z wszystkich urządzeń", func(t *testing.T) {
//		uRepo := new(userRepoMock)
//		rtSvc := new(refreshServiceMock)
//		rtSvc.On("DeleteByUserID", userID).Return(nil)
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		err := svc.Logout(ctx, userID, "", true)
//		assert.NoError(t, err)
//	})
//
//	t.Run("logout z konkretnego urządzenia", func(t *testing.T) {
//		uRepo := new(userRepoMock)
//		rtSvc := new(refreshServiceMock)
//
//		rt := &refresh_token.RefreshToken{ID: 42, Token: token}
//		rtSvc.On("FindByToken", token).Return(rt, nil)
//		rtSvc.On("Delete", rt.ID).Return(nil)
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		err := svc.Logout(ctx, userID, token, false)
//		assert.NoError(t, err)
//	})
//
//	t.Run("brak tokena przy pojedynczym urządzeniu – błąd", func(t *testing.T) {
//		uRepo := new(userRepoMock)
//		rtSvc := new(refreshServiceMock)
//
//		svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//
//		err := svc.Logout(ctx, userID, "", false)
//		assert.Error(t, err)
//	})
//}
//
///*
//	---  DODATKOWO: weryfikacja podpisów JWT (opcjonalnie) ---------------
//   Przykład jak można sprawdzić, że wygenerowany access ma poprawne claimy.
//*/
//
//func TestJWTClaimsAfterLogin(t *testing.T) {
//	ctx := context.Background()
//	uRepo := new(userRepoMock)
//	rtSvc := new(refreshServiceMock)
//
//	inp := LoginInput{Email: "marta@example.com", Password: "sekret"}
//	hashed := hashPass(t, inp.Password)
//	existing := &user.User{ID: 1, Email: inp.Email, Password: hashed}
//
//	uRepo.On("GetByEmail", inp.Email).Return(existing, nil)
//	rtSvc.On("Create", mock.AnythingOfType("refresh_token.RefreshToken")).Return(nil)
//
//	svc := &service{repo: uRepo, refreshTokenService: rtSvc, jwtKey: jwtKey}
//	access, _, err := svc.Login(ctx, inp)
//	assert.NoError(t, err)
//
//	claims, err := jwt.VerifyToken(access, jwtKey)
//	assert.NoError(t, err)
//	assert.Equal(t, inp.Email, claims.Subject)
//	assert.Equal(t, float64(existing.ID), claims["uid"])
//}
