package main

import (
	"fmt"
	"log"
	"net/http"

	"project.com/internal"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := internal.InitServerConfig()

	storage := internal.NewMemStorage()

	r := NewRouter(storage)

	// Создаем HTTP-сервер с настройками
	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: r,
	}

	// Запуск HTTP-сервера через http.ListenAndServe()
	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", serverCfg.Address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}
