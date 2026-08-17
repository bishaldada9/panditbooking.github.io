package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/bees/hindu-ritual-platform/internal/auth/dto"
	"github.com/bees/hindu-ritual-platform/internal/auth/service"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type AuthHandler struct {
	service    *service.AuthService
	mfaService *service.MFAService
	log        *logger.Logger
}

func NewAuthHandler(authService *service.AuthService, mfaService *service.MFAService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		service:    authService,
		mfaService: mfaService,
		log:        log,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.service.Register(&req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.Error(c, http.StatusConflict, "Registration failed", err.Error())
		return
	}

	c.Set("user_id", result.User.ID)
	response.Created(c, "Registration successful", result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.service.Login(&req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		if err.Error() == "mfa_required" {
			response.Error(c, http.StatusAccepted, "MFA required", result)
			return
		}
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	c.Set("user_id", result.User.ID)
	response.Success(c, "Login successful", result)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	result, err := h.service.RefreshToken(req.RefreshToken, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	response.Success(c, "Token refreshed", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken := c.GetString("access_token")
	refreshToken := ""
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err == nil {
		refreshToken = req.RefreshToken
	}

	userID, _ := c.Get("user_id")
	if err := h.service.Logout(accessToken, refreshToken, userID.(uuid.UUID)); err != nil {
		h.log.Error("Logout failed", zap.Error(err))
	}

	response.Success(c, "Logged out successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.service.ForgotPassword(&req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process request", err.Error())
		return
	}

	response.Success(c, "If the email exists, a reset code has been sent", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Password reset failed", err.Error())
		return
	}

	response.Success(c, "Password has been reset successfully", nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.service.ChangePassword(userID.(uuid.UUID), &req); err != nil {
		response.Error(c, http.StatusBadRequest, "Password change failed", err.Error())
		return
	}

	response.Success(c, "Password changed successfully", nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.service.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	response.Success(c, "Profile retrieved", user)
}

func (h *AuthHandler) SetupMFA(c *gin.Context) {
	userID, _ := c.Get("user_id")
	result, err := h.mfaService.SetupMFA(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "MFA setup failed", err.Error())
		return
	}

	response.Success(c, "MFA setup initiated", result)
}

func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.MFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.mfaService.VerifyMFA(userID.(uuid.UUID), req.Code); err != nil {
		response.Error(c, http.StatusBadRequest, "MFA verification failed", err.Error())
		return
	}

	response.Success(c, "MFA enabled successfully", nil)
}

func (h *AuthHandler) DisableMFA(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.mfaService.DisableMFA(userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to disable MFA", err.Error())
		return
	}

	response.Success(c, "MFA disabled successfully", nil)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}

	if err := h.service.VerifyEmail(userID.(uuid.UUID), req.Code); err != nil {
		response.Error(c, http.StatusBadRequest, "Email verification failed", err.Error())
		return
	}

	response.Success(c, "Email verified successfully", nil)
}

func (h *AuthHandler) GetSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessions, err := h.service.GetUserSessions(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get sessions", err.Error())
		return
	}

	response.Success(c, "Sessions retrieved", sessions)
}

func (h *AuthHandler) LogoutAllSessions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.service.LogoutAllSessions(userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to logout all sessions", err.Error())
		return
	}

	response.Success(c, "All sessions logged out", nil)
}
