package main

import (
	"project.com/internal"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := internal.InitServerConfig()

	storage := internal.NewMemStorage()
	logger := internal.NewLogger() // Инициализируйте логгер, если он есть

	deps := internal.NewHandlerDependencies(storage, logger)

	r := internal.NewRouter(deps)

	// Запуск HTTP-сервера
	internal.StartServer(serverCfg.Address, r)
}
