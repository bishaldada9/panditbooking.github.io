package dto

import "github.com/google/uuid"

type InitiatePaymentRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Gateway   string    `json:"gateway" validate:"required,oneof=esewa khalti cash mock"`
}

type PaymentResponse struct {
	ID            uuid.UUID `json:"id"`
	BookingID     uuid.UUID `json:"booking_id"`
	Amount        float64   `json:"amount"`
	Gateway       string    `json:"gateway"`
	Status        string    `json:"status"`
	TransactionID string    `json:"transaction_id"`
	GatewayRefID  string    `json:"gateway_ref_id"`
	GatewayURL    string    `json:"gateway_url,omitempty"`
	CreatedAt     string    `json:"created_at"`
}

type VerifyPaymentRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
	GatewayRefID  string `json:"gateway_ref_id"`
	Status        string `json:"status" validate:"required"`
}

type RefundRequest struct {
	PaymentID uuid.UUID `json:"payment_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,min=0"`
	Reason    string    `json:"reason" validate:"required"`
}
