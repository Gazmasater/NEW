package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CustomLogger - пользовательский тип, включающий в себя zap.Logger
type CustomLogger struct {
	*zap.Logger
}

func Create() (*CustomLogger, error) {
	// Настройки логгера
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("не может создать логгер: %w", err)
	}
	return &CustomLogger{logger}, nil
}

func (l *CustomLogger) Info(tmpl string, args ...any) {
	l.Sugar().Infof(tmpl, args...)
}

// Sync - метод для синхронизации логгер
func (l *CustomLogger) Sync() error {
	if l == nil {
		return fmt.Errorf("логгер не был создан")
	}
	return l.Sync()
}
