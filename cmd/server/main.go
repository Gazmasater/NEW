package main

import (
	"github.com/gin-gonic/gin"
	"project.com/internal"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстрации
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))
	r.GET("/metrics", internal.HandleMetrics(storage))

	r.POST("/update/:metricType/:metricName/:metricValue", internal.HandleUpdate(storage))

	// Обработчик для получения всех метрик
	//r.GET("/metrics", func(c *gin.Context) {
	//	// Получаем все известные метрики и их значения
	//	metrics := storage.GetAllMetrics()

	// Формируем JSON-ответ с метриками
	//			c.JSON(http.StatusOK, metrics)
	//		})

	r.GET("/:metricValue/:metricType/:metricName", internal.HandleUpdate(storage))

	r.Run(":8080")
}
