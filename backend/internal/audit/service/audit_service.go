package service

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	"github.com/bees/hindu-ritual-platform/internal/audit/model"
	"github.com/bees/hindu-ritual-platform/internal/audit/repository"
)

type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(userID *uuid.UUID, action, resource, resourceID, detail, oldValue, newValue, ip, userAgent, status string) error {
	entry := &model.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Detail:     detail,
		OldValue:   datatypes.JSON([]byte(oldValue)),
		NewValue:   datatypes.JSON([]byte(newValue)),
		IP:         ip,
		UserAgent:  userAgent,
		Status:     status,
	}
	return s.repo.Log(entry)
}

func (s *AuditService) GetLogs(page, limit int) ([]model.AuditLog, int64, error) {
	logs, total, err := s.repo.FindAll(page, limit)
	if logs == nil {
		logs = make([]model.AuditLog, 0)
	}
	return logs, total, err
}

func (s *AuditService) GetUserLogs(userID uuid.UUID, page, limit int) ([]model.AuditLog, int64, error) {
	logs, total, err := s.repo.FindByUserID(userID, page, limit)
	if logs == nil {
		logs = make([]model.AuditLog, 0)
	}
	return logs, total, err
}

func (s *AuditService) GetActionLogs(action string, page, limit int) ([]model.AuditLog, int64, error) {
	logs, total, err := s.repo.FindByAction(action, page, limit)
	if logs == nil {
		logs = make([]model.AuditLog, 0)
	}
	return logs, total, err
}
