package utils

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	Email  string `json:"email"`
	UserID int64  `json:"user_id"`
	jwt.RegisteredClaims
}
