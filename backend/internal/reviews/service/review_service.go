package service

import (
	"errors"
	bookingRepo "github.com/bees/hindu-ritual-platform/internal/bookings/repository"
	panditRepo "github.com/bees/hindu-ritual-platform/internal/pandits/repository"
	"github.com/bees/hindu-ritual-platform/internal/reviews/dto"
	"github.com/bees/hindu-ritual-platform/internal/reviews/model"
	"github.com/bees/hindu-ritual-platform/internal/reviews/repository"
	"github.com/google/uuid"
)

type ReviewService struct {
	repo        *repository.ReviewRepository
	bookingRepo *bookingRepo.BookingRepository
	panditRepo  *panditRepo.PanditRepository
}

func NewReviewService(repo *repository.ReviewRepository, bookingRepo *bookingRepo.BookingRepository, panditRepo *panditRepo.PanditRepository) *ReviewService {
	return &ReviewService{
		repo:        repo,
		bookingRepo: bookingRepo,
		panditRepo:  panditRepo,
	}
}

func (s *ReviewService) CreateReview(customerID uuid.UUID, req *dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	booking, err := s.bookingRepo.FindByID(req.BookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	if booking.CustomerID != customerID {
		return nil, errors.New("unauthorized")
	}
	if booking.Status != "completed" {
		return nil, errors.New("can only review completed bookings")
	}

	existing, _ := s.repo.FindByBookingID(req.BookingID)
	if existing != nil {
		return nil, errors.New("already reviewed this booking")
	}

	review := &model.Review{
		BookingID:  req.BookingID,
		CustomerID: customerID,
		PanditID:   booking.PanditID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsVisible:  true,
	}

	if err := s.repo.Create(review); err != nil {
		return nil, err
	}

	avg, _ := s.repo.GetAverageRating(booking.PanditID)
	pandit, _ := s.panditRepo.FindByID(booking.PanditID)
	if pandit != nil {
		pandit.Rating = avg
		pandit.TotalReviews++
		s.panditRepo.Update(pandit)
	}

	return s.toReviewResponse(review), nil
}

func (s *ReviewService) GetPanditReviews(panditID uuid.UUID, page, limit int) ([]dto.ReviewResponse, int64, error) {
	reviews, total, err := s.repo.FindByPandit(panditID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		responses = append(responses, *s.toReviewResponse(&r))
	}
	return responses, total, nil
}

func (s *ReviewService) ModerateReview(reviewID, adminID uuid.UUID, isVisible bool, adminReply string) error {
	review, err := s.repo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}
	review.IsVisible = isVisible
	review.AdminReply = adminReply
	return s.repo.Update(review)
}

func (s *ReviewService) ListReviews(page, limit int) ([]dto.ReviewResponse, int64, error) {
	reviews, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		responses = append(responses, *s.toReviewResponse(&r))
	}
	return responses, total, nil
}

func (s *ReviewService) toReviewResponse(r *model.Review) *dto.ReviewResponse {
	return &dto.ReviewResponse{
		ID:         r.ID,
		BookingID:  r.BookingID,
		CustomerID: r.CustomerID,
		PanditID:   r.PanditID,
		Rating:     r.Rating,
		Comment:    r.Comment,
		IsVerified: r.IsVerified,
		IsVisible:  r.IsVisible,
		AdminReply: r.AdminReply,
		CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
