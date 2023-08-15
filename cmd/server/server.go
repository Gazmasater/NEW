package main

import (
	"github.com/go-chi/chi"
	"project.com/internal/serverin"
)

func main() {
	config := serverin.InitServerConfig()
	storage := serverin.NewMemStorage() // Создание объекта MemStorage
	logger := serverin.NewLogger()      // Создание объекта *log.Logger

	controller := serverin.NewHandlerDependencies(storage, logger)

	r := chi.NewRouter()
	r.Mount("/", controller.Route()) // Монтирование главного роутера

	serverin.StartServer(config.Address, r) // Запуск сервера с использованием адреса из конфигурации
}
