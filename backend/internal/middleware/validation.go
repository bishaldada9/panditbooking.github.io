package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/bees/hindu-ritual-platform/pkg/response"
	v "github.com/bees/hindu-ritual-platform/pkg/validator"
)

var validate *validator.Validate

func init() {
	validate = v.InitializeValidator()
}

func ValidateBody[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			handleValidationError(c, err)
			c.Abort()
			return
		}

		if err := validate.Struct(req); err != nil {
			errors := v.FormatValidationErrors(err)
			response.ValidationError(c, "Validation failed", errors)
			c.Abort()
			return
		}

		c.Set("validated_body", req)
		c.Next()
	}
}

func ValidateQuery[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindQuery(&req); err != nil {
			handleValidationError(c, err)
			c.Abort()
			return
		}

		if err := validate.Struct(req); err != nil {
			errors := v.FormatValidationErrors(err)
			response.ValidationError(c, "Validation failed", errors)
			c.Abort()
			return
		}

		c.Set("validated_query", req)
		c.Next()
	}
}

func ValidateURI[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindUri(&req); err != nil {
			handleValidationError(c, err)
			c.Abort()
			return
		}

		if err := validate.Struct(req); err != nil {
			errors := v.FormatValidationErrors(err)
			response.ValidationError(c, "Validation failed", errors)
			c.Abort()
			return
		}

		c.Set("validated_uri", req)
		c.Next()
	}
}

func ValidateForm[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindWith(&req, binding.Form); err != nil {
			handleValidationError(c, err)
			c.Abort()
			return
		}

		if err := validate.Struct(req); err != nil {
			errors := v.FormatValidationErrors(err)
			response.ValidationError(c, "Validation failed", errors)
			c.Abort()
			return
		}

		c.Set("validated_form", req)
		c.Next()
	}
}

func ValidateJSON[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			handleValidationError(c, err)
			c.Abort()
			return
		}

		if err := validate.Struct(req); err != nil {
			errors := v.FormatValidationErrors(err)
			response.ValidationError(c, "Validation failed", errors)
			c.Abort()
			return
		}

		c.Set("validated_json", req)
		c.Next()
	}
}

func validateContentType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodDelete || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")

		if contentType == "" {
			if len(allowedTypes) > 0 {
				response.Error(c, http.StatusUnsupportedMediaType, "Content-Type header is required")
				c.Abort()
				return
			}
			c.Next()
			return
		}

		contentType = strings.Split(contentType, ";")[0]
		contentType = strings.TrimSpace(contentType)

		allowed := false
		for _, t := range allowedTypes {
			if contentType == t {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, http.StatusUnsupportedMediaType,
				fmt.Sprintf("Unsupported Content-Type: %s. Allowed types: %s", contentType, strings.Join(allowedTypes, ", ")))
			c.Abort()
			return
		}

		c.Next()
	}
}

func ContentTypeValidationMiddleware() gin.HandlerFunc {
	return validateContentType("application/json")
}

func FormContentTypeValidationMiddleware() gin.HandlerFunc {
	return validateContentType("application/x-www-form-urlencoded", "multipart/form-data")
}

func handleValidationError(c *gin.Context, err error) {
	var errorMessages []v.ValidationError

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errorMessages = append(errorMessages, v.ValidationError{
				Field:   strings.ToLower(e.Field()),
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
				Message: getCustomErrorMessage(e),
			})
		}
		response.ValidationError(c, "Validation failed", errorMessages)
		return
	}

	response.ValidationError(c, "Invalid request format", nil)
}

func getCustomErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	case "nepali_phone":
		return "Invalid Nepali phone number format"
	case "ritual_date":
		return "Ritual date must be in the future"
	case "ritual_time":
		return "Invalid time format"
	case "amount":
		return "Invalid amount"
	case "nepali_name":
		return "Name contains invalid characters"
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
