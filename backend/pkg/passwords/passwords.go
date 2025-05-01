package passwords

import "golang.org/x/crypto/bcrypt"

const cost = bcrypt.DefaultCost

func Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	return string(b), err
}

func Match(raw, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}
