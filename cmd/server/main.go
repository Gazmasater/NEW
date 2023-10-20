package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/go-chi/chi"
	"project.com/internal"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := internal.InitServerConfig()

	// Создание логгера
	logger := internal.CreateLogger()

	r := chi.NewRouter()

	db, err := sql.Open("postgres", serverCfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения с базой данных: %v", err)
		return
	}
	defer db.Close()

	storage := internal.NewMemStorage()
	controller := internal.NewHandlerDependencies(storage, logger, serverCfg, db)
	controller.SetupDatabase()
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
			jsonData := storage.GetAllMetricsJSON()
			if jsonData == "" {
				log.Println("Ошибка при получении JSON-представления метрик")
				continue
			}

			println("!!!!!jsonData!!!!!", jsonData)

			if err := internal.WriteJSONToFile(serverCfg.FileStoragePath, jsonData); err != nil {
				log.Fatalf("Ошибка при записи в файл: %v", err)
			}

		}

	}

}
