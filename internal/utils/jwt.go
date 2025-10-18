package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtKey = []byte("your_secret_key") // need to mmove to env
// utils/jwt.go
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func JwtKey() []byte {
	return jwtKey
}

type Claims struct {
	AdminID uint
	Email   string
	jwt.RegisteredClaims
}

func GenerateJWT(adminID uint, email string) (string, error) {
	claims := &Claims{
		AdminID: adminID,
		Email:   email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtKey)
}
