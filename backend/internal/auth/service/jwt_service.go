package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/bees/hindu-ritual-platform/internal/auth/model"
	"github.com/bees/hindu-ritual-platform/pkg/configs"
	"github.com/bees/hindu-ritual-platform/pkg/redis"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	DeviceID string    `json:"device_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	config      *configs.Config
	redisClient *redis.RedisClient
}

func NewJWTService(config *configs.Config, redisClient *redis.RedisClient) *JWTService {
	return &JWTService{
		config:      config,
		redisClient: redisClient,
	}
}

func (s *JWTService) GenerateTokenPair(user *model.User, deviceID string) (*TokenPair, error) {
	// Access token - short lived (1 hour)
	accessClaims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hindu-ritual-platform",
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.config.App.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshClaims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     user.Role,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hindu-ritual-platform",
			Subject:   user.ID.String(),
			ID:        uuid.New().String(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.config.App.JWTRefreshSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.App.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	blacklisted, err := s.redisClient.IsBlacklisted(tokenString)
	if err != nil || blacklisted {
		return nil, errors.New("token has been revoked")
	}

	return claims, nil
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.App.JWTRefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

func (s *JWTService) BlacklistToken(tokenString string) error {
	// Parse token to get expiry
	claims := &Claims{}
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok {
		ttl := time.Until(claims.ExpiresAt.Time)
		if ttl > 0 {
			return s.redisClient.AddToBlacklist(tokenString, ttl)
		}
	}

	return nil
}

func (s *JWTService) IsTokenBlacklisted(tokenString string) (bool, error) {
	return s.redisClient.IsBlacklisted(tokenString)
}
