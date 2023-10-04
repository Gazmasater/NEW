package main

import (
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/go-chi/chi"
	internal "project.com/internal/server"
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

	if serverCfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(serverCfg.StoreInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metricsToSave := storage.GetAllMetrics()
			for _, metric := range metricsToSave {
				if err := controller.WriteMetricToFile(&metric); err != nil {
					log.Printf("Ошибка при сохранении метрики в файл: %s", err)
				}
			}
		}

	}

}
