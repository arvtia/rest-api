package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AdminID uint   `json:"admin_id,omitempty"`
	UserID  uint   `json:"user_id,omitempty"`
	Email   string `json:"email"`
	Role    string `json:"role"` // "admin" or "user"
	jwt.RegisteredClaims
}

func getSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	return secret, nil
}

func GenerateJWT(id uint, email, role string) (string, error) {
	secret, err := getSecret()
	if err != nil {
		return "", err
	}

	claims := Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(id),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	if role == "admin" {
		claims.AdminID = id
	} else {
		claims.UserID = id
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
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
