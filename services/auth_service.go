package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(username, password string) (string, error)
}

type authService struct {
	secret string
}

func NewAuthService(secret string) AuthService {
	return &authService{secret: secret}
}

func (a *authService) Login(username, password string) (string, error) {

	fmt.Println("JWT SECRET:", a.secret)

	if username != "admin" || password != "password" {
		return "", errors.New("invalid credentials")
	}

	// claims := jwt.MapClaims{
	// 	"user_id": "demo-user-id",
	// 	"exp":     time.Now().Add(time.Hour * 1).Unix(),
	// }

	claims := jwt.MapClaims{
		"user_id": uuid.NewString(), // ✅ valid UUID
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
