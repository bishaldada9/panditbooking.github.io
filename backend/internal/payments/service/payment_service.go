package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	bookingRepo "github.com/bees/hindu-ritual-platform/internal/bookings/repository"
	"github.com/bees/hindu-ritual-platform/internal/payments/dto"
	"github.com/bees/hindu-ritual-platform/internal/payments/model"
	"github.com/bees/hindu-ritual-platform/internal/payments/repository"
	"github.com/bees/hindu-ritual-platform/pkg/configs"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
)

type PaymentService struct {
	repo        *repository.PaymentRepository
	bookingRepo *bookingRepo.BookingRepository
	config      *configs.Config
	log         *logger.Logger
	gateways    map[string]PaymentGateway
}

type PaymentGateway interface {
	Initiate(amount float64, bookingID, transactionID string) (string, error)
	Verify(transactionID, refID string) (bool, error)
	Refund(transactionID string, amount float64) error
}

func NewPaymentService(repo *repository.PaymentRepository, bookingRepo *bookingRepo.BookingRepository, config *configs.Config, log *logger.Logger) *PaymentService {
	svc := &PaymentService{
		repo:        repo,
		bookingRepo: bookingRepo,
		config:      config,
		log:         log,
		gateways:    make(map[string]PaymentGateway),
	}
	svc.gateways["mock"] = &MockGateway{}
	svc.gateways["esewa"] = &MockGateway{}
	svc.gateways["khalti"] = &MockGateway{}
	svc.gateways["cash"] = &MockGateway{}
	return svc
}

func (s *PaymentService) InitiatePayment(customerID uuid.UUID, req *dto.InitiatePaymentRequest) (*dto.PaymentResponse, error) {
	booking, err := s.bookingRepo.FindByID(req.BookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	if booking.CustomerID != customerID {
		return nil, errors.New("unauthorized")
	}
	if booking.Status != "pending" && booking.Status != "confirmed" {
		return nil, errors.New("booking cannot be paid in current status")
	}

	existingPayment, _ := s.repo.FindByBookingID(req.BookingID)
	if existingPayment != nil && existingPayment.Status == model.PaymentCompleted {
		return nil, errors.New("booking already paid")
	}

	gateway, ok := s.gateways[req.Gateway]
	if !ok {
		return nil, errors.New("unsupported payment gateway")
	}

	transactionID := uuid.New().String()
	gatewayURL, err := gateway.Initiate(booking.TotalAmount, req.BookingID.String(), transactionID)
	if err != nil {
		return nil, fmt.Errorf("gateway initiation failed: %w", err)
	}

	payment := &model.Payment{
		BookingID:     req.BookingID,
		CustomerID:    customerID,
		Amount:        booking.TotalAmount,
		Gateway:       model.PaymentGateway(req.Gateway),
		Status:        model.PaymentPending,
		TransactionID: transactionID,
	}
	if req.Gateway == string(model.GatewayCash) {
		now := time.Now()
		payment.Status = model.PaymentCompleted
		payment.GatewayRefID = "cash-on-site"
		payment.VerifiedAt = &now
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		ID:            payment.ID,
		BookingID:     payment.BookingID,
		Amount:        payment.Amount,
		Gateway:       string(payment.Gateway),
		Status:        string(payment.Status),
		TransactionID: payment.TransactionID,
		GatewayURL:    gatewayURL,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *PaymentService) VerifyPayment(req *dto.VerifyPaymentRequest) (*dto.PaymentResponse, error) {
	payment, err := s.repo.FindByTransactionID(req.TransactionID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	gateway, ok := s.gateways[string(payment.Gateway)]
	if !ok {
		return nil, errors.New("unsupported payment gateway")
	}

	verified, err := gateway.Verify(req.TransactionID, req.GatewayRefID)
	if err != nil || !verified {
		payment.Status = model.PaymentFailed
		s.repo.Update(payment)
		return nil, errors.New("payment verification failed")
	}

	now := time.Now()
	payment.Status = model.PaymentCompleted
	payment.GatewayRefID = req.GatewayRefID
	payment.VerifiedAt = &now
	payment.WebhookReceived = true

	if err := s.repo.Update(payment); err != nil {
		return nil, err
	}

	tx := &model.Transaction{
		PaymentID:   payment.ID,
		UserID:      payment.CustomerID,
		Type:        "payment",
		Amount:      payment.Amount,
		Description: fmt.Sprintf("Payment for booking %s", payment.BookingID),
		ReferenceID: payment.TransactionID,
	}
	s.repo.CreateTransaction(tx)

	return s.toPaymentResponse(payment), nil
}

func (s *PaymentService) GetPayment(id uuid.UUID) (*dto.PaymentResponse, error) {
	payment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return s.toPaymentResponse(payment), nil
}

func (s *PaymentService) GetPaymentByBooking(bookingID uuid.UUID) (*dto.PaymentResponse, error) {
	payment, err := s.repo.FindByBookingID(bookingID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return s.toPaymentResponse(payment), nil
}

func (s *PaymentService) RefundPayment(adminID uuid.UUID, req *dto.RefundRequest) (*dto.PaymentResponse, error) {
	payment, err := s.repo.FindByID(req.PaymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	if payment.Status != model.PaymentCompleted {
		return nil, errors.New("payment cannot be refunded in current status")
	}
	if req.Amount <= 0 {
		req.Amount = payment.Amount
	}
	if req.Amount > payment.Amount {
		return nil, errors.New("refund amount cannot exceed payment amount")
	}

	gateway, ok := s.gateways[string(payment.Gateway)]
	if !ok {
		return nil, errors.New("unsupported payment gateway")
	}

	if err := gateway.Refund(payment.TransactionID, req.Amount); err != nil {
		return nil, fmt.Errorf("refund failed: %w", err)
	}

	now := time.Now()
	payment.Status = model.PaymentRefunded
	payment.RefundedAt = &now
	payment.RefundAmount = req.Amount
	payment.RefundReason = req.Reason

	if err := s.repo.Update(payment); err != nil {
		return nil, err
	}

	tx := &model.Transaction{
		PaymentID:   payment.ID,
		UserID:      adminID,
		Type:        "refund",
		Amount:      -req.Amount,
		Description: fmt.Sprintf("Refund for booking %s: %s", payment.BookingID, req.Reason),
		ReferenceID: payment.TransactionID,
	}
	s.repo.CreateTransaction(tx)

	return s.toPaymentResponse(payment), nil
}

func (s *PaymentService) HandleWebhook(gateway string, payload []byte) error {
	s.log.Info("Webhook received", zap.String("gateway", gateway), zap.String("payload", string(payload)))
	return nil
}

func (s *PaymentService) ListPayments(page, limit int) ([]dto.PaymentResponse, int64, error) {
	payments, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		responses = append(responses, *s.toPaymentResponse(&p))
	}
	return responses, total, nil
}

func (s *PaymentService) toPaymentResponse(p *model.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:            p.ID,
		BookingID:     p.BookingID,
		Amount:        p.Amount,
		Gateway:       string(p.Gateway),
		Status:        string(p.Status),
		TransactionID: p.TransactionID,
		GatewayRefID:  p.GatewayRefID,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

type MockGateway struct{}

func (g *MockGateway) Initiate(amount float64, bookingID, transactionID string) (string, error) {
	return fmt.Sprintf("https://mock-gateway.test/pay?tx=%s&amount=%.2f", transactionID, amount), nil
}

func (g *MockGateway) Verify(transactionID, refID string) (bool, error) {
	return true, nil
}

func (g *MockGateway) Refund(transactionID string, amount float64) error {
	return nil
}
