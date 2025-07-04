package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"

	"Go-Exercise/pkg/model"
	"Go-Exercise/pkg/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
	ErrUserNotFound = errors.New("user not found")
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthService struct {
	repo      repository.AuthRepository
	secretKey []byte
}

func NewAuthService(repo repository.AuthRepository, secretKey string) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: []byte(secretKey),
	}
}

func (s *AuthService) Register(name, email, password string, age int) (*model.User, error) {
	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		Age:      age,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(email, password, deviceInfo string) (*TokenResponse, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	ok, err := argon2.VerifyEncoded([]byte(password), []byte(user.Password))
	if err != nil || !ok {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	token := generateRefreshToken()
	refreshTokenExpiry := getRefreshTokenExpiry()
	refreshToken := &model.RefreshToken{
		ID:         uuid.New(),
		UserID:     user.ID,
		Token:      token,
		ExpiresAt:  time.Now().Add(time.Duration(refreshTokenExpiry) * time.Second).Unix(),
		DeviceInfo: deviceInfo,
	}

	if err := s.repo.SaveRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: token,
		TokenType:    "Bearer",
		ExpiresIn:    getAccessTokenExpiry(),
	}, nil
}

func (s *AuthService) RefreshToken(refreshTokenStr, deviceInfo string) (*TokenResponse, error) {
	token, err := s.repo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if token.IsRevoked || time.Now().After(time.Unix(token.ExpiresAt, 0)) {
		return nil, ErrExpiredToken
	}

	user, err := s.repo.FindUserByID(token.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	token.IsRevoked = true
	if err := s.repo.UpdateRefreshToken(token); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshTokenExpiry := getRefreshTokenExpiry()
	newRefreshToken := &model.RefreshToken{
		ID:         uuid.New(),
		UserID:     user.ID,
		Token:      generateRefreshToken(),
		ExpiresAt:  time.Now().Add(time.Duration(refreshTokenExpiry) * time.Second).Unix(),
		DeviceInfo: deviceInfo,
	}

	if err := s.repo.SaveRefreshToken(newRefreshToken); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
		TokenType:    "Bearer",
		ExpiresIn:    getAccessTokenExpiry(),
	}, nil
}

func (s *AuthService) Logout(refreshTokenStr string, logoutAll bool) error {
	token, err := s.repo.GetRefreshToken(refreshTokenStr)
	if err != nil {
		return ErrInvalidToken
	}

	if logoutAll {
		return s.repo.RevokeAllUserTokens(token.UserID)
	}

	token.IsRevoked = true
	return s.repo.UpdateRefreshToken(token)
}

func (s *AuthService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(getAccessTokenExpiry()) * time.Second).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func generateRefreshToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
}

func getAccessTokenExpiry() int64 {
	expiry := os.Getenv("ACCESS_TOKEN_EXPIRY")
	if expiry == "" {
		return 3600
	}
	val, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		return 3600
	}
	return val
}

func getRefreshTokenExpiry() int64 {
	expiry := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if expiry == "" {
		return 2592000
	}
	val, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		return 2592000
	}
	return val
}
