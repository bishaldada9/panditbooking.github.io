package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_strength"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone" validate:"required,nepal_phone"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	DeviceID string `json:"device_id"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         UserDTO `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type MFASetupResponse struct {
	Secret     string   `json:"secret"`
	QRCodeURL  string   `json:"qr_code_url"`
	RecoveryCodes []string `json:"recovery_codes"`
}

type MFAVerifyRequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

type MFALoginRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	MFACode      string `json:"mfa_code"`
	RecoveryCode string `json:"recovery_code"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,password_strength"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,password_strength"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type UserDTO struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	FullName        string    `json:"full_name"`
	Phone           string    `json:"phone"`
	Role            string    `json:"role"`
	IsEmailVerified bool      `json:"is_email_verified"`
	MFAEnabled      bool      `json:"mfa_enabled"`
	IsActive        bool      `json:"is_active"`
	IsSuspended     bool      `json:"is_suspended"`
}
