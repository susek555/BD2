package auth

import (
	"errors"
	"net/http"
)

var (
	ErrEmailTaken           = errors.New("email already in use")
	ErrUsernameTaken        = errors.New("username already in use")
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrInvalidBody          = errors.New("invalid body")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrNipAlreadyTaken      = errors.New("NIP already taken")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrCreateRefreshToken   = errors.New("error - create refresh token")
	ErrRefreshTokenRequired = errors.New("refresh token required")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrUnauthorized         = errors.New("unauthorized")
)

var ErrorMap = map[error]int{
	ErrEmailTaken:           http.StatusConflict,
	ErrUsernameTaken:        http.StatusConflict,
	ErrInvalidCredentials:   http.StatusUnauthorized,
	ErrInvalidBody:          http.StatusBadRequest,
	ErrInvalidRefreshToken:  http.StatusBadRequest,
	ErrNipAlreadyTaken:      http.StatusConflict,
	ErrRefreshTokenExpired:  http.StatusBadRequest,
	ErrCreateRefreshToken:   http.StatusBadRequest,
	ErrRefreshTokenRequired: http.StatusBadRequest,
	ErrRefreshTokenNotFound: http.StatusBadRequest,
	ErrUnauthorized:         http.StatusUnauthorized,
}
