package dto

import "github.com/google/uuid"

type NotificationResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Message       string    `json:"message"`
	IsRead        bool      `json:"is_read"`
	ReferenceID   string    `json:"reference_id,omitempty"`
	ReferenceType string    `json:"reference_type,omitempty"`
	CreatedAt     string    `json:"created_at"`
}
