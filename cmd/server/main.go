package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"project.com/internal"
)

func main() {

	// Определение и инициализация флага -a с значением по умолчанию "localhost:8080"
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	storage := internal.NewMemStorage()

	// Пример сохранения метрик для демонстрации
	storage.SaveMetric("gauge", "temperature", 25.0)
	storage.SaveMetric("counter", "requests", int64(10))
	storage.SaveMetric("counter", "", int64(10))

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

	// Запуск HTTP-сервера на указанном адресе
	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", *internal.Addr)
	err := http.ListenAndServe(*internal.Addr, r)
	if err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}
