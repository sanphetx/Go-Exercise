package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken() string
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService(secretKey string) Service {
	return &jwtService{
		secretKey: []byte(secretKey),
	}
}

func (s *jwtService) GenerateAccessToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"iat":  time.Now().Unix(),
		"type": "access",
	})
	return token.SignedString(s.secretKey)
}

func (s *jwtService) GenerateRefreshToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
}
