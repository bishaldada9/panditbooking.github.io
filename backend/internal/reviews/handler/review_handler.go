package handler

import (
	"github.com/bees/hindu-ritual-platform/internal/reviews/dto"
	"github.com/bees/hindu-ritual-platform/internal/reviews/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.CreateReview(userID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Review creation failed", err.Error())
		return
	}
	response.Created(c, "Review created", result)
}

func (h *ReviewHandler) GetPanditReviews(c *gin.Context) {
	panditID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid pandit ID", err.Error())
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	reviews, total, err := h.service.GetPanditReviews(panditID, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get reviews", err.Error())
		return
	}
	response.Paginated(c, "Reviews retrieved", reviews, total, page, limit)
}

func (h *ReviewHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	reviews, total, err := h.service.ListReviews(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get reviews", err.Error())
		return
	}
	response.Paginated(c, "Reviews retrieved", reviews, total, page, limit)
}

func (h *ReviewHandler) Moderate(c *gin.Context) {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid review ID", err.Error())
		return
	}
	adminID, _ := c.Get("user_id")
	var req struct {
		IsVisible  bool   `json:"is_visible"`
		AdminReply string `json:"admin_reply"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	if err := h.service.ModerateReview(reviewID, adminID.(uuid.UUID), req.IsVisible, req.AdminReply); err != nil {
		response.Error(c, http.StatusBadRequest, "Moderation failed", err.Error())
		return
	}
	response.Success(c, "Review moderated", nil)
}
