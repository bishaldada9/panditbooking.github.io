package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
)

type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationApproved VerificationStatus = "approved"
	VerificationRejected VerificationStatus = "rejected"
)

type Pandit struct {
	ID                 uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID          `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	User               authModel.User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Bio                string             `gorm:"type:text" json:"bio"`
	ExperienceYears    int                `gorm:"default:0" json:"experience_years"`
	Specialization     string             `gorm:"size:500" json:"specialization"`
	Languages          datatypes.JSONSlice[string] `gorm:"type:jsonb" json:"languages"`
	BasePrice          float64            `gorm:"not null;default:0" json:"base_price"`
	Rating             float64            `gorm:"default:0" json:"rating"`
	TotalReviews       int                `gorm:"default:0" json:"total_reviews"`
	TotalBookings      int                `gorm:"default:0" json:"total_bookings"`
	IsAvailable        bool               `gorm:"default:true" json:"is_available"`
	VerificationStatus VerificationStatus `gorm:"default:pending" json:"verification_status"`
	VerificationNotes  string             `gorm:"type:text" json:"verification_notes"`
	VerifiedBy         *uuid.UUID          `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt         *time.Time          `json:"verified_at"`
	Latitude           float64            `json:"latitude"`
	Longitude          float64            `json:"longitude"`
	ServiceArea        string             `gorm:"size:255" json:"service_area"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`
}

type PanditDocument struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PanditID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"pandit_id"`
	DocumentType    string     `gorm:"not null;size:50" json:"document_type"`
	DocumentURL     string     `gorm:"not null;size:500" json:"document_url"`
	IsVerified      bool       `gorm:"default:false" json:"is_verified"`
	VerifiedBy      *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt      *time.Time `json:"verified_at"`
	RejectionReason string     `gorm:"type:text" json:"rejection_reason"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Availability struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PanditID  uuid.UUID `gorm:"type:uuid;not null;index:idx_pandit_date" json:"pandit_id"`
	Date      string    `gorm:"not null;size:10;index:idx_pandit_date" json:"date"`
	StartTime string    `gorm:"not null;size:5" json:"start_time"`
	EndTime   string    `gorm:"not null;size:5" json:"end_time"`
	IsBooked  bool      `gorm:"default:false" json:"is_booked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Pandit) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
