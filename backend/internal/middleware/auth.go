package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authService "github.com/bees/hindu-ritual-platform/internal/auth/service"
	"github.com/bees/hindu-ritual-platform/pkg/redis"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type Claims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	DeviceID  string   `json:"device_id,omitempty"`
	TokenID   string   `json:"token_id"`
	Roles     []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type AuthMiddleware struct {
	jwtService  *authService.JWTService
	redisClient *redis.RedisClient
}

func NewAuthMiddleware(jwtService *authService.JWTService, redisClient *redis.RedisClient) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		redisClient: redisClient,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			response.Unauthorized(c, "Authentication token is required")
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		userID, _ := uuid.Parse(claims.UserID.String())
		c.Set("user_id", userID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("access_token", tokenString)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		if userRole.(string) != role {
			response.Forbidden(c, fmt.Sprintf("Role '%s' is required", role))
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
	}

	token, err := c.Cookie("access_token")
	if err == nil && token != "" {
		return token, nil
	}

	return "", fmt.Errorf("no authorization token found")
}
