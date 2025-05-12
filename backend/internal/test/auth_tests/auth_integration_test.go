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
	"github.com/susek555/BD2/car-dealer-api/pkg/passwords"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ------
// Setup
// ------

func setupDB(users []user.User, refreshTokens []refresh_token.RefreshToken) (user.UserRepositoryInterface, refresh_token.RefreshTokenServiceInterface, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&user.User{},
		&user.Person{},
		&user.Company{},
		&refresh_token.RefreshToken{},
	)
	if err != nil {
		return nil, nil, err
	}
	repo := user.NewUserRepository(db)
	for _, user := range users {
		repo.Create(&user)
	}
	refreshTokenService := refresh_token.NewRefreshTokenService(db)
	for _, rt := range refreshTokens {
		refreshTokenService.Create(&rt)
	}
	return repo, refreshTokenService, nil
}

func newTestServer(seedUsers []user.User, refreshTokens []refresh_token.RefreshToken) (*gin.Engine, error) {
	repo, rtSvc, err := setupDB(seedUsers, refreshTokens)
	if err != nil {
		return nil, err
	}

	svc := &auth.AuthService{Repo: repo, RefreshTokenService: rtSvc, JwtKey: []byte("secret")}
	h := &auth.Handler{Service: svc}

	r := gin.Default()
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
	return r, nil
}

func TestRegisterPersonSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestRegisterCompanySuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestRegisterPersonEmailAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{
		{
			Email:    "taken@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	assert.Equal(t, "Email already taken", response["errors"].(map[string]any)["email"].([]any)[0])
}

func TestRegisterPersonUsernameAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{
		{
			Email:    "unique@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "P",
			Person: &user.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	assert.Equal(t, "Username already taken", response["errors"].(map[string]any)["username"].([]any)[0])
}

func TestRegisterPersonInvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{}
	seedTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedTokens)
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
	seedUsers := []user.User{
		{
			Email:    "taken@example.com",
			Username: "unique_username",
			Password: "PolskaGurom",
			Selector: "C",
			Company: &user.Company{
				Name: "Herakles",
				NIP:  "1234567890",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	assert.Equal(t, "Email already taken", response["errors"].(map[string]any)["email"].([]any)[0])
}

func TestRegisterCompanyUsernameAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{
		{
			Email:    "unique@example.com",
			Username: "taken_username",
			Password: "PolskaGurom",
			Selector: "C",
			Company: &user.Company{
				Name: "Herakles",
				NIP:  "1234567890",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	assert.Equal(t, "Username already taken", response["errors"].(map[string]any)["username"].([]any)[0])
}

func TestRegisterCompanyNipAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	seedUsers := []user.User{
		{
			Email:    "unique@examle.com",
			Username: "unique_username",
			Password: "PolskaGurom",
			Selector: "C",
			Company: &user.Company{
				Name: "Herakles",
				NIP:  "1234567890",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &user.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &user.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &user.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{
		{
			Email:    "herkules@gmail.com",
			Username: "herkules",
			Password: hashedPassword,
			Selector: "P",
			Person: &user.Person{
				Name:    "Herakles",
				Surname: "Wielki",
			},
		},
	}
	seedRefreshTokens := []refresh_token.RefreshToken{
		{
			UserId:     1,
			Token:      "valid_refresh_token",
			ExpiryDate: time.Now().Add(30 * 24 * time.Hour),
		},
	}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{}
	seedRefreshTokens := []refresh_token.RefreshToken{}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
	seedUsers := []user.User{}
	seedRefreshTokens := []refresh_token.RefreshToken{
		{
			UserId:     1,
			Token:      "expired_refresh_token",
			ExpiryDate: time.Now().Add(-30 * 24 * time.Hour),
		},
	}
	server, err := newTestServer(seedUsers, seedRefreshTokens)
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
