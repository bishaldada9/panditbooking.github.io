package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/rituals/model"
)

type RitualRepository struct {
	db *gorm.DB
}

func NewRitualRepository(db *gorm.DB) *RitualRepository {
	return &RitualRepository{db: db}
}

func (r *RitualRepository) CreateCategory(category *model.RitualCategory) error {
	return r.db.Create(category).Error
}

func (r *RitualRepository) FindCategoryByID(id uuid.UUID) (*model.RitualCategory, error) {
	var category model.RitualCategory
	err := r.db.First(&category, "id = ?", id).Error
	return &category, err
}

func (r *RitualRepository) FindAllCategories() ([]model.RitualCategory, error) {
	var categories []model.RitualCategory
	err := r.db.Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

func (r *RitualRepository) UpdateCategory(category *model.RitualCategory) error {
	return r.db.Save(category).Error
}

func (r *RitualRepository) CreateRitual(ritual *model.Ritual) error {
	return r.db.Create(ritual).Error
}

func (r *RitualRepository) FindByID(id uuid.UUID) (*model.Ritual, error) {
	var ritual model.Ritual
	err := r.db.Preload("Category").First(&ritual, "id = ?", id).Error
	return &ritual, err
}

func (r *RitualRepository) FindByCategory(categoryID uuid.UUID) ([]model.Ritual, error) {
	var rituals []model.Ritual
	err := r.db.Where("category_id = ? AND is_active = ?", categoryID, true).Find(&rituals).Error
	return rituals, err
}

func (r *RitualRepository) FindAll() ([]model.Ritual, error) {
	var rituals []model.Ritual
	err := r.db.Preload("Category").Where("is_active = ?", true).Find(&rituals).Error
	return rituals, err
}

func (r *RitualRepository) UpdateRitual(ritual *model.Ritual) error {
	return r.db.Save(ritual).Error
}
