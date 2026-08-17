package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/notification/model"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) FindByUserID(userID uuid.UUID, page, limit int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64
	r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error
	return notifications, total, err
}

func (r *NotificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
