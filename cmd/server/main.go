package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"project.com/internal"
)

func main() {

	r := gin.Default()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстр
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))

	// Обработчик для обновления и получения метрик
	r.Any("/", internal.HandleUpdate(storage))

	log.Fatal(http.ListenAndServe(":8080", r))
}
