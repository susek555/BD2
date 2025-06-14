package auth_tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/auth"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/refresh_token"
	"github.com/susek555/BD2/car-dealer-api/internal/domains/user"
	"github.com/susek555/BD2/car-dealer-api/internal/models"
	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
	"github.com/susek555/BD2/car-dealer-api/pkg/middleware"
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(users []models.User, refreshTokens []models.RefreshToken) (user.UserRepositoryInterface, refresh_token.RefreshTokenServiceInterface, error) {
	dsn := "host=localhost user=bd2_user password=bd2_password dbname=bd2_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	db.Exec("TRUNCATE TABLE refresh_tokens, companies, people, users RESTART IDENTITY CASCADE")
	repo := user.NewUserRepository(db)
	for _, user_ := range users {
		err = repo.Create(&user_)
		if err != nil {
			return nil, nil, err
		}
	}
	refreshTokenRepo := refresh_token.NewRefreshTokenRepository(db)
	refreshTokenService := refresh_token.NewRefreshTokenService(refreshTokenRepo)
	for _, rt := range refreshTokens {
		err = refreshTokenService.Create(&rt)
		if err != nil {
			return nil, nil, err
		}
	}
	return repo, refreshTokenService, nil
}

func newTestServer(seedUsers []models.User, refreshTokens []models.RefreshToken) (*gin.Engine, *auth.Handler, refresh_token.RefreshTokenServiceInterface, error) {
	repo, rtSvc, err := setupDB(seedUsers, refreshTokens)
	if err != nil {
		return nil, nil, nil, err
	}

	svc := &auth.AuthService{Repo: repo, RefreshTokenService: rtSvc, JwtKey: []byte("secret")}
	h := &auth.Handler{Service: svc}
	verifier := jwt.NewJWTVerifier("secret")
	r := gin.Default()
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/logout", middleware.Authenticate(verifier), h.Logout)
	return r, h, rtSvc, nil
}

func getValidToken(userID uint, email string) (string, error) {
	secret := []byte("secret")
	return jwt.GenerateToken(email, int64(userID), secret, time.Now().Add(1*time.Hour))
}
func TestRegisterPersonSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusCreated
	assert.NoError(t, err)
	payload := `
	{
		"username": "unique_username",
		"email": "unique_email@example.com",
		"password": "PolskaGurom",
		"selector": "P",
		"person_name": "Herakles",
		"person_surname": "Wielki"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
}

func TestRegisterCompanySuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusCreated
	assert.NoError(t, err)
	payload := `
	{
		"username": "unique_username",
		"email": "unique_mail@example.com",
		"password": "PolskaGurom",
		"selector": "C",
		"company_name": "Herakles",
		"company_nip": "1234567890"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)

}

func TestRegisterPersonEmailAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "taken@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusConflict
	assert.NoError(t, err)
	payload := `
	{
		"username": "unique_username",
		"email": "taken@example.com",
		"password": "PolskaGurom",
		"selector": "P",
		"person_name": "Herakles",
		"person_surname": "Wielki"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	body := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, body)
	assert.Equal(t, "email already in use", response["errors"].(map[string]any)["email"].([]any)[0])
}

func TestRegisterPersonUsernameAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "unique@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusConflict
	assert.NoError(t, err)
	payload := `
	{
		"username": "taken_username",
		"email": "not_taken@example.com",
		"password": "PolskaGurom",
		"selector": "P",
		"person_name": "Herakles",
		"person_surname": "Wielki"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	body := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, body)
	assert.Equal(t, "username already in use", response["errors"].(map[string]any)["username"].([]any)[0])
}

func TestRegisterPersonInvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedTokens)
	wantStatus := http.StatusBadRequest
	assert.NoError(t, err)
	payload := `
	{
		"username": "unique_username",
		"email": "invalid_email",
		"password": "PolskaGurom",
		"selector": "P",
		"person_name": "Herakles",
		"person_surname": "Wielki"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	body := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, body)
	assert.Equal(t, "Key: 'CreateUserDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag", response["errors"].(map[string]any)["Email"].([]any)[0])
}

func TestRegisterCompanyEmailAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "taken@example.com",
			Username: "unique_username",
			Password: "PolskaGurom",
			Selector: "C",
			Company: &models.Company{
				Name: "Herakles",
				Nip:  "1234567890",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusConflict
	assert.NoError(t, err)
	payload := `
	{
		"username": "unique_username2",
		"email": "taken@example.com",
		"password": "PolskaGurom",
		"selector": "C",
		"company_name": "Herakles",
		"company_nip": "1234567890"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	body := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, body)
	assert.Equal(t, "email already in use", response["errors"].(map[string]any)["email"].([]any)[0])
}

func TestRegisterCompanyUsernameAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "unique@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "C",
			Company: &models.Company{
				Name: "Herakles",
				Nip:  "1234567890",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusConflict
	assert.NoError(t, err)
	payload := `
	{
		"username": "taken_username",
		"email": "unique2@example.com",
		"password": "PolskaGurom",
		"selector": "C",
		"company_name": "Herakles",
		"company_nip": "1234567890"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	body := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, body)
	assert.Equal(t, "username already in use", response["errors"].(map[string]any)["username"].([]any)[0])
}

func TestRegisterCompanyNipAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "unique@examle.com",
			Username: "unique_username",
			Password: "PolskaGurom",
			Selector: "C",
			Company: &models.Company{
				Name: "Herakles",
				Nip:  "1234567890",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusConflict
	assert.NoError(t, err)
	payload := `
	{
		"username": "unique_username2",
		"email": "unique@example.com",
		"password": "PolskaGurom",
		"selector": "C",
		"company_name": "Herakles",
		"company_nip": "1234567890"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	body := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, body)
	assert.Equal(t, "NIP already taken", response["errors"].(map[string]any)["company_nip"].([]any)[0])
}

func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hashedPassword, err := passwords.Hash("PolskaGurom")
	assert.NoError(t, err)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusOK
	assert.NoError(t, err)
	payload := `
	{
		"login": "herkules@gmail.com",
		"password": "PolskaGurom"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "herkules", response["username"])
	assert.Equal(t, "P", response["selector"])
	assert.Equal(t, "herkules@gmail.com", response["email"])
	assert.Equal(t, "Herakles", response["person_name"])
	assert.Equal(t, "Wielki", response["person_surname"])
	assert.NotEmpty(t, response["access_token"])
	assert.NotEmpty(t, response["refresh_token"])
}

func TestLoginInvalidLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hashedPassword, err := passwords.Hash("PolskaGurom")
	assert.NoError(t, err)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusUnauthorized
	assert.NoError(t, err)
	payload := `
	{
		"login": "invalid@gmail.com",
		"password": "PolskaGurom"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid credentials", response["errors"].(map[string]any)["credentials"].([]any)[0])
}

func TestLoginInvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hashedPassword, err := passwords.Hash("PolskaGurom")
	assert.NoError(t, err)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusUnauthorized
	assert.NoError(t, err)
	payload := `
	{
		"login": "herkules@gmail.com",
		"password": "invalid_password"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid credentials", response["errors"].(map[string]any)["credentials"].([]any)[0])
}

func TestLoginInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusBadRequest
	assert.NoError(t, err)
	payload := `
	{
		"login": "invalid",
		"password": "invalid"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid body", response["errors"].(map[string]any)["server"].([]any)[0])
}

func TestRefreshSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hashedPassword, err := passwords.Hash("PolskaGurom")
	assert.NoError(t, err)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []models.RefreshToken{
		{
			UserID:     1,
			Token:      "valid_refresh_token",
			ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusOK
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "valid_refresh_token"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])
}

func TestRefreshInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusUnauthorized
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "invalid_refresh_token"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid refresh token", response["error_description"])
}

func TestRefreshExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Username: "test_user",
			Email:    "test@example.com",
			Password: "hashed_password",
			Selector: "P",
			Person: &models.Person{
				Name:    "Test",
				Surname: "User",
			},
		},
	}
	seedRefreshTokens := []models.RefreshToken{
		{
			UserID:     1,
			Token:      "expired_refresh_token",
			ExpiryDate: time.Now().Add(-30 * 24 * time.Hour),
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusUnauthorized
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "expired_refresh_token"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "refresh token expired", response["error_description"])
}

func TestRefreshInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusBadRequest
	assert.NoError(t, err)
	payload := `
	{
		"invalid_field": "invalid"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid body", response["error_description"])
}

func TestLogoutSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []models.RefreshToken{
		{
			UserID:     1,
			Token:      "valid_refresh_token",
			ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
		},
	}
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusNoContent
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "valid_refresh_token"
	}
	`
	accessToken, err := getValidToken(1, "herkules@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
}

func TestLogoutNoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusUnauthorized
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "invalid_refresh_token"
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestLogoutNonExistingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusNotFound
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "invalid_refresh_token"
	}
	`
	accessToken, err := getValidToken(1, "herakles@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "refresh token not found", response["error_description"])
}

func TestLogoutEmptyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusBadRequest
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": ""
	}
	`
	accessToken, err := getValidToken(1, "herakles@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "refresh token required", response["error_description"])
}

func TestLogoutAllDevicesSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
		{
			Email:    "herakles2@gmail.com",
			Username: "herkules2",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []models.RefreshToken{
		{
			UserID:     1,
			Token:      "valid_refresh_token",
			ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
		},
		{
			UserID:     1,
			Token:      "valid_refresh_token_2",
			ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
		},
		{
			UserID:     2,
			Token:      "valid_refresh_token_3",
			ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
		},
	}
	server, _, rtSvc, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusNoContent
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "valid_refresh_token",
		"all_devices": true
	}
	`
	accessToken, err := getValidToken(1, "herakles@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	user1Tokens, err := rtSvc.FindByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, user1Tokens, 0)
	users2Tokens, err := rtSvc.FindByUserID(2)
	assert.NoError(t, err)
	assert.Len(t, users2Tokens, 1)
	assert.Equal(t, "valid_refresh_token_3", users2Tokens[0].Token)
	assert.Equal(t, uint(2), users2Tokens[0].UserID)
}

func TestLogoutAllDevicesNoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seedUsers []models.User
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusUnauthorized
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "invalid_refresh_token",
		"all_devices": true
	}
	`
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", response["message"])
}

func TestLogoutAllDevicesNonExistingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []models.RefreshToken{
		{
			UserID:     1,
			Token:      "valid_refresh_token",
			ExpiryDate: time.Now().UTC().Add(30 * 24 * time.Hour),
		},
	}
	server, _, rtSvc, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusNotFound
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "invalid_refresh_token",
		"all_devices": true
	}
	`
	accessToken, err := getValidToken(1, "herakles@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "refresh token not found", response["error_description"])
	user1Tokens, err := rtSvc.FindByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, user1Tokens, 1)
	assert.Equal(t, "valid_refresh_token", user1Tokens[0].Token)
	assert.Equal(t, uint(1), user1Tokens[0].UserID)
	warsawLocation, _ := time.LoadLocation("Europe/Warsaw")
	assert.Equal(t, time.Now().In(warsawLocation).Add(30*24*time.Hour).Format(time.RFC3339),
		user1Tokens[0].ExpiryDate.In(warsawLocation).Format(time.RFC3339))
}

func TestLogoutAllDevicesEmptyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []models.RefreshToken{
		{
			UserID:     1,
			Token:      "valid_refresh_token",
			ExpiryDate: time.Now().UTC().Add(30 * 24 * time.Hour),
		},
	}
	server, _, rtSvc, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusBadRequest
	assert.NoError(t, err)
	payload := `
	{
		"refresh_token": "",
		"all_devices": true
	}
	`
	accessToken, err := getValidToken(1, "herkules@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "refresh token required", response["error_description"])
	user1Tokens, err := rtSvc.FindByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, user1Tokens, 1)
	assert.Equal(t, "valid_refresh_token", user1Tokens[0].Token)
	assert.Equal(t, uint(1), user1Tokens[0].UserID)
	assert.Equal(t, time.Now().UTC().Add(30*24*time.Hour).Format(time.RFC3339), user1Tokens[0].ExpiryDate.UTC().Format(time.RFC3339))
}

func TestLogoutAllDevicesInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []models.User{
		{
			Email:    "herakles@gmail.com",
			Username: "herkules",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &models.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	var seedRefreshTokens []models.RefreshToken
	server, _, _, err := newTestServer(seedUsers, seedRefreshTokens)
	wantStatus := http.StatusBadRequest
	assert.NoError(t, err)
	payload := `
	{
		"invalid_field": "invalid",
		"all_devices": true
	}
	`
	accessToken, err := getValidToken(1, "herakles@gmail.com")
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	assert.Equal(t, wantStatus, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "refresh token required", response["error_description"])
}
