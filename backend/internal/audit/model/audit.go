package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
)

type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User       *authModel.User `gorm:"foreignKey:UserID" json:"-"`
	Action     string     `gorm:"not null;size:100;index" json:"action"`
	Resource   string     `gorm:"size:100" json:"resource"`
	ResourceID string     `gorm:"size:100" json:"resource_id"`
	Detail     string     `gorm:"type:text" json:"detail"`
	OldValue   datatypes.JSON `gorm:"type:jsonb" json:"-"`
	NewValue   datatypes.JSON `gorm:"type:jsonb" json:"-"`
	IP         string     `gorm:"size:45" json:"ip"`
	UserAgent  string     `gorm:"size:500" json:"user_agent"`
	Status     string     `gorm:"size:20" json:"status"`
	CreatedAt  time.Time  `gorm:"index" json:"created_at"`
}
