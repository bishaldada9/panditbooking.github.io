package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Code    int         `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type PaginatedResponse struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Meta     PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, details ...interface{}) {
	errResp := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	}
	if len(details) > 0 {
		errResp.Details = details[0]
	}
	c.JSON(statusCode, errResp)
}

func ValidationError(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusUnprocessableEntity,
		Details: details,
	})
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusUnauthorized,
	})
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusForbidden,
	})
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusNotFound,
	})
}

func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusInternalServerError,
	})
}

func Paginated(c *gin.Context, message string, data interface{}, total int64, page, limit int) {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			PerPage:    limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
