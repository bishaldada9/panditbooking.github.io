package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	auditSvc "github.com/bees/hindu-ritual-platform/internal/audit/service"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func RequestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		}

		if requestID, exists := c.Get("request_id"); exists {
			fields = append(fields, zap.String("request_id", requestID.(string)))
		}

		if status >= 500 {
			log.Error("Request failed", fields...)
		} else if status >= 400 {
			log.Warn("Request warning", fields...)
		} else {
			log.Info("Request completed", fields...)
		}
	}
}

func PanicRecovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("request_id", requestID.(string)),
					zap.String("path", c.Request.URL.Path),
					zap.Stack("stack"),
				)
				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

func AuditTrail(auditService *auditSvc.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if auditService == nil || !isAuditableMethod(c.Request.Method) {
			return
		}

		var userID *uuid.UUID
		if value, exists := c.Get("user_id"); exists {
			if id, ok := value.(uuid.UUID); ok {
				userID = &id
			}
		}

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		resourceID := firstRouteID(c)
		status := "success"
		if c.Writer.Status() >= 400 {
			status = "failed"
		}

		detail := fmt.Sprintf("%s %s returned %d", c.Request.Method, route, c.Writer.Status())
		_ = auditService.Log(
			userID,
			strings.ToLower(c.Request.Method)+" "+route,
			route,
			resourceID,
			detail,
			"{}",
			"{}",
			c.ClientIP(),
			c.GetHeader("User-Agent"),
			status,
		)
	}
}

func isAuditableMethod(method string) bool {
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

func firstRouteID(c *gin.Context) string {
	for _, name := range []string{"id", "bookingId", "paymentId", "reviewId", "gateway"} {
		if value := c.Param(name); value != "" {
			return value
		}
	}
	return ""
}
