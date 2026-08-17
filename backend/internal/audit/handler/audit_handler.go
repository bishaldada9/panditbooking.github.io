package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/bees/hindu-ritual-platform/internal/audit/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler(service *service.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if userID := c.Query("user_id"); userID != "" {
		id, err := uuid.Parse(userID)
		if err == nil {
			logs, total, err := h.service.GetUserLogs(id, page, limit)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, "Failed to get logs", err.Error())
				return
			}
			response.Paginated(c, "Audit logs retrieved", logs, total, page, limit)
			return
		}
	}

	if action := c.Query("action"); action != "" {
		logs, total, err := h.service.GetActionLogs(action, page, limit)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to get logs", err.Error())
			return
		}
		response.Paginated(c, "Audit logs retrieved", logs, total, page, limit)
		return
	}

	logs, total, err := h.service.GetLogs(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get logs", err.Error())
		return
	}
	response.Paginated(c, "Audit logs retrieved", logs, total, page, limit)
}
