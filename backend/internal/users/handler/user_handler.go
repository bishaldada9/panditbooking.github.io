package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/bees/hindu-ritual-platform/internal/users/dto"
	"github.com/bees/hindu-ritual-platform/internal/users/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.service.GetProfile(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Profile not found", err.Error())
		return
	}
	response.Success(c, "Profile retrieved", profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	profile, err := h.service.UpdateProfile(userID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Profile update failed", err.Error())
		return
	}
	response.Success(c, "Profile updated", profile)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page := 1
	limit := 20
	users, total, err := h.service.ListUsers(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list users", err.Error())
		return
	}
	response.Paginated(c, "Users retrieved", users, total, page, limit)
}
