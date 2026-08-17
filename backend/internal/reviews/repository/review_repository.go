package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/reviews/model"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByID(id uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.First(&review, "id = ?", id).Error
	return &review, err
}

func (r *ReviewRepository) FindByPandit(panditID uuid.UUID, page, limit int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64
	r.db.Model(&model.Review{}).Where("pandit_id = ? AND is_visible = ?", panditID, true).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("pandit_id = ? AND is_visible = ?", panditID, true).Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) FindByBookingID(bookingID uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("booking_id = ?", bookingID).First(&review).Error
	return &review, err
}

func (r *ReviewRepository) Update(review *model.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) FindAll(page, limit int) ([]model.Review, int64, error) {
	var reviews []model.Review
	var total int64
	r.db.Model(&model.Review{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

func (r *ReviewRepository) GetAverageRating(panditID uuid.UUID) (float64, error) {
	var avg float64
	err := r.db.Model(&model.Review{}).Where("pandit_id = ? AND is_visible = ?", panditID, true).Select("COALESCE(AVG(rating), 0)").Scan(&avg).Error
	return avg, err
}
