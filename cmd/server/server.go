package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"project.com/internal"
)

func NewRouter(storage *internal.MemStorage) http.Handler {
	r := chi.NewRouter()

	r.Get("/metrics", internal.HandleMetrics(storage))

	r.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		internal.HandlePostRequest(w, r)
	})

	r.Get("/value/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		internal.HandleGetRequest(w, r, storage)
	})

	return r
}

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
