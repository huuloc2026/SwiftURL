package jwtService

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SimpleJWTService struct {
	secret string
}

func NewSimpleJWTService(secret string) *SimpleJWTService {
	return &SimpleJWTService{secret: secret}
}

func (s *SimpleJWTService) GenerateToken(userID int64, username string, exp time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}
