package handler

import (
	"net/http"
	"strconv"

	"github.com/bees/hindu-ritual-platform/internal/bookings/dto"
	"github.com/bees/hindu-ritual-platform/internal/bookings/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.CreateBooking(userID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Booking failed", err.Error())
		return
	}
	response.Created(c, "Booking created", result)
}

func (h *BookingHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid booking ID", err.Error())
		return
	}
	booking, err := h.service.GetBooking(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Booking not found", err.Error())
		return
	}
	response.Success(c, "Booking retrieved", booking)
}

func (h *BookingHandler) ListMyBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	bookings, total, err := h.service.GetCustomerBookings(userID.(uuid.UUID), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get bookings", err.Error())
		return
	}
	response.Paginated(c, "Bookings retrieved", bookings, total, page, limit)
}

func (h *BookingHandler) ListPanditBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	bookings, total, err := h.service.GetPanditBookings(userID.(uuid.UUID), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get bookings", err.Error())
		return
	}
	response.Paginated(c, "Bookings retrieved", bookings, total, page, limit)
}

func (h *BookingHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	bookings, total, err := h.service.GetAllBookings(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get bookings", err.Error())
		return
	}
	response.Paginated(c, "Bookings retrieved", bookings, total, page, limit)
}

func (h *BookingHandler) Confirm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid booking ID", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.service.ConfirmBooking(id, userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusBadRequest, "Confirmation failed", err.Error())
		return
	}
	response.Success(c, "Booking confirmed", nil)
}

func (h *BookingHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid booking ID", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.service.CompleteBooking(id, userID.(uuid.UUID)); err != nil {
		response.Error(c, http.StatusBadRequest, "Completion failed", err.Error())
		return
	}
	response.Success(c, "Booking completed", nil)
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid booking ID", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	var req dto.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	if err := h.service.CancelBooking(id, userID.(uuid.UUID), role.(string), req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, "Cancellation failed", err.Error())
		return
	}
	response.Success(c, "Booking cancelled", nil)
}

func (h *BookingHandler) Reject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid booking ID", err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	var req dto.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	if err := h.service.RejectBooking(id, userID.(uuid.UUID), req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, "Rejection failed", err.Error())
		return
	}
	response.Success(c, "Booking rejected", nil)
}
