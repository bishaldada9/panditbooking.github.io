package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	panditModel "github.com/bees/hindu-ritual-platform/internal/pandits/model"
	ritualModel "github.com/bees/hindu-ritual-platform/internal/rituals/model"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCompleted BookingStatus = "completed"
	BookingCancelled BookingStatus = "cancelled"
	BookingRejected  BookingStatus = "rejected"
)

type Booking struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CustomerID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"customer_id"`
	PanditID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"pandit_id"`
	RitualID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"ritual_id"`
	Customer      authModel.User     `gorm:"foreignKey:CustomerID" json:"-"`
	Pandit        panditModel.Pandit   `gorm:"foreignKey:PanditID" json:"-"`
	Ritual        ritualModel.Ritual   `gorm:"foreignKey:RitualID" json:"-"`
	Status        BookingStatus  `gorm:"not null;default:pending;index" json:"status"`
	ScheduledDate string         `gorm:"not null;size:10" json:"scheduled_date"`
	StartTime     string         `gorm:"not null;size:5" json:"start_time"`
	EndTime       string         `gorm:"not null;size:5" json:"end_time"`
	Address       string         `gorm:"not null;type:text" json:"address"`
	Latitude      float64        `json:"latitude"`
	Longitude     float64        `json:"longitude"`
	SpecialNotes  string         `gorm:"type:text" json:"special_notes"`
	TotalAmount   float64        `gorm:"not null" json:"total_amount"`
	PlatformFee   float64        `gorm:"default:0" json:"platform_fee"`
	PanditAmount  float64        `gorm:"default:0" json:"pandit_amount"`
	CancelledBy   *uuid.UUID    `gorm:"type:uuid" json:"cancelled_by"`
	CancelReason  string         `gorm:"type:text" json:"cancel_reason"`
	CancelledAt   *time.Time     `json:"cancelled_at"`
	ConfirmedAt   *time.Time     `json:"confirmed_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Complaint struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"booking_id"`
	ComplainantID uuid.UUID  `gorm:"type:uuid;not null" json:"complainant_id"`
	Booking       Booking       `gorm:"foreignKey:BookingID" json:"-"`
	Complainant   authModel.User `gorm:"foreignKey:ComplainantID" json:"-"`
	Subject       string     `gorm:"not null;size:255" json:"subject"`
	Description   string     `gorm:"type:text" json:"description"`
	Status        string     `gorm:"default:pending;size:20" json:"status"`
	AdminNotes    string     `gorm:"type:text" json:"admin_notes"`
	ResolvedBy    *uuid.UUID `gorm:"type:uuid" json:"resolved_by"`
	ResolvedAt    *time.Time `json:"resolved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
