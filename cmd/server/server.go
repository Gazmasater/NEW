package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

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

	var addr string

	// Чтение переменной окружения или установка значения по умолчанию
	addrEnv := os.Getenv("SERVER_ADDRESS")
	println("addr = addrEnv", addrEnv)

	if addrEnv != "" {
		addr = addrEnv
	} else {
		flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")

		if _, err := url.Parse(addr); err != nil {
			fmt.Printf("Ошибка: неверный формат адреса сервера: %s\n", addr)
			flag.Usage()
			os.Exit(1)
		}

	}

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
