package handler

import (
	"net/http"

	"github.com/bees/hindu-ritual-platform/internal/pandits/dto"
	"github.com/bees/hindu-ritual-platform/internal/pandits/model"
	"github.com/bees/hindu-ritual-platform/internal/pandits/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PanditHandler struct {
	service *service.PanditService
}

func NewPanditHandler(service *service.PanditService) *PanditHandler {
	return &PanditHandler{service: service}
}

func (h *PanditHandler) Register(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.RegisterPanditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.RegisterPandit(userID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusConflict, "Registration failed", err.Error())
		return
	}
	response.Created(c, "Pandit profile created", result)
}

func (h *PanditHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.service.GetPanditProfile(userID.(uuid.UUID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Profile not found", err.Error())
		return
	}
	response.Success(c, "Profile retrieved", profile)
}

func (h *PanditHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.UpdatePanditProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	profile, err := h.service.UpdatePanditProfile(userID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Profile update failed", err.Error())
		return
	}
	response.Success(c, "Pandit profile updated", profile)
}

func (h *PanditHandler) GetPandit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid pandit ID", err.Error())
		return
	}
	pandit, err := h.service.GetPanditByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Pandit not found", err.Error())
		return
	}
	response.Success(c, "Pandit retrieved", pandit)
}

func (h *PanditHandler) List(c *gin.Context) {
	page := 1
	limit := 20
	filters := make(map[string]interface{})
	if specialization := c.Query("specialization"); specialization != "" {
		filters["specialization"] = specialization
	}
	pandits, total, err := h.service.ListPandits(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list pandits", err.Error())
		return
	}
	response.Paginated(c, "Pandits retrieved", pandits, total, page, limit)
}

func (h *PanditHandler) ListForAdmin(c *gin.Context) {
	page := 1
	limit := 100
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["verification_status"] = status
	}
	pandits, total, err := h.service.ListPanditsForAdmin(page, limit, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list pandits", err.Error())
		return
	}
	response.Paginated(c, "Pandits retrieved", pandits, total, page, limit)
}

func (h *PanditHandler) UpdateAvailability(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	if err := h.service.UpdateAvailability(userID.(uuid.UUID), &req); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update availability", err.Error())
		return
	}
	response.Success(c, "Availability updated", nil)
}

func (h *PanditHandler) GetAvailability(c *gin.Context) {
	panditID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid pandit ID", err.Error())
		return
	}
	date := c.Query("date")
	avail, err := h.service.GetAvailability(panditID, date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get availability", err.Error())
		return
	}
	response.Success(c, "Availability retrieved", avail)
}

func (h *PanditHandler) VerifyPandit(c *gin.Context) {
	panditID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid pandit ID", err.Error())
		return
	}
	adminID, _ := c.Get("user_id")
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending approved rejected"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request", err.Error())
		return
	}
	status := model.VerificationStatus(req.Status)
	if err := h.service.VerifyPandit(panditID, adminID.(uuid.UUID), status, req.Notes); err != nil {
		response.Error(c, http.StatusBadRequest, "Verification failed", err.Error())
		return
	}
	response.Success(c, "Pandit verification updated", nil)
}
