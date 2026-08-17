package service

import (
	"github.com/bees/hindu-ritual-platform/internal/notification/dto"
	"github.com/bees/hindu-ritual-platform/internal/notification/model"
	"github.com/bees/hindu-ritual-platform/internal/notification/repository"
	"github.com/google/uuid"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(userID uuid.UUID, notifType, title, message, refID, refType string) error {
	notification := &model.Notification{
		UserID:        userID,
		Type:          notifType,
		Title:         title,
		Message:       message,
		ReferenceID:   refID,
		ReferenceType: refType,
	}
	return s.repo.Create(notification)
}

func (s *NotificationService) GetUserNotifications(userID uuid.UUID, page, limit int) ([]dto.NotificationResponse, int64, error) {
	notifications, total, err := s.repo.FindByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		responses = append(responses, *s.toNotificationResponse(&n))
	}
	return responses, total, nil
}

func (s *NotificationService) MarkAsRead(id uuid.UUID) error {
	return s.repo.MarkAsRead(id)
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

func (s *NotificationService) toNotificationResponse(n *model.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
		ID:            n.ID,
		UserID:        n.UserID,
		Type:          n.Type,
		Title:         n.Title,
		Message:       n.Message,
		IsRead:        n.IsRead,
		ReferenceID:   n.ReferenceID,
		ReferenceType: n.ReferenceType,
		CreatedAt:     n.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
