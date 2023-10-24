package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Create() *zap.Logger {
	// Настройки логгера
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, _ := config.Build()
	return logger
}
