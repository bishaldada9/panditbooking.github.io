package service

import (
	"errors"
	"time"

	"github.com/bees/hindu-ritual-platform/internal/bookings/dto"
	"github.com/bees/hindu-ritual-platform/internal/bookings/model"
	"github.com/bees/hindu-ritual-platform/internal/bookings/repository"
	panditRepo "github.com/bees/hindu-ritual-platform/internal/pandits/repository"
	ritualRepo "github.com/bees/hindu-ritual-platform/internal/rituals/repository"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
	"github.com/google/uuid"
)

type BookingService struct {
	repo       *repository.BookingRepository
	panditRepo *panditRepo.PanditRepository
	ritualRepo *ritualRepo.RitualRepository
	log        *logger.Logger
}

func NewBookingService(repo *repository.BookingRepository, panditRepo *panditRepo.PanditRepository, ritualRepo *ritualRepo.RitualRepository, log *logger.Logger) *BookingService {
	return &BookingService{
		repo:       repo,
		panditRepo: panditRepo,
		ritualRepo: ritualRepo,
		log:        log,
	}
}

func (s *BookingService) CreateBooking(customerID uuid.UUID, req *dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	pandit, err := s.panditRepo.FindByID(req.PanditID)
	if err != nil {
		return nil, errors.New("pandit not found")
	}
	if pandit.VerificationStatus != "approved" {
		return nil, errors.New("pandit is not verified")
	}

	ritual, err := s.ritualRepo.FindByID(req.RitualID)
	if err != nil {
		return nil, errors.New("ritual not found")
	}

	platformFee := ritual.BasePrice * 0.10
	totalAmount := ritual.BasePrice + platformFee

	booking := &model.Booking{
		CustomerID:    customerID,
		PanditID:      req.PanditID,
		RitualID:      req.RitualID,
		Status:        model.BookingPending,
		ScheduledDate: req.ScheduledDate,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Address:       req.Address,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		SpecialNotes:  req.SpecialNotes,
		TotalAmount:   totalAmount,
		PlatformFee:   platformFee,
		PanditAmount:  ritual.BasePrice,
	}

	if err := s.repo.Create(booking); err != nil {
		return nil, err
	}

	return s.toBookingResponse(booking), nil
}

func (s *BookingService) GetBooking(id uuid.UUID) (*dto.BookingResponse, error) {
	booking, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	return s.toBookingResponse(booking), nil
}

func (s *BookingService) GetCustomerBookings(customerID uuid.UUID, page, limit int) ([]dto.BookingResponse, int64, error) {
	bookings, total, err := s.repo.FindByCustomer(customerID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		responses = append(responses, *s.toBookingResponse(&b))
	}
	return responses, total, nil
}

func (s *BookingService) GetPanditBookings(panditUserID uuid.UUID, page, limit int) ([]dto.BookingResponse, int64, error) {
	pandit, err := s.panditRepo.FindByUserID(panditUserID)
	if err != nil {
		return nil, 0, errors.New("pandit profile not found")
	}
	bookings, total, err := s.repo.FindByPandit(pandit.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		responses = append(responses, *s.toBookingResponse(&b))
	}
	return responses, total, nil
}

func (s *BookingService) ConfirmBooking(bookingID, panditUserID uuid.UUID) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	pandit, err := s.panditRepo.FindByUserID(panditUserID)
	if err != nil || booking.PanditID != pandit.ID {
		return errors.New("unauthorized")
	}
	if booking.Status != model.BookingPending {
		return errors.New("booking cannot be confirmed in current status")
	}
	now := time.Now()
	booking.Status = model.BookingConfirmed
	booking.ConfirmedAt = &now
	return s.repo.Update(booking)
}

func (s *BookingService) CompleteBooking(bookingID, panditUserID uuid.UUID) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	pandit, err := s.panditRepo.FindByUserID(panditUserID)
	if err != nil || booking.PanditID != pandit.ID {
		return errors.New("unauthorized")
	}
	if booking.Status != model.BookingConfirmed {
		return errors.New("booking cannot be completed in current status")
	}
	now := time.Now()
	booking.Status = model.BookingCompleted
	booking.CompletedAt = &now
	return s.repo.Update(booking)
}

func (s *BookingService) CancelBooking(bookingID, userID uuid.UUID, role string, reason string) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	if booking.CustomerID != userID && role != "admin" {
		return errors.New("unauthorized to cancel this booking")
	}
	if booking.Status == model.BookingCompleted || booking.Status == model.BookingCancelled {
		return errors.New("booking cannot be cancelled in current status")
	}
	now := time.Now()
	booking.Status = model.BookingCancelled
	booking.CancelledBy = &userID
	booking.CancelReason = reason
	booking.CancelledAt = &now
	return s.repo.Update(booking)
}

func (s *BookingService) RejectBooking(bookingID, panditUserID uuid.UUID, reason string) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}
	pandit, err := s.panditRepo.FindByUserID(panditUserID)
	if err != nil || booking.PanditID != pandit.ID {
		return errors.New("unauthorized")
	}
	if booking.Status != model.BookingPending {
		return errors.New("booking cannot be rejected in current status")
	}
	booking.Status = model.BookingRejected
	booking.CancelReason = reason
	return s.repo.Update(booking)
}

func (s *BookingService) GetAllBookings(page, limit int) ([]dto.BookingResponse, int64, error) {
	bookings, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		responses = append(responses, *s.toBookingResponse(&b))
	}
	return responses, total, nil
}

func (s *BookingService) toBookingResponse(booking *model.Booking) *dto.BookingResponse {
	return &dto.BookingResponse{
		ID:            booking.ID,
		CustomerID:    booking.CustomerID,
		PanditID:      booking.PanditID,
		RitualID:      booking.RitualID,
		Status:        booking.Status,
		ScheduledDate: booking.ScheduledDate,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		Address:       booking.Address,
		TotalAmount:   booking.TotalAmount,
		PlatformFee:   booking.PlatformFee,
		SpecialNotes:  booking.SpecialNotes,
		CreatedAt:     booking.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
