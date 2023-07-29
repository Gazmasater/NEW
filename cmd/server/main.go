package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"project.com/internal"
)

func main() {

	// Вызываем новую функцию для парсинга флага и получения адреса сервера
	Addr, err := internal.ParseAddr()
	if err != nil {
		fmt.Println("Ошибка парсинга адреса сервера:", err)
		return
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстр
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))
	storage.SaveMetric("counter", "", int64(10))

	r.GET("/metrics", internal.HandleMetrics(storage))

	r.POST("/update/:metricType/:metricName/:metricValue", internal.HandleUpdate(storage))

	r.GET("/:metricValue/:metricType/:metricName", internal.HandleUpdate(storage))

	// Запуск HTTP-сервера на указанном адресе

	println("serverURL  main server", Addr)

	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", *Addr)
	err1 := http.ListenAndServe(*Addr, r)
	if err1 != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}
