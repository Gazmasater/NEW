package main

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"project.com/internal/serverin"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	serverin.InitLogger()

	defer logger.Sync()

	serverin.Sugar.Info("Initializing logger...") // Используем sugar для логирования

	// Инициализация конфигурации и хранилища
	config := serverin.InitServerConfig(logger)
	storage := serverin.NewMemStorage()

	serverin.Sugar.Info("Initializing configuration and storage...")
	// Создание контроллера
	controller := serverin.NewHandlerDependencies(storage, logger)

	// Создание маршрутизатора
	r := chi.NewRouter()
	r.Use(serverin.WithLogging)

	serverin.Sugar.Info("Initializing router...")

	// Монтирование главного роутера
	// Монтирование главного роутера с использованием анонимной функции
	r.Mount("/", serverin.WithLogging(controller.Route()))

	serverin.Sugar.Info("Mounting main router...")

	// Запуск сервера
	serverin.StartServer(config.Address, r)

	serverin.Sugar.Info("Server started on address:", config.Address)
}
