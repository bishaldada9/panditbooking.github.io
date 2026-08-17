package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email              string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Phone              string         `gorm:"uniqueIndex;size:20" json:"phone"`
	PasswordHash       string         `gorm:"not null;size:255" json:"-"`
	FullName           string         `gorm:"not null;size:255" json:"full_name"`
	Role               string         `gorm:"not null;size:50;default:customer;index" json:"role"`
	IsEmailVerified    bool           `gorm:"default:false" json:"is_email_verified"`
	IsPhoneVerified    bool           `gorm:"default:false" json:"is_phone_verified"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	IsSuspended        bool           `gorm:"default:false" json:"is_suspended"`
	MFAEnabled         bool           `gorm:"default:false" json:"mfa_enabled"`
	MFASecret          string         `gorm:"size:255" json:"-"`
	RecoveryCodes      datatypes.JSONSlice[string] `gorm:"type:jsonb" json:"-"`
	FailedLoginAttempts int           `gorm:"default:0" json:"-"`
	LockedUntil        *time.Time     `json:"-"`
	LastLoginAt        *time.Time     `json:"last_login_at"`
	LastLoginIP        string         `gorm:"size:45" json:"-"`
	LastUserAgent      string         `gorm:"size:500" json:"-"`
	PasswordChangedAt  *time.Time     `json:"-"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefreshToken struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"-"`
	Token      string         `gorm:"not null;uniqueIndex;size:500" json:"-"`
	DeviceID   string         `gorm:"size:255" json:"device_id"`
	DeviceName string         `gorm:"size:255" json:"device_name"`
	IP         string         `gorm:"size:45" json:"ip"`
	UserAgent  string         `gorm:"size:500" json:"user_agent"`
	IsRevoked  bool           `gorm:"default:false" json:"-"`
	ExpiresAt  time.Time      `json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type OTP struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Code      string    `gorm:"not null;size:6" json:"-"`
	Purpose   string    `gorm:"not null;size:50" json:"purpose"`
	IsUsed    bool      `gorm:"default:false" json:"-"`
	ExpiresAt time.Time `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
