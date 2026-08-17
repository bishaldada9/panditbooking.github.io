package migrations

import (
	"fmt"

	"gorm.io/gorm"

	authModel "github.com/bees/hindu-ritual-platform/internal/auth/model"
	auditModel "github.com/bees/hindu-ritual-platform/internal/audit/model"
	bookingModel "github.com/bees/hindu-ritual-platform/internal/bookings/model"
	notifModel "github.com/bees/hindu-ritual-platform/internal/notification/model"
	panditModel "github.com/bees/hindu-ritual-platform/internal/pandits/model"
	paymentModel "github.com/bees/hindu-ritual-platform/internal/payments/model"
	reviewModel "github.com/bees/hindu-ritual-platform/internal/reviews/model"
	ritualModel "github.com/bees/hindu-ritual-platform/internal/rituals/model"
)

func Run(db *gorm.DB) error {
	err := db.AutoMigrate(
		&authModel.User{},
		&authModel.RefreshToken{},
		&authModel.OTP{},
		&panditModel.Pandit{},
		&panditModel.PanditDocument{},
		&panditModel.Availability{},
		&ritualModel.RitualCategory{},
		&ritualModel.Ritual{},
		&bookingModel.Booking{},
		&bookingModel.Complaint{},
		&paymentModel.Payment{},
		&paymentModel.Transaction{},
		&reviewModel.Review{},
		&auditModel.AuditLog{},
		&notifModel.Notification{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create indexes for common queries
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func createIndexes(db *gorm.DB) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_customer ON bookings(customer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_pandit ON bookings(pandit_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status)`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_scheduled_date ON bookings(scheduled_date)`,
		`CREATE INDEX IF NOT EXISTS idx_payments_booking ON payments(booking_id)`,
		`CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status)`,
		`CREATE INDEX IF NOT EXISTS idx_payments_transaction ON payments(transaction_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read)`,
		`CREATE INDEX IF NOT EXISTS idx_pandits_verification ON pandits(verification_status)`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_pandit ON reviews(pandit_id)`,
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			return fmt.Errorf("failed to create index: %s, error: %w", idx, err)
		}
	}
	return nil
}
