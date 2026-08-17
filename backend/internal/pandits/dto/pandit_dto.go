package dto

import (
	"github.com/bees/hindu-ritual-platform/internal/pandits/model"
	"github.com/google/uuid"
)

type RegisterPanditRequest struct {
	Bio             string   `json:"bio" validate:"required"`
	ExperienceYears int      `json:"experience_years" validate:"required,min=0,max=70"`
	Specialization  string   `json:"specialization" validate:"required"`
	Languages       []string `json:"languages" validate:"required,min=1"`
	BasePrice       float64  `json:"base_price" validate:"required,min=0"`
	ServiceArea     string   `json:"service_area" validate:"required"`
}

type UpdatePanditProfileRequest struct {
	Bio             string   `json:"bio"`
	ExperienceYears int      `json:"experience_years"`
	Specialization  string   `json:"specialization"`
	Languages       []string `json:"languages"`
	BasePrice       float64  `json:"base_price"`
	ServiceArea     string   `json:"service_area"`
	IsAvailable     *bool    `json:"is_available"`
}

type PanditResponse struct {
	ID                 uuid.UUID                `json:"id"`
	UserID             uuid.UUID                `json:"user_id"`
	FullName           string                   `json:"full_name"`
	Email              string                   `json:"email"`
	Phone              string                   `json:"phone"`
	Bio                string                   `json:"bio"`
	ExperienceYears    int                      `json:"experience_years"`
	Specialization     string                   `json:"specialization"`
	Languages          []string                 `json:"languages"`
	BasePrice          float64                  `json:"base_price"`
	Rating             float64                  `json:"rating"`
	TotalReviews       int                      `json:"total_reviews"`
	TotalBookings      int                      `json:"total_bookings"`
	IsAvailable        bool                     `json:"is_available"`
	VerificationStatus model.VerificationStatus `json:"verification_status"`
	ServiceArea        string                   `json:"service_area"`
	CreatedAt          string                   `json:"created_at"`
}

type UpdateAvailabilityRequest struct {
	Date      string `json:"date" validate:"required"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
	IsBooked  bool   `json:"is_booked"`
}
