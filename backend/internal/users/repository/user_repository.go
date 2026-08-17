package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uuid.UUID) (*authModel.User, error) {
	var user authModel.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *authModel.User) error {
	return r.db.Model(&authModel.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"full_name": user.FullName,
		"phone":     user.Phone,
	}).Error
}

func (r *UserRepository) FindAll(page, limit int) ([]authModel.User, int64, error) {
	var users []authModel.User
	var total int64
	r.db.Model(&authModel.User{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}
