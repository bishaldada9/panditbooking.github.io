package dto

import "github.com/google/uuid"

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone" validate:"required,nepal_phone"`
	Bio      string `json:"bio"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	FullName        string    `json:"full_name"`
	Phone           string    `json:"phone"`
	Role            string    `json:"role"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsActive        bool      `json:"is_active"`
	IsSuspended     bool      `json:"is_suspended"`
	CreatedAt       string    `json:"created_at"`
}
