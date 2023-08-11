package main

import (
	"project.com/internal/serverin"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := serverin.InitServerConfig()

	storage := serverin.NewMemStorage()
	logger := serverin.NewLogger()

	deps := serverin.NewHandlerDependencies(storage, logger)

	r := serverin.NewRouter(deps)

	serverin.StartServer(serverCfg.Address, r)
}
