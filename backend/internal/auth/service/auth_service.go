package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/bees/hindu-ritual-platform/internal/auth/dto"
	"github.com/bees/hindu-ritual-platform/internal/auth/model"
	authRepo "github.com/bees/hindu-ritual-platform/internal/auth/repository"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
	"github.com/bees/hindu-ritual-platform/pkg/security"
)

type AuthService struct {
	repo       *authRepo.AuthRepository
	jwtService *JWTService
	mfaService *MFAService
	log        *logger.Logger
}

func NewAuthService(repo *authRepo.AuthRepository, jwtService *JWTService, mfaService *MFAService, log *logger.Logger) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
		mfaService: mfaService,
		log:        log,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest, ip, userAgent string) (*dto.LoginResponse, error) {
	existingUser, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Role:         "customer",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate email verification OTP
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	otp := &model.OTP{
		UserID:    user.ID,
		Code:      code,
		Purpose:   "email_verification",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.repo.CreateOTP(otp); err != nil {
		s.log.Error("Failed to create verification OTP", zap.Error(err))
	}

	tokens, err := s.jwtService.GenerateTokenPair(user, "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save refresh token
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		IP:        ip,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(rt); err != nil {
		s.log.Error("Failed to save refresh token", zap.Error(err))
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    3600,
		User:         s.toUserDTO(user),
	}, nil
}

func (s *AuthService) Login(req *dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error) {
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if user.IsSuspended {
		return nil, errors.New("account is suspended")
	}

	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, fmt.Errorf("account is locked until %s", user.LockedUntil.Format(time.RFC3339))
	}

	if !security.CheckPassword(user.PasswordHash, req.Password) {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= 5 {
			lockDuration := 15 * time.Minute
			lockTime := time.Now().Add(lockDuration)
			user.LockedUntil = &lockTime
			s.log.Warn("Account locked due to failed attempts", zap.String("user_id", user.ID.String()), zap.String("email", user.Email))
		}
		if err := s.repo.UpdateUser(user); err != nil {
			s.log.Error("Failed to update user after failed login", zap.Error(err))
		}
		return nil, errors.New("invalid email or password")
	}

	// Reset failed attempts on successful login
	if err := s.repo.ResetFailedAttempts(user.ID); err != nil {
		s.log.Error("Failed to reset failed attempts", zap.Error(err))
	}

	// Update last login
	if err := s.repo.UpdateLastLogin(user.ID, ip, userAgent); err != nil {
		s.log.Error("Failed to update last login", zap.Error(err))
	}

	// Check MFA
	if user.MFAEnabled {
		return &dto.LoginResponse{
			User: s.toUserDTO(user),
		}, errors.New("mfa_required")
	}

	tokens, err := s.jwtService.GenerateTokenPair(user, req.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		DeviceID:  req.DeviceID,
		IP:        ip,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(rt); err != nil {
		s.log.Error("Failed to save refresh token", zap.Error(err))
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    3600,
		User:         s.toUserDTO(user),
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string, ip, userAgent string) (*dto.RefreshTokenResponse, error) {
	// Validate the refresh token
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Check if token is in Redis blacklist
	blacklisted, err := s.jwtService.IsTokenBlacklisted(refreshToken)
	if err != nil || blacklisted {
		return nil, errors.New("token has been revoked")
	}

	// Find in database
	storedToken, err := s.repo.FindRefreshToken(refreshToken)
	if err != nil || storedToken == nil {
		return nil, errors.New("refresh token not found or revoked")
	}

	// Revoke old token
	if err := s.repo.RevokeRefreshToken(storedToken.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	user, err := s.repo.FindUserByID(claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new token pair
	tokens, err := s.jwtService.GenerateTokenPair(user, storedToken.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Save new refresh token
	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     tokens.RefreshToken,
		DeviceID:  storedToken.DeviceID,
		IP:        ip,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.repo.CreateRefreshToken(rt); err != nil {
		s.log.Error("Failed to save new refresh token", zap.Error(err))
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (s *AuthService) Logout(accessToken, refreshToken string, userID uuid.UUID) error {
	// Blacklist tokens in Redis
	if err := s.jwtService.BlacklistToken(accessToken); err != nil {
		s.log.Error("Failed to blacklist access token", zap.Error(err))
	}

	if refreshToken != "" {
		if err := s.jwtService.BlacklistToken(refreshToken); err != nil {
			s.log.Error("Failed to blacklist refresh token", zap.Error(err))
		}

		// Revoke in database
		storedToken, err := s.repo.FindRefreshToken(refreshToken)
		if err == nil && storedToken != nil {
			if err := s.repo.RevokeRefreshToken(storedToken.ID); err != nil {
				s.log.Error("Failed to revoke refresh token", zap.Error(err))
			}
		}
	}

	return nil
}

func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil || user == nil {
		// Don't reveal if email exists
		return nil
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	otp := &model.OTP{
		UserID:    user.ID,
		Code:      code,
		Purpose:   "password_reset",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.repo.CreateOTP(otp); err != nil {
		return fmt.Errorf("failed to create reset code: %w", err)
	}

	s.log.Info("Password reset code generated", zap.String("user_id", user.ID.String()), zap.String("code", code))
	return nil
}

func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// TODO: Parse token to get user ID and validate
	// For now, find any valid OTP with the token as code
	var userID uuid.UUID
	// This would be decoded from the reset token

	otp, err := s.repo.FindValidOTP(userID, req.Token, "password_reset")
	if err != nil || otp == nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.UpdatePassword(otp.UserID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := s.repo.MarkOTPUsed(otp.ID); err != nil {
		s.log.Error("Failed to mark OTP as used", zap.Error(err))
	}

	// Revoke all existing sessions
	if err := s.repo.RevokeAllUserTokens(otp.UserID); err != nil {
		s.log.Error("Failed to revoke all user tokens", zap.Error(err))
	}

	return nil
}

func (s *AuthService) ChangePassword(userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.FindUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	if !security.CheckPassword(user.PasswordHash, req.OldPassword) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.UpdatePassword(userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) GetUserByID(id uuid.UUID) (*dto.UserDTO, error) {
	user, err := s.repo.FindUserByID(id)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	dto := s.toUserDTO(user)
	return &dto, nil
}

func (s *AuthService) VerifyEmail(userID uuid.UUID, code string) error {
	otp, err := s.repo.FindValidOTP(userID, code, "email_verification")
	if err != nil || otp == nil {
		return errors.New("invalid or expired verification code")
	}

	if err := s.repo.VerifyEmail(userID); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	if err := s.repo.MarkOTPUsed(otp.ID); err != nil {
		s.log.Error("Failed to mark OTP as used", zap.Error(err))
	}

	return nil
}

func (s *AuthService) GetUserSessions(userID uuid.UUID) ([]model.RefreshToken, error) {
	return s.repo.GetUserSessions(userID)
}

func (s *AuthService) LogoutAllSessions(userID uuid.UUID) error {
	return s.repo.RevokeAllUserTokens(userID)
}

func (s *AuthService) toUserDTO(user *model.User) dto.UserDTO {
	return dto.UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		FullName:        user.FullName,
		Phone:           user.Phone,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		MFAEnabled:      user.MFAEnabled,
		IsActive:        user.IsActive,
		IsSuspended:     user.IsSuspended,
	}
}
