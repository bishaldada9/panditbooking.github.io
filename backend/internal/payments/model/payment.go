package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	bookingModel "github.com/bees/hindu-ritual-platform/internal/bookings/model"
)

type PaymentStatus string
type PaymentGateway string

const (
	PaymentPending           PaymentStatus = "pending"
	PaymentCompleted         PaymentStatus = "completed"
	PaymentFailed            PaymentStatus = "failed"
	PaymentRefunded          PaymentStatus = "refunded"
	PaymentPartiallyRefunded PaymentStatus = "partially_refunded"

	GatewayEsewa  PaymentGateway = "esewa"
	GatewayKhalti PaymentGateway = "khalti"
	GatewayMock   PaymentGateway = "mock"
	GatewayCash   PaymentGateway = "cash"
)

type Payment struct {
	ID              uuid.UUID            `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BookingID       uuid.UUID            `gorm:"type:uuid;uniqueIndex;not null" json:"booking_id"`
	CustomerID      uuid.UUID            `gorm:"type:uuid;not null;index" json:"customer_id"`
	Booking         bookingModel.Booking `gorm:"foreignKey:BookingID" json:"-"`
	Customer        authModel.User       `gorm:"foreignKey:CustomerID" json:"-"`
	Amount          float64              `gorm:"not null" json:"amount"`
	Gateway         PaymentGateway       `gorm:"not null;size:20" json:"gateway"`
	GatewayRefID    string               `gorm:"size:255" json:"gateway_ref_id"`
	Status          PaymentStatus        `gorm:"not null;default:pending;index" json:"status"`
	TransactionID   string               `gorm:"uniqueIndex;size:255" json:"transaction_id"`
	VerifiedAt      *time.Time           `json:"verified_at"`
	RefundedAt      *time.Time           `json:"refunded_at"`
	RefundAmount    float64              `json:"refund_amount"`
	RefundReason    string               `gorm:"type:text" json:"refund_reason"`
	GatewayResponse datatypes.JSON       `gorm:"type:jsonb" json:"-"`
	WebhookReceived bool                 `gorm:"default:false" json:"webhook_received"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

type Transaction struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PaymentID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"payment_id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Payment       Payment        `gorm:"foreignKey:PaymentID" json:"-"`
	User          authModel.User `gorm:"foreignKey:UserID" json:"-"`
	Type          string         `gorm:"not null;size:50" json:"type"`
	Amount        float64        `gorm:"not null" json:"amount"`
	BalanceBefore float64        `json:"balance_before"`
	BalanceAfter  float64        `json:"balance_after"`
	Description   string         `gorm:"type:text" json:"description"`
	ReferenceID   string         `gorm:"size:255" json:"reference_id"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
