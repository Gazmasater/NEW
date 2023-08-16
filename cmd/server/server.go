package main

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"project.com/internal/serverin"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	config := serverin.InitServerConfig(logger)
	storage := serverin.NewMemStorage() // Создание объекта MemStorage

	controller := serverin.NewHandlerDependencies(storage, logger)

	r := chi.NewRouter()
	r.Mount("/", controller.Route()) // Монтирование главного роутера

	serverin.StartServer(config.Address, r) // Запуск сервера с использованием адреса из конфигурации
}
