package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"project.com/internal"
)

//func parseAddr() (string, error) {
// Определение и парсинг флага
//	var addr string

//	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")

//	fmt.Println("here is address server", addr)

//	return addr, nil
//}

func main() {

	// Вызыв новую функцию для парсинга флага и получения адреса сервера
	// Вызыв новую функцию для парсинга флага и получения адреса сервера
	//addr, err := parseAddr()
	//if err != nil {
	//	fmt.Println("Ошибка парсинга адреса сервера:", err)
	//	return
	//}
	var addr string

	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")

	flag.Parse()

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

	println("serverURL  main server", addr)

	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", addr)
	if err := r.Run(addr); err != nil {

		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}

}
