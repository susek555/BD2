package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JWTVerifier struct {
	secret []byte
}

func NewJWTVerifier(secret string) *JWTVerifier {
	return &JWTVerifier{secret: []byte(secret)}
}

func (j *JWTVerifier) VerifyToken(token string) (int64, error) {
	claims := &CustomClaims{}

	parsed, err := jwt.ParseWithClaims(
		token,
		claims,
		func(t *jwt.Token) (interface{}, error) { return j.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, fmt.Errorf("token expired: %w", err)
		}
		return 0, fmt.Errorf("token parse error: %w", err)
	}

	if !parsed.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}
