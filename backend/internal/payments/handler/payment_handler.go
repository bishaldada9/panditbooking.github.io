package handler

import (
	"net/http"
	"io"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/bees/hindu-ritual-platform/internal/payments/dto"
	"github.com/bees/hindu-ritual-platform/internal/payments/service"
	"github.com/bees/hindu-ritual-platform/pkg/response"
)

type PaymentHandler struct {
	service *service.PaymentService
}

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Initiate(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.InitiatePayment(userID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Payment initiation failed", err.Error())
		return
	}
	response.Success(c, "Payment initiated", result)
}

func (h *PaymentHandler) Verify(c *gin.Context) {
	var req dto.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.VerifyPayment(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Payment verification failed", err.Error())
		return
	}
	response.Success(c, "Payment verified", result)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ValidationError(c, "Invalid payment ID", err.Error())
		return
	}
	payment, err := h.service.GetPayment(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found", err.Error())
		return
	}
	response.Success(c, "Payment retrieved", payment)
}

func (h *PaymentHandler) GetPaymentByBooking(c *gin.Context) {
	bookingID, err := uuid.Parse(c.Param("bookingId"))
	if err != nil {
		response.ValidationError(c, "Invalid booking ID", err.Error())
		return
	}
	payment, err := h.service.GetPaymentByBooking(bookingID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found", err.Error())
		return
	}
	response.Success(c, "Payment retrieved", payment)
}

func (h *PaymentHandler) Refund(c *gin.Context) {
	adminID, _ := c.Get("user_id")
	var req dto.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "Invalid request body", err.Error())
		return
	}
	result, err := h.service.RefundPayment(adminID.(uuid.UUID), &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Refund failed", err.Error())
		return
	}
	response.Success(c, "Refund processed", result)
}

func (h *PaymentHandler) Webhook(c *gin.Context) {
	gateway := c.Param("gateway")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid webhook payload", err.Error())
		return
	}
	if err := h.service.HandleWebhook(gateway, body); err != nil {
		response.Error(c, http.StatusBadRequest, "Webhook processing failed", err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *PaymentHandler) List(c *gin.Context) {
	page := 1
	limit := 20
	payments, total, err := h.service.ListPayments(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list payments", err.Error())
		return
	}
	response.Paginated(c, "Payments retrieved", payments, total, page, limit)
}
