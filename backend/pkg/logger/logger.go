package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger

var (
	logger *zap.Logger
	once   sync.Once
)

type LogContext struct {
	RequestID string
	UserID    string
	IP        string
	Path      string
	Method    string
}

func NewLogger(env string) *zap.Logger {
	once.Do(func() {
		var zapConfig zap.Config

		if env == "production" {
			zapConfig = zap.NewProductionConfig()
			zapConfig.EncoderConfig.TimeKey = "timestamp"
			zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		} else {
			zapConfig = zap.NewDevelopmentConfig()
			zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		zapConfig.OutputPaths = []string{"stdout"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}

		var err error
		logger, err = zapConfig.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		if err != nil {
			logger = zap.NewNop()
		}

		if env != "production" {
			logger = logger.WithOptions(zap.Development())
		}
	})

	return logger
}

func GetLogger() *zap.Logger {
	if logger == nil {
		return zap.NewNop()
	}
	return logger
}

func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
	os.Exit(1)
}

func InfoWithContext(ctx *LogContext, msg string, extra ...zap.Field) {
	fields := append(
		[]zap.Field{
			zap.String("request_id", ctx.RequestID),
			zap.String("user_id", ctx.UserID),
			zap.String("ip", ctx.IP),
			zap.String("path", ctx.Path),
			zap.String("method", ctx.Method),
		},
		extra...,
	)
	GetLogger().Info(msg, fields...)
}

func ErrorWithContext(ctx *LogContext, msg string, extra ...zap.Field) {
	fields := append(
		[]zap.Field{
			zap.String("request_id", ctx.RequestID),
			zap.String("user_id", ctx.UserID),
			zap.String("ip", ctx.IP),
			zap.String("path", ctx.Path),
			zap.String("method", ctx.Method),
		},
		extra...,
	)
	GetLogger().Error(msg, fields...)
}

func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}
