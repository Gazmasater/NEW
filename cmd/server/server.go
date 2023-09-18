package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"project.com/internal"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := internal.InitServerConfig()

	// Создание логгера
	logger := internal.CreateLogger()

	r := chi.NewRouter()

	storage := internal.NewMemStorage()
	controller := internal.NewHandlerDependencies(storage, logger, serverCfg)

	r.Route("/", func(r chi.Router) {

		r.Mount("/", controller.Route())
	})

	// Создаем HTTP-сервер с настройками
	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: r,
	}

	// Запуск HTTP-сервера через http.ListenAndServe()
	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", serverCfg.Address)
	err1 := server.ListenAndServe()
	if err1 != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err1)
	}
	defer logger.Sync()

}
