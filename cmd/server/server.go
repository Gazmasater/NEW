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
	log.Println("main storage logger ", storage, logger)

	deps := &serverin.HandlerDependencies{} // Создайте свои зависимости
	controller := serverin.NewMyController(deps)
	router := controller.Route()

	serverin.StartServer("localhost:8080", router)

}
