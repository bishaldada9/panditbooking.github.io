package service

import (
	"errors"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	"github.com/bees/hindu-ritual-platform/internal/users/dto"
	"github.com/bees/hindu-ritual-platform/internal/users/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	user.FullName = req.FullName
	user.Phone = req.Phone
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *UserService) ListUsers(page, limit int) ([]dto.UserResponse, int64, error) {
	users, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *toUserResponse(&u))
	}
	return responses, total, nil
}

func toUserResponse(user *authModel.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		FullName:        user.FullName,
		Phone:           user.Phone,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		IsActive:        user.IsActive,
		IsSuspended:     user.IsSuspended,
		CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
