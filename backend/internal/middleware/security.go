package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bees/hindu-ritual-platform/pkg/configs"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type SecurityMiddleware struct {
	config *configs.Config
}

func NewSecurityMiddleware(cfg *configs.Config) *SecurityMiddleware {
	return &SecurityMiddleware{config: cfg}
}

func (m *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

func (m *SecurityMiddleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" && m.isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID, X-Device-ID, Accept, Origin")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Request-ID")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (m *SecurityMiddleware) isAllowedOrigin(origin string) bool {
	for _, allowed := range m.config.CORS.AllowedOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

type RateLimiter struct {
	mu                sync.RWMutex
	visitors          map[string]*tokenBucket
	requestsPerMinute int
	burstSize         int
}

type tokenBucket struct {
	tokens    float64
	lastCheck time.Time
	rate      float64
	burst     int
}

func NewRateLimiter(requestsPerMinute, burstSize int) *RateLimiter {
	rl := &RateLimiter{
		visitors:          make(map[string]*tokenBucket),
		requestsPerMinute: requestsPerMinute,
		burstSize:         burstSize,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for key, bucket := range rl.visitors {
			if time.Since(bucket.lastCheck) > 5*time.Minute {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.visitors[key]
	now := time.Now()

	if !exists {
		rl.visitors[key] = &tokenBucket{
			tokens:    float64(rl.burstSize) - 1,
			lastCheck: now,
			rate:      float64(rl.requestsPerMinute) / 60.0,
			burst:     rl.burstSize,
		}
		return true
	}

	elapsed := now.Sub(bucket.lastCheck).Seconds()
	bucket.tokens += elapsed * bucket.rate
	if bucket.tokens > float64(bucket.burst) {
		bucket.tokens = float64(bucket.burst)
	}
	bucket.lastCheck = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			key = fmt.Sprintf("user:%v", userID)
		}

		if !rl.allow(key) {
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(\bSELECT\b.*\bFROM\b|\bINSERT\b.*\bINTO\b|\bUPDATE\b.*\bSET\b|\bDELETE\b.*\bFROM\b|\bDROP\b.*\bTABLE\b|\bUNION\b.*\bSELECT\b)`),
	regexp.MustCompile(`(?i)(\bOR\b\s+1\s*=\s*1|\bOR\b\s+true\b)`),
	regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onload=)`),
}

func SQLInjectionDetectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, pattern := range dangerousPatterns {
			if pattern.MatchString(c.Request.RequestURI) {
				response.Error(c, http.StatusBadRequest, "Invalid request parameters", nil)
				c.Abort()
				return
			}
		}

		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				for _, pattern := range dangerousPatterns {
					if pattern.MatchString(value) {
						response.Error(c, http.StatusBadRequest, "Invalid input detected", nil)
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

func sanitizeInput(input string) string {
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(input)
}
