package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/bees/hindu-ritual-platform/internal/notification/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type NotificationHandler struct {
	service *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	notifications, total, err := h.service.GetUserNotifications(userID.(uuid.UUID), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get notifications", err.Error())
		return
	}
	response.Paginated(c, "Notifications retrieved", notifications, total, page, limit)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid notification ID", err.Error())
		return
	}
	if err := h.service.MarkAsRead(id); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to mark as read", err.Error())
		return
	}
	response.Success(c, "Notification marked as read", nil)
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.service.MarkAllAsRead(userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to mark all as read", err.Error())
		return
	}
	response.Success(c, "All notifications marked as read", nil)
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	count, err := h.service.GetUnreadCount(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get count", err.Error())
		return
	}
	response.Success(c, "Unread count retrieved", map[string]int64{"count": count})
}
