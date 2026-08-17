package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/bees/hindu-ritual-platform/internal/routes"
	"github.com/bees/hindu-ritual-platform/migrations"
	"github.com/bees/hindu-ritual-platform/pkg/configs"
	"github.com/bees/hindu-ritual-platform/pkg/database"
	"github.com/bees/hindu-ritual-platform/pkg/logger"
	redisPkg "github.com/bees/hindu-ritual-platform/pkg/redis"
)

func main() {
	cfg := configs.LoadConfig()

	log := logger.NewLogger(cfg.App.Env)
	log.Info("Starting server",
		zap.String("app", cfg.App.AppName),
		zap.String("env", cfg.App.Env),
		zap.String("version", cfg.App.Version),
	)

	db, err := database.InitDatabase(&cfg.DB)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}
	log.Info("Database connected successfully")

	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	if err := migrations.Run(db); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}
	log.Info("Database migrations completed")

	redisClient, err := redisPkg.InitRedis(&cfg.Redis)
	if err != nil {
		log.Warn("Redis not available, continuing without it", zap.Error(err))
	} else {
		log.Info("Redis connected successfully")
		defer redisPkg.CloseRedis(redisClient)
	}

	router := gin.New()

	routes.SetupRoutes(router, db, redisClient, cfg, log)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Info(fmt.Sprintf("Server listening on port %s", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}
