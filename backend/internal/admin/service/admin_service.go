package service

import (
	"go.uber.org/zap"

	"github.com/bees/hindu-ritual-platform/internal/admin/dto"
	"github.com/bees/hindu-ritual-platform/internal/admin/repository"
	panditRepo "github.com/bees/hindu-ritual-platform/internal/pandits/repository"
	bookingRepo "github.com/bees/hindu-ritual-platform/internal/bookings/repository"
	auditRepo "github.com/bees/hindu-ritual-platform/internal/audit/repository"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
)

type AdminService struct {
	repo        *repository.AdminRepository
	panditRepo  *panditRepo.PanditRepository
	bookingRepo *bookingRepo.BookingRepository
	auditRepo   *auditRepo.AuditRepository
	log         *logger.Logger
}

func NewAdminService(repo *repository.AdminRepository, panditRepo *panditRepo.PanditRepository, bookingRepo *bookingRepo.BookingRepository, auditRepo *auditRepo.AuditRepository, log *logger.Logger) *AdminService {
	return &AdminService{
		repo:        repo,
		panditRepo:  panditRepo,
		bookingRepo: bookingRepo,
		auditRepo:   auditRepo,
		log:         log,
	}
}

func (s *AdminService) GetDashboard() (*dto.DashboardMetrics, error) {
	metrics, err := s.repo.GetDashboardMetrics()
	if err != nil {
		return nil, err
	}
	return &dto.DashboardMetrics{
		TotalUsers:           metrics.TotalUsers,
		TotalPandits:         metrics.TotalPandits,
		TotalBookings:        metrics.TotalBookings,
		TotalRevenue:         metrics.TotalRevenue,
		PendingVerifications: metrics.PendingVerifications,
		ActiveBookings:       metrics.ActiveBookings,
		FailedLogins:         metrics.FailedLogins,
		NewUsersToday:        metrics.NewUsersToday,
	}, nil
}

func (s *AdminService) SuspendUser(adminID, userID string, reason string) error {
	if err := s.repo.SuspendUser(userID, reason); err != nil {
		return err
	}
	s.log.Info("User suspended", zap.String("admin_id", adminID), zap.String("user_id", userID), zap.String("reason", reason))
	return nil
}

func (s *AdminService) ActivateUser(adminID, userID string) error {
	if err := s.repo.ActivateUser(userID); err != nil {
		return err
	}
	s.log.Info("User activated", zap.String("admin_id", adminID), zap.String("user_id", userID))
	return nil
}
