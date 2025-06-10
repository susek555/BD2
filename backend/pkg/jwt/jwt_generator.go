package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(email string, userID int64, secret []byte, expiryTime time.Time) (string, error) {
	claims := CustomClaims{
		Email:  email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
			Issuer:    "car-dealer-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
