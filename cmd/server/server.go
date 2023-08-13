package main

import (
	"log"

	"github.com/go-chi/chi"
	"project.com/internal/serverin"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := serverin.InitServerConfig()
	log.Println("main serverCfg", serverCfg)

	storage := serverin.NewMemStorage()
	logger := serverin.NewLogger()

	deps := serverin.NewHandlerDependencies(storage, logger)
	controller := serverin.NewMyController(deps) // Создаем новый контроллер

	rootRouter := chi.NewRouter()                      // Создаем корневой роутер
	rootRouter.Mount("/", controller.Route())          // Подключаем роутер из контроллера к корневому роутеру
	serverin.StartServer("localhost:8080", rootRouter) // Запускаем сервер
}
