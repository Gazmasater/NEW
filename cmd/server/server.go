package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"project.com/internal"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := internal.InitServerConfig()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстрации
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))
	storage.SaveMetric("counter", "", int64(10))

	r.GET("/metrics", internal.HandleMetrics(storage))
	r.POST("/update/:metricType/:metricName/:metricValue", internal.HandleUpdate(storage))
	r.GET("/value/:metricType/:metricName", internal.HandleUpdate(storage))

	// Запуск HTTP-сервера на указанном адресе
	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", serverCfg.Address)
	if err := r.Run(serverCfg.Address); err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}
