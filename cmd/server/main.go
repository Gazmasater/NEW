package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"fmt"
	"log"
	"net/http"

	"time"

	"project.com/internal/app/server"
	"project.com/internal/config"
	"project.com/internal/logger"
	"project.com/internal/storage"
)

func main() {

	serverCfg := config.InitServerConfig()

	// Создание логгера
	appLogger, err := logger.Create()
	if err != nil {
		// Обработка ошибки
		fmt.Printf("Ошибка при создании логгера: %v", err)
		return
	}

	db, err := sql.Open("postgres", serverCfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Ошибка при открытии соединения с базой данных: %v", err)
		return
	}
	defer db.Close()

	mStorage := storage.NewMemStorage()
	app := server.Init(mStorage, serverCfg, db)
	app.SetupDatabase()
	// Создаем HTTP-сервер с настройками
	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: app.Route(),
	}

	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", serverCfg.Address)
	err1 := server.ListenAndServe()
	if err1 != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err1)
	}
	defer appLogger.Sync()

	if serverCfg.StoreInterval > 0 {
		ticker := time.NewTicker(time.Duration(serverCfg.StoreInterval) * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			jsonData := mStorage.GetAllMetricsJSON()
			if jsonData == "" {
				log.Println("Ошибка при получении JSON-представления метрик")
				continue
			}

			println("!!!!!jsonData!!!!!", jsonData)

			if err := app.WriteJSONToFile(serverCfg.FileStoragePath, jsonData); err != nil {
				log.Fatalf("Ошибка при записи в файл: %v", err)
			}

		}

	}

}
