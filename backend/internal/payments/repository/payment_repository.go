package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/payments/model"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) FindByID(id uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.First(&payment, "id = ?", id).Error
	return &payment, err
}

func (r *PaymentRepository) FindByBookingID(bookingID uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("booking_id = ?", bookingID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) FindByTransactionID(txID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("transaction_id = ?", txID).First(&payment).Error
	return &payment, err
}

func (r *PaymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) FindAll(page, limit int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64
	r.db.Model(&model.Payment{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error
	return payments, total, err
}

func (r *PaymentRepository) CreateTransaction(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}
