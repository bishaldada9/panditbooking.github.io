package dto

import "github.com/google/uuid"

type CreateReviewRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Comment   string    `json:"comment" validate:"required"`
}

type ReviewResponse struct {
	ID         uuid.UUID `json:"id"`
	BookingID  uuid.UUID `json:"booking_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	PanditID   uuid.UUID `json:"pandit_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	IsVerified bool      `json:"is_verified"`
	IsVisible  bool      `json:"is_visible"`
	AdminReply string    `json:"admin_reply,omitempty"`
	CreatedAt  string    `json:"created_at"`
}
