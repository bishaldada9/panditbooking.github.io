package repository

import (
	"time"
	"gorm.io/gorm"
	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	panditModel "github.com/bees/hindu-ritual-platform/internal/pandits/model"
	bookingModel "github.com/bees/hindu-ritual-platform/internal/bookings/model"
	paymentModel "github.com/bees/hindu-ritual-platform/internal/payments/model"
	auditModel "github.com/bees/hindu-ritual-platform/internal/audit/model"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetDashboardMetrics() (*DashboardMetrics, error) {
	var metrics DashboardMetrics
	today := time.Now().Truncate(24 * time.Hour)

	r.db.Model(&authModel.User{}).Count(&metrics.TotalUsers)
	r.db.Model(&panditModel.Pandit{}).Count(&metrics.TotalPandits)
	r.db.Model(&bookingModel.Booking{}).Count(&metrics.TotalBookings)
	r.db.Model(&panditModel.Pandit{}).Where("verification_status = ?", "pending").Count(&metrics.PendingVerifications)
	r.db.Model(&bookingModel.Booking{}).Where("status IN ?", []string{"pending", "confirmed"}).Count(&metrics.ActiveBookings)
	r.db.Model(&authModel.User{}).Where("failed_login_attempts > 0").Count(&metrics.FailedLogins)
	r.db.Model(&authModel.User{}).Where("created_at >= ?", today).Count(&metrics.NewUsersToday)
	r.db.Model(&paymentModel.Payment{}).Where("status = ?", "completed").Select("COALESCE(SUM(amount), 0)").Scan(&metrics.TotalRevenue)

	return &metrics, nil
}

func (r *AdminRepository) SuspendUser(userID string, reason string) error {
	return r.db.Model(&authModel.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_suspended": true,
	}).Error
}

func (r *AdminRepository) ActivateUser(userID string) error {
	return r.db.Model(&authModel.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_suspended": false,
	}).Error
}

func (r *AdminRepository) GetAuditLogs(page, limit int) ([]auditModel.AuditLog, int64, error) {
	var logs []auditModel.AuditLog
	var total int64
	r.db.Model(&auditModel.AuditLog{}).Count(&total)
	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

type DashboardMetrics struct {
	TotalUsers           int64
	TotalPandits         int64
	TotalBookings        int64
	TotalRevenue         float64
	PendingVerifications int64
	ActiveBookings       int64
	FailedLogins         int64
	NewUsersToday        int64
}
