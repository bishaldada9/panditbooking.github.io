package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RitualCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:255" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"size:255" json:"icon"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Ritual struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CategoryID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Name             string         `gorm:"not null;size:255" json:"name"`
	Slug             string         `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Description      string         `gorm:"type:text" json:"description"`
	Duration         string         `gorm:"size:50" json:"duration"`
	BasePrice        float64        `gorm:"not null;default:0" json:"base_price"`
	RequiredItems    string         `gorm:"type:text" json:"required_items"`
	Procedure        string         `gorm:"type:text" json:"procedure"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	PanditCommission float64        `gorm:"default:0" json:"pandit_commission"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Category RitualCategory `gorm:"foreignKey:CategoryID" json:"-"`
}

func (r *Ritual) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
