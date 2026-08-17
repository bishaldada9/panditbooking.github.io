package dto

import (
	"github.com/google/uuid"
	"github.com/bees/hindu-ritual-platform/internal/bookings/model"
)

type CreateBookingRequest struct {
	PanditID      uuid.UUID `json:"pandit_id" validate:"required"`
	RitualID      uuid.UUID `json:"ritual_id" validate:"required"`
	ScheduledDate string    `json:"scheduled_date" validate:"required"`
	StartTime     string    `json:"start_time" validate:"required"`
	EndTime       string    `json:"end_time" validate:"required"`
	Address       string    `json:"address" validate:"required"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	SpecialNotes  string    `json:"special_notes"`
}

type BookingResponse struct {
	ID            uuid.UUID           `json:"id"`
	CustomerID    uuid.UUID           `json:"customer_id"`
	PanditID      uuid.UUID           `json:"pandit_id"`
	RitualID      uuid.UUID           `json:"ritual_id"`
	Status        model.BookingStatus `json:"status"`
	ScheduledDate string              `json:"scheduled_date"`
	StartTime     string              `json:"start_time"`
	EndTime       string              `json:"end_time"`
	Address       string              `json:"address"`
	TotalAmount   float64             `json:"total_amount"`
	PlatformFee   float64             `json:"platform_fee"`
	SpecialNotes  string              `json:"special_notes"`
	CreatedAt     string              `json:"created_at"`
}

type CancelBookingRequest struct {
	Reason string `json:"reason" validate:"required"`
}
