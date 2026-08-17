package service

import (
	"errors"
	"time"

	"github.com/bees/hindu-ritual-platform/internal/pandits/dto"
	"github.com/bees/hindu-ritual-platform/internal/pandits/model"
	"github.com/bees/hindu-ritual-platform/internal/pandits/repository"
	"github.com/google/uuid"
)

type PanditService struct {
	repo *repository.PanditRepository
}

func NewPanditService(repo *repository.PanditRepository) *PanditService {
	return &PanditService{repo: repo}
}

func (s *PanditService) RegisterPandit(userID uuid.UUID, req *dto.RegisterPanditRequest) (*dto.PanditResponse, error) {
	existing, _ := s.repo.FindByUserID(userID)
	if existing != nil {
		return nil, errors.New("pandit profile already exists")
	}

	pandit := &model.Pandit{
		UserID:             userID,
		Bio:                req.Bio,
		ExperienceYears:    req.ExperienceYears,
		Specialization:     req.Specialization,
		Languages:          req.Languages,
		BasePrice:          req.BasePrice,
		ServiceArea:        req.ServiceArea,
		VerificationStatus: model.VerificationPending,
		IsAvailable:        true,
	}

	if err := s.repo.Create(pandit); err != nil {
		return nil, err
	}
	return s.toPanditResponse(pandit), nil
}

func (s *PanditService) GetPanditProfile(userID uuid.UUID) (*dto.PanditResponse, error) {
	pandit, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("pandit profile not found")
	}
	return s.toPanditResponse(pandit), nil
}

func (s *PanditService) GetPanditByID(id uuid.UUID) (*dto.PanditResponse, error) {
	pandit, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("pandit not found")
	}
	return s.toPanditResponse(pandit), nil
}

func (s *PanditService) UpdatePanditProfile(userID uuid.UUID, req *dto.UpdatePanditProfileRequest) (*dto.PanditResponse, error) {
	pandit, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("pandit profile not found")
	}

	if req.Bio != "" {
		pandit.Bio = req.Bio
	}
	if req.ExperienceYears >= 0 {
		pandit.ExperienceYears = req.ExperienceYears
	}
	if req.Specialization != "" {
		pandit.Specialization = req.Specialization
	}
	if len(req.Languages) > 0 {
		pandit.Languages = req.Languages
	}
	if req.BasePrice >= 0 {
		pandit.BasePrice = req.BasePrice
	}
	if req.ServiceArea != "" {
		pandit.ServiceArea = req.ServiceArea
	}
	if req.IsAvailable != nil {
		pandit.IsAvailable = *req.IsAvailable
	}

	if err := s.repo.Update(pandit); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.toPanditResponse(updated), nil
}

func (s *PanditService) ListPandits(page, limit int, filters map[string]interface{}) ([]dto.PanditResponse, int64, error) {
	filters["verification_status"] = model.VerificationApproved
	return s.listPandits(page, limit, filters)
}

func (s *PanditService) ListPanditsForAdmin(page, limit int, filters map[string]interface{}) ([]dto.PanditResponse, int64, error) {
	return s.listPandits(page, limit, filters)
}

func (s *PanditService) listPandits(page, limit int, filters map[string]interface{}) ([]dto.PanditResponse, int64, error) {
	pandits, total, err := s.repo.FindAll(page, limit, filters)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.PanditResponse, 0, len(pandits))
	for _, p := range pandits {
		responses = append(responses, *s.toPanditResponse(&p))
	}
	return responses, total, nil
}

func (s *PanditService) UpdateAvailability(userID uuid.UUID, req *dto.UpdateAvailabilityRequest) error {
	pandit, err := s.repo.FindByUserID(userID)
	if err != nil {
		return errors.New("pandit profile not found")
	}

	avail := &model.Availability{
		PanditID:  pandit.ID,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsBooked:  req.IsBooked,
	}
	return s.repo.CreateAvailability(avail)
}

func (s *PanditService) GetAvailability(panditID uuid.UUID, date string) ([]model.Availability, error) {
	return s.repo.FindAvailability(panditID, date)
}

func (s *PanditService) UploadDocument(userID uuid.UUID, docType, docURL string) error {
	pandit, err := s.repo.FindByUserID(userID)
	if err != nil {
		return errors.New("pandit profile not found")
	}
	doc := &model.PanditDocument{
		PanditID:     pandit.ID,
		DocumentType: docType,
		DocumentURL:  docURL,
	}
	return s.repo.CreateDocument(doc)
}

func (s *PanditService) VerifyPandit(panditID, adminID uuid.UUID, status model.VerificationStatus, notes string) error {
	pandit, err := s.repo.FindByID(panditID)
	if err != nil {
		return errors.New("pandit not found")
	}
	pandit.VerificationStatus = status
	pandit.VerificationNotes = notes
	pandit.VerifiedBy = &adminID
	now := time.Now()
	pandit.VerifiedAt = &now
	return s.repo.Update(pandit)
}

func (s *PanditService) GetPendingVerifications() ([]dto.PanditResponse, error) {
	pandits, err := s.repo.FindPendingVerifications()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.PanditResponse, 0, len(pandits))
	for _, p := range pandits {
		responses = append(responses, *s.toPanditResponse(&p))
	}
	return responses, nil
}

func (s *PanditService) toPanditResponse(pandit *model.Pandit) *dto.PanditResponse {
	return &dto.PanditResponse{
		ID:                 pandit.ID,
		UserID:             pandit.UserID,
		FullName:           pandit.User.FullName,
		Email:              pandit.User.Email,
		Phone:              pandit.User.Phone,
		Bio:                pandit.Bio,
		ExperienceYears:    pandit.ExperienceYears,
		Specialization:     pandit.Specialization,
		Languages:          pandit.Languages,
		BasePrice:          pandit.BasePrice,
		Rating:             pandit.Rating,
		TotalReviews:       pandit.TotalReviews,
		TotalBookings:      pandit.TotalBookings,
		IsAvailable:        pandit.IsAvailable,
		VerificationStatus: pandit.VerificationStatus,
		ServiceArea:        pandit.ServiceArea,
		CreatedAt:          pandit.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
