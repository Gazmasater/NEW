package server

import (
	"database/sql"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
	"project.com/internal/config"
	"project.com/internal/storage"
)

type app struct {
	Storage *storage.MemStorage
	Logger  *zap.Logger
	Config  *config.ServerConfig
	DB      *sql.DB // Добавлено поле для базы данных
}

//var logger *zap.Logger

func Init(storage *storage.MemStorage, config *config.ServerConfig, db *sql.DB) *app {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger")
	}
	defer logger.Sync()

	// Ваш код инициализации и других полей структуры app

	return &app{
		Storage: storage,
		Logger:  logger,
		Config:  config,
		DB:      db,
	}
}
