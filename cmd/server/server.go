package main

import (
	"project.com/internal"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := internal.InitServerConfig()

	storage := internal.NewMemStorage()

	r := internal.NewRouter(storage)

	// Запуск HTTP-сервера
	internal.StartServer(serverCfg.Address, r)
}
