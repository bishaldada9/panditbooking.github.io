package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bees/hindu-ritual-platform/internal/admin/dto"
	"github.com/bees/hindu-ritual-platform/internal/admin/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(service *service.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	metrics, err := h.service.GetDashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get dashboard", err.Error())
		return
	}
	response.Success(c, "Dashboard data retrieved", metrics)
}

func (h *AdminHandler) SuspendUser(c *gin.Context) {
	adminID, _ := c.Get("user_id")
	userID := c.Param("id")
	var req dto.SuspendUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	if err := h.service.SuspendUser(adminID.(uuid.UUID).String(), userID, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to suspend user", err.Error())
		return
	}
	response.Success(c, "User suspended", nil)
}

func (h *AdminHandler) ActivateUser(c *gin.Context) {
	adminID, _ := c.Get("user_id")
	userID := c.Param("id")
	if err := h.service.ActivateUser(adminID.(uuid.UUID).String(), userID); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to activate user", err.Error())
		return
	}
	response.Success(c, "User activated", nil)
}
