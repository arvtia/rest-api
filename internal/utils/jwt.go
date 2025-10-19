package utils

import (
	"errors"
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID uint   `json:"admin_id"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

func getSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	return secret, nil
}

func GenerateJWT(adminID uint, email string) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	claims := Claims{
		AdminID: adminID,
		Email:   email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(adminID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ParseJWT(tokenStr string) (*Claims, error) {
	secret, err := getSecret()
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	parser := jwt.NewParser()
	_, err = parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
