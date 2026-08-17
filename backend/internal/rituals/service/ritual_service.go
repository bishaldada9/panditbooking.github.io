package service

import (
	"errors"
	"strings"

	"github.com/bees/hindu-ritual-platform/internal/rituals/dto"
	"github.com/bees/hindu-ritual-platform/internal/rituals/model"
	"github.com/bees/hindu-ritual-platform/internal/rituals/repository"
	"github.com/google/uuid"
)

type RitualService struct {
	repo *repository.RitualRepository
}

func NewRitualService(repo *repository.RitualRepository) *RitualService {
	return &RitualService{repo: repo}
}

func (s *RitualService) CreateCategory(req *dto.CreateRitualCategoryRequest) (*dto.RitualCategoryResponse, error) {
	category := &model.RitualCategory{
		Name:        req.Name,
		Slug:        strings.ToLower(strings.ReplaceAll(req.Name, " ", "-")),
		Description: req.Description,
		Icon:        req.Icon,
	}
	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}
	return s.toCategoryResponse(category), nil
}

func (s *RitualService) GetCategories() ([]dto.RitualCategoryResponse, error) {
	categories, err := s.repo.FindAllCategories()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.RitualCategoryResponse, 0, len(categories))
	for _, c := range categories {
		responses = append(responses, *s.toCategoryResponse(&c))
	}
	return responses, nil
}

func (s *RitualService) CreateRitual(req *dto.CreateRitualRequest) (*dto.RitualResponse, error) {
	ritual := &model.Ritual{
		CategoryID:       req.CategoryID,
		Name:             req.Name,
		Slug:             strings.ToLower(strings.ReplaceAll(req.Name, " ", "-")),
		Description:      req.Description,
		Duration:         req.Duration,
		BasePrice:        req.BasePrice,
		RequiredItems:    req.RequiredItems,
		Procedure:        req.Procedure,
		PanditCommission: req.PanditCommission,
	}
	if err := s.repo.CreateRitual(ritual); err != nil {
		return nil, err
	}
	return s.toRitualResponse(ritual), nil
}

func (s *RitualService) GetRituals() ([]dto.RitualResponse, error) {
	rituals, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	responses := make([]dto.RitualResponse, 0, len(rituals))
	for _, r := range rituals {
		responses = append(responses, *s.toRitualResponse(&r))
	}
	return responses, nil
}

func (s *RitualService) GetRitualsByCategory(categoryID uuid.UUID) ([]dto.RitualResponse, error) {
	rituals, err := s.repo.FindByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	responses := make([]dto.RitualResponse, 0, len(rituals))
	for _, r := range rituals {
		responses = append(responses, *s.toRitualResponse(&r))
	}
	return responses, nil
}

func (s *RitualService) GetRitual(id uuid.UUID) (*dto.RitualResponse, error) {
	ritual, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("ritual not found")
	}
	return s.toRitualResponse(ritual), nil
}

func (s *RitualService) toCategoryResponse(c *model.RitualCategory) *dto.RitualCategoryResponse {
	return &dto.RitualCategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Icon:        c.Icon,
		IsActive:    c.IsActive,
	}
}

func (s *RitualService) toRitualResponse(r *model.Ritual) *dto.RitualResponse {
	categoryName := ""
	if r.Category.ID != uuid.Nil {
		categoryName = r.Category.Name
	}
	return &dto.RitualResponse{
		ID:               r.ID,
		CategoryID:       r.CategoryID,
		CategoryName:     categoryName,
		Name:             r.Name,
		Slug:             r.Slug,
		Description:      r.Description,
		Duration:         r.Duration,
		BasePrice:        r.BasePrice,
		RequiredItems:    r.RequiredItems,
		Procedure:        r.Procedure,
		PanditCommission: r.PanditCommission,
	}
}
