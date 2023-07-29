package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"project.com/internal"
)

func main() {

	// Вызыв новую функцию для парсинга флага и получения адреса сервера
	addr, err := internal.ParseAddr() //   internal/serverdat.go
	if err != nil {
		fmt.Println("Ошибка парсинга адреса сервера:", err)
		return
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстр
	// Пример сохранения метрик для демонстр
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))
	storage.SaveMetric("counter", "", int64(10))

	r.GET("/metrics", internal.HandleMetrics(storage))

	r.POST("/update/:metricType/:metricName/:metricValue", internal.HandleUpdate(storage))

	r.GET("/value/:metricType/:metricName", internal.HandleUpdate(storage))

	// Запуск HTTP-сервера на указанном адресе

	println("serverURL  main server", *addr)

	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", *addr)
	err1 := http.ListenAndServe(*addr, r)
	if err1 != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}
