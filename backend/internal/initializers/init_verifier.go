package initializers

import (
	"log"
	"os"

	"github.com/susek555/BD2/car-dealer-api/pkg/jwt"
)

var Verifier *jwt.JWTVerifier

func InitializeVerifier() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	Verifier = jwt.NewJWTVerifier(secret)
}
