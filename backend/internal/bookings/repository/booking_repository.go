package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/bookings/model"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Customer").Preload("Pandit").Preload("Ritual").First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindByCustomer(customerID uuid.UUID, page, limit int) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64
	r.db.Model(&model.Booking{}).Where("customer_id = ?", customerID).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("customer_id = ?", customerID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) FindByPandit(panditID uuid.UUID, page, limit int) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64
	r.db.Model(&model.Booking{}).Where("pandit_id = ?", panditID).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("pandit_id = ?", panditID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) FindAll(page, limit int) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64
	r.db.Model(&model.Booking{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) FindByPanditAndDate(panditID uuid.UUID, date string) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where("pandit_id = ? AND scheduled_date = ? AND status NOT IN ?",
		panditID, date, []string{"cancelled", "rejected"}).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) CreateComplaint(complaint *model.Complaint) error {
	return r.db.Create(complaint).Error
}
