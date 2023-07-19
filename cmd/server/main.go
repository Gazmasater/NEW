package main

import (
	"log"
	"net/http"

	"project.com/internal"
)

func main() {

	mux := http.NewServeMux()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстр
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))

	mux.HandleFunc("/", internal.HandleUpdate(storage))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
