package main

import (
	"log"

	"project.com/internal/serverin"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := serverin.InitServerConfig()
	log.Println("serverCfg", serverCfg)

	storage := serverin.NewMemStorage()
	logger := serverin.NewLogger()

	deps := serverin.NewHandlerDependencies(storage, logger)
	controller := serverin.NewMyController(deps) // Создаем новый контроллер

	rootRouter := controller.Route() // Получаем роутер из контроллера
	serverin.StartServer("localhost:8080", rootRouter)
}
