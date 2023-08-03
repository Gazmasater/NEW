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

	r := chi.NewRouter()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))
	storage.SaveMetric("counter", "", int64(10))

	r.Get("/metrics", internal.HandleMetrics(storage))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", internal.HandleUpdate(storage))
	r.Get("/value/{metricType}/{metricName}", internal.HandleUpdate(storage))

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
