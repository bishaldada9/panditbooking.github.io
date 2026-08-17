package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	bookingModel "github.com/bees/hindu-ritual-platform/internal/bookings/model"
	panditModel "github.com/bees/hindu-ritual-platform/internal/pandits/model"
)

type Review struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID  uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"booking_id"`
	CustomerID uuid.UUID      `gorm:"type:uuid;not null;index" json:"customer_id"`
	PanditID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"pandit_id"`
	Booking    bookingModel.Booking `gorm:"foreignKey:BookingID" json:"-"`
	Customer   authModel.User      `gorm:"foreignKey:CustomerID" json:"-"`
	Pandit     panditModel.Pandit  `gorm:"foreignKey:PanditID" json:"-"`
	Rating     int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment    string         `gorm:"type:text" json:"comment"`
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	IsVisible  bool           `gorm:"default:true" json:"is_visible"`
	AdminReply string         `gorm:"type:text" json:"admin_reply"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
