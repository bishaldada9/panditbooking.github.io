package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/bees/hindu-ritual-platform/internal/pandits/model"
	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
)

type PanditRepository struct {
	db *gorm.DB
}

func NewPanditRepository(db *gorm.DB) *PanditRepository {
	return &PanditRepository{db: db}
}

func (r *PanditRepository) Create(pandit *model.Pandit) error {
	return r.db.Create(pandit).Error
}

func (r *PanditRepository) FindByUserID(userID uuid.UUID) (*model.Pandit, error) {
	var pandit model.Pandit
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&pandit).Error
	if err != nil {
		return nil, err
	}
	return &pandit, nil
}

func (r *PanditRepository) FindByID(id uuid.UUID) (*model.Pandit, error) {
	var pandit model.Pandit
	err := r.db.Preload("User").First(&pandit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pandit, nil
}

func (r *PanditRepository) FindAll(page, limit int, filters map[string]interface{}) ([]model.Pandit, int64, error) {
	var pandits []model.Pandit
	var total int64
	query := r.db.Model(&model.Pandit{}).Preload("User")
	for key, val := range filters {
		query = query.Where(key, val)
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("rating DESC, total_bookings DESC").Find(&pandits).Error
	return pandits, total, err
}

func (r *PanditRepository) Update(pandit *model.Pandit) error {
	return r.db.Save(pandit).Error
}

func (r *PanditRepository) CreateDocument(doc *model.PanditDocument) error {
	return r.db.Create(doc).Error
}

func (r *PanditRepository) FindDocumentsByPanditID(panditID uuid.UUID) ([]model.PanditDocument, error) {
	var docs []model.PanditDocument
	err := r.db.Where("pandit_id = ?", panditID).Find(&docs).Error
	return docs, err
}

func (r *PanditRepository) UpdateDocument(doc *model.PanditDocument) error {
	return r.db.Save(doc).Error
}

func (r *PanditRepository) CreateAvailability(avail *model.Availability) error {
	return r.db.Create(avail).Error
}

func (r *PanditRepository) FindAvailability(panditID uuid.UUID, date string) ([]model.Availability, error) {
	var availabilities []model.Availability
	err := r.db.Where("pandit_id = ? AND date = ?", panditID, date).Order("start_time").Find(&availabilities).Error
	return availabilities, err
}

func (r *PanditRepository) DeleteAvailability(id uuid.UUID) error {
	return r.db.Delete(&model.Availability{}, "id = ?", id).Error
}

func (r *PanditRepository) FindPendingVerifications() ([]model.Pandit, error) {
	var pandits []model.Pandit
	err := r.db.Where("verification_status = ?", model.VerificationPending).Preload("User").Find(&pandits).Error
	return pandits, err
}

func (r *PanditRepository) GetUserByID(id uuid.UUID) (*authModel.User, error) {
	var user authModel.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}
