package main

import (
	"github.com/go-chi/chi"
	"project.com/internal/serverin"
)

func main() {
	storage := serverin.NewMemStorage() // Создание объекта MemStorage
	logger := serverin.NewLogger()      // Создание объекта *log.Logger

	deps := serverin.NewHandlerDependencies(storage, logger) // Создание объекта HandlerDependencies с передачей зависимостей

	controller := serverin.NewMyController(deps) // Создание контроллера с переданными зависимостями

	r := chi.NewRouter()
	r.Mount("/", controller.Route()) // Монтирование главного роутера

	serverin.StartServer("localhost:8080", r) // Запуск сервера
}
