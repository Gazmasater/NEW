package main

import (
	"log"
	"net/http"
	"os"

	"project.com/internal"
)

func main() {

	file, err := os.Create("server.log")
	if err != nil {
		log.Fatal("Cannot create log file:", err)
	}
	defer file.Close()

	log.SetOutput(file)

	mux := http.NewServeMux()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстр
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))

	mux.HandleFunc("/update/", internal.HandleUpdate(storage))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
