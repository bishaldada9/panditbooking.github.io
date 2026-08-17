package routes

import (
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"

	adminHandler "github.com/bees/hindu-ritual-platform/internal/admin/handler"
	adminRepo "github.com/bees/hindu-ritual-platform/internal/admin/repository"
	adminSvc "github.com/bees/hindu-ritual-platform/internal/admin/service"
	auditHandler "github.com/bees/hindu-ritual-platform/internal/audit/handler"
	auditRepo "github.com/bees/hindu-ritual-platform/internal/audit/repository"
	auditSvc "github.com/bees/hindu-ritual-platform/internal/audit/service"
	authHandler "github.com/bees/hindu-ritual-platform/internal/auth/handler"
	authRepo "github.com/bees/hindu-ritual-platform/internal/auth/repository"
	authSvc "github.com/bees/hindu-ritual-platform/internal/auth/service"
	bookingHandler "github.com/bees/hindu-ritual-platform/internal/bookings/handler"
	bookingRepo "github.com/bees/hindu-ritual-platform/internal/bookings/repository"
	bookingSvc "github.com/bees/hindu-ritual-platform/internal/bookings/service"
	"github.com/bees/hindu-ritual-platform/internal/middleware"
	notifHandler "github.com/bees/hindu-ritual-platform/internal/notification/handler"
	notifRepo "github.com/bees/hindu-ritual-platform/internal/notification/repository"
	notifSvc "github.com/bees/hindu-ritual-platform/internal/notification/service"
	panditHandler "github.com/bees/hindu-ritual-platform/internal/pandits/handler"
	panditRepo "github.com/bees/hindu-ritual-platform/internal/pandits/repository"
	panditSvc "github.com/bees/hindu-ritual-platform/internal/pandits/service"
	paymentHandler "github.com/bees/hindu-ritual-platform/internal/payments/handler"
	paymentRepo "github.com/bees/hindu-ritual-platform/internal/payments/repository"
	paymentSvc "github.com/bees/hindu-ritual-platform/internal/payments/service"
	reviewHandler "github.com/bees/hindu-ritual-platform/internal/reviews/handler"
	reviewRepo "github.com/bees/hindu-ritual-platform/internal/reviews/repository"
	reviewSvc "github.com/bees/hindu-ritual-platform/internal/reviews/service"
	ritualHandler "github.com/bees/hindu-ritual-platform/internal/rituals/handler"
	ritualRepo "github.com/bees/hindu-ritual-platform/internal/rituals/repository"
	ritualSvc "github.com/bees/hindu-ritual-platform/internal/rituals/service"
	userHandler "github.com/bees/hindu-ritual-platform/internal/users/handler"
	userRepo "github.com/bees/hindu-ritual-platform/internal/users/repository"
	userSvc "github.com/bees/hindu-ritual-platform/internal/users/service"
	"github.com/bees/hindu-ritual-platform/pkg/configs"
	"github.com/bees/hindu-ritual-platform/pkg/redis"
)

func SetupRoutes(
	router *gin.Engine,
	db *gorm.DB,
	rdb *goredis.Client,
	cfg *configs.Config,
	log *zap.Logger,
) {
	redisClient := redis.NewRedisClient(rdb)

	authRepoInstance := authRepo.NewAuthRepository(db)
	userRepoInstance := userRepo.NewUserRepository(db)
	panditRepoInstance := panditRepo.NewPanditRepository(db)
	bookingRepoInstance := bookingRepo.NewBookingRepository(db)
	ritualRepoInstance := ritualRepo.NewRitualRepository(db)
	paymentRepoInstance := paymentRepo.NewPaymentRepository(db)
	reviewRepoInstance := reviewRepo.NewReviewRepository(db)
	adminRepoInstance := adminRepo.NewAdminRepository(db)
	auditRepoInstance := auditRepo.NewAuditRepository(db)
	notifRepoInstance := notifRepo.NewNotificationRepository(db)

	jwtService := authSvc.NewJWTService(cfg, redisClient)
	mfaService := authSvc.NewMFAService(authRepoInstance, cfg)
	authServiceInstance := authSvc.NewAuthService(authRepoInstance, jwtService, mfaService, log)
	userServiceInstance := userSvc.NewUserService(userRepoInstance)
	panditServiceInstance := panditSvc.NewPanditService(panditRepoInstance)
	ritualServiceInstance := ritualSvc.NewRitualService(ritualRepoInstance)
	bookingServiceInstance := bookingSvc.NewBookingService(bookingRepoInstance, panditRepoInstance, ritualRepoInstance, log)
	paymentServiceInstance := paymentSvc.NewPaymentService(paymentRepoInstance, bookingRepoInstance, cfg, log)
	reviewServiceInstance := reviewSvc.NewReviewService(reviewRepoInstance, bookingRepoInstance, panditRepoInstance)
	adminServiceInstance := adminSvc.NewAdminService(adminRepoInstance, panditRepoInstance, bookingRepoInstance, auditRepoInstance, log)
	auditServiceInstance := auditSvc.NewAuditService(auditRepoInstance)
	notifServiceInstance := notifSvc.NewNotificationService(notifRepoInstance)

	authHandlerInstance := authHandler.NewAuthHandler(authServiceInstance, mfaService, log)
	userHandlerInstance := userHandler.NewUserHandler(userServiceInstance)
	panditHandlerInstance := panditHandler.NewPanditHandler(panditServiceInstance)
	ritualHandlerInstance := ritualHandler.NewRitualHandler(ritualServiceInstance)
	bookingHandlerInstance := bookingHandler.NewBookingHandler(bookingServiceInstance)
	paymentHandlerInstance := paymentHandler.NewPaymentHandler(paymentServiceInstance)
	reviewHandlerInstance := reviewHandler.NewReviewHandler(reviewServiceInstance)
	adminHandlerInstance := adminHandler.NewAdminHandler(adminServiceInstance)
	auditHandlerInstance := auditHandler.NewAuditHandler(auditServiceInstance)
	notifHandlerInstance := notifHandler.NewNotificationHandler(notifServiceInstance)

	authMiddleware := middleware.NewAuthMiddleware(jwtService, redisClient)
	securityMiddleware := middleware.NewSecurityMiddleware(cfg)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.RequestsPerMinute, cfg.RateLimit.BurstSize)

	router.Use(
		securityMiddleware.SecurityHeaders(),
		securityMiddleware.CORSMiddleware(),
		middleware.SQLInjectionDetectionMiddleware(),
		middleware.RateLimitMiddleware(rateLimiter),
		middleware.RequestLogger(log),
		middleware.PanicRecovery(log),
		middleware.RequestIDMiddleware(),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "hindu-ritual-platform", "version": cfg.App.Version})
	})

	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuditTrail(auditServiceInstance))
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandlerInstance.Register)
			auth.POST("/login", authHandlerInstance.Login)
			auth.POST("/refresh", authHandlerInstance.RefreshToken)
			auth.POST("/forgot-password", authHandlerInstance.ForgotPassword)
			auth.POST("/reset-password", authHandlerInstance.ResetPassword)
		}

		v1.GET("/rituals", ritualHandlerInstance.GetRituals)
		v1.GET("/rituals/:id", ritualHandlerInstance.GetRitual)
		v1.GET("/categories", ritualHandlerInstance.GetCategories)
		v1.GET("/pandits", panditHandlerInstance.List)
		v1.GET("/pandits/:id", panditHandlerInstance.GetPandit)
		v1.GET("/pandits/:id/availability", panditHandlerInstance.GetAvailability)
		v1.GET("/pandits/:id/reviews", reviewHandlerInstance.GetPanditReviews)

		webhook := v1.Group("/payments/webhook")
		{
			webhook.POST("/:gateway", paymentHandlerInstance.Webhook)
		}

		authenticated := v1.Group("")
		authenticated.Use(authMiddleware.RequireAuth())
		{
			authenticated.GET("/profile", userHandlerInstance.GetProfile)
			authenticated.PUT("/profile", userHandlerInstance.UpdateProfile)

			authenticated.POST("/auth/logout", authHandlerInstance.Logout)
			authenticated.POST("/auth/change-password", authHandlerInstance.ChangePassword)
			authenticated.POST("/auth/verify-email", authHandlerInstance.VerifyEmail)
			authenticated.GET("/auth/sessions", authHandlerInstance.GetSessions)
			authenticated.POST("/auth/logout-all", authHandlerInstance.LogoutAllSessions)

			authenticated.POST("/auth/mfa/setup", authHandlerInstance.SetupMFA)
			authenticated.POST("/auth/mfa/verify", authHandlerInstance.VerifyMFA)
			authenticated.POST("/auth/mfa/disable", authHandlerInstance.DisableMFA)

			authenticated.POST("/bookings", bookingHandlerInstance.Create)
			authenticated.GET("/bookings", bookingHandlerInstance.ListMyBookings)
			authenticated.GET("/bookings/:id", bookingHandlerInstance.Get)

			authenticated.PUT("/bookings/:id/confirm", bookingHandlerInstance.Confirm)
			authenticated.PUT("/bookings/:id/complete", bookingHandlerInstance.Complete)
			authenticated.PUT("/bookings/:id/cancel", bookingHandlerInstance.Cancel)
			authenticated.PUT("/bookings/:id/reject", bookingHandlerInstance.Reject)

			pandit := authenticated.Group("/pandit")
			pandit.Use(authMiddleware.RequireRole("pandit"))
			{
				pandit.POST("/register", panditHandlerInstance.Register)
				pandit.GET("/profile", panditHandlerInstance.GetProfile)
				pandit.PUT("/profile", panditHandlerInstance.UpdateProfile)
				pandit.POST("/availability", panditHandlerInstance.UpdateAvailability)
				pandit.GET("/bookings", bookingHandlerInstance.ListPanditBookings)
			}

			authenticated.POST("/reviews", reviewHandlerInstance.Create)

			authenticated.POST("/payments/initiate", paymentHandlerInstance.Initiate)
			authenticated.POST("/payments/verify", paymentHandlerInstance.Verify)
			authenticated.GET("/payments/:id", paymentHandlerInstance.GetPayment)
			authenticated.GET("/payments/booking/:bookingId", paymentHandlerInstance.GetPaymentByBooking)

			authenticated.GET("/notifications", notifHandlerInstance.List)
			authenticated.PUT("/notifications/:id/read", notifHandlerInstance.MarkRead)
			authenticated.PUT("/notifications/read-all", notifHandlerInstance.MarkAllRead)
			authenticated.GET("/notifications/unread-count", notifHandlerInstance.UnreadCount)
		}

		admin := v1.Group("/admin")
		admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("admin"))
		{
			admin.GET("/dashboard", adminHandlerInstance.GetDashboard)
			admin.GET("/users", userHandlerInstance.ListUsers)
			admin.GET("/pandits", panditHandlerInstance.ListForAdmin)
			admin.PUT("/users/:id/suspend", adminHandlerInstance.SuspendUser)
			admin.PUT("/users/:id/activate", adminHandlerInstance.ActivateUser)
			admin.PUT("/pandits/:id/verify", panditHandlerInstance.VerifyPandit)
			admin.GET("/bookings", bookingHandlerInstance.ListAll)
			admin.GET("/payments", paymentHandlerInstance.List)
			admin.GET("/reviews", reviewHandlerInstance.List)
			admin.GET("/audit-logs", auditHandlerInstance.GetLogs)
			admin.POST("/reviews/:id/moderate", reviewHandlerInstance.Moderate)
			admin.POST("/payments/refund", paymentHandlerInstance.Refund)
			admin.POST("/categories", ritualHandlerInstance.CreateCategory)
			admin.POST("/rituals", ritualHandlerInstance.CreateRitual)
		}
	}
}
