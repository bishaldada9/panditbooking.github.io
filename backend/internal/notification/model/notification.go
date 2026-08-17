package model

import (
	"time"

	"github.com/google/uuid"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
)

type Notification struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User          authModel.User `gorm:"foreignKey:UserID" json:"-"`
	Type          string    `gorm:"not null;size:50;index" json:"type"`
	Title         string    `gorm:"not null;size:255" json:"title"`
	Message       string    `gorm:"type:text" json:"message"`
	IsRead        bool      `gorm:"default:false;index" json:"is_read"`
	ReferenceID   string    `gorm:"size:100" json:"reference_id"`
	ReferenceType string    `gorm:"size:50" json:"reference_type"`
	CreatedAt     time.Time `json:"created_at"`
}
