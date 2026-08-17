package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/bees/hindu-ritual-platform/internal/rituals/dto"
	"github.com/bees/hindu-ritual-platform/internal/rituals/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type RitualHandler struct {
	service *service.RitualService
}

func NewRitualHandler(service *service.RitualService) *RitualHandler {
	return &RitualHandler{service: service}
}

func (h *RitualHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateRitualCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.CreateCategory(&req)
	if err != nil {
		response.Error(c, http.StatusConflict, "Category creation failed", err.Error())
		return
	}
	response.Created(c, "Category created", result)
}

func (h *RitualHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get categories", err.Error())
		return
	}
	response.Success(c, "Categories retrieved", categories)
}

func (h *RitualHandler) CreateRitual(c *gin.Context) {
	var req dto.CreateRitualRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.CreateRitual(&req)
	if err != nil {
		response.Error(c, http.StatusConflict, "Ritual creation failed", err.Error())
		return
	}
	response.Created(c, "Ritual created", result)
}

func (h *RitualHandler) GetRituals(c *gin.Context) {
	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			response.ValidationError(c, "Invalid category ID", err.Error())
			return
		}
		rituals, err := h.service.GetRitualsByCategory(id)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to get rituals", err.Error())
			return
		}
		response.Success(c, "Rituals retrieved", rituals)
		return
	}
	rituals, err := h.service.GetRituals()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get rituals", err.Error())
		return
	}
	response.Success(c, "Rituals retrieved", rituals)
}

func (h *RitualHandler) GetRitual(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid ritual ID", err.Error())
		return
	}
	ritual, err := h.service.GetRitual(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Ritual not found", err.Error())
		return
	}
	response.Success(c, "Ritual retrieved", ritual)
}
