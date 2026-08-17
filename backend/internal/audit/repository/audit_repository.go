package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/audit/model"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Log(entry *model.AuditLog) error {
	return r.db.Create(entry).Error
}

func (r *AuditRepository) FindAll(page, limit int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	r.db.Model(&model.AuditLog{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func (r *AuditRepository) FindByUserID(userID uuid.UUID, page, limit int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	r.db.Model(&model.AuditLog{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func (r *AuditRepository) FindByAction(action string, page, limit int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	r.db.Model(&model.AuditLog{}).Where("action = ?", action).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Where("action = ?", action).Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}
