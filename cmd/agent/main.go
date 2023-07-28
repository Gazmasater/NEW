package main

import (
	"fmt"
	"net/http"
	"time"

	"project.com/internal"
)

func sendDataToServer(metrics []*internal.Metric, serverURL string) {

	for _, metric := range metrics {
		serverURL := fmt.Sprintf("%s/update/%s/%s/%v", serverURL, metric.Type, metric.Name, metric.Value)
		println("serverURL sendDataToServer  ", serverURL)
		// Отправка POST-запроса
		resp, err := http.Post(serverURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

	}
}

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	// Получение адреса сервера с помощью функции GetAddr()
	serverURL := internal.GetAddr()
	println(" serverURL  main", serverURL)

	metricsChan := internal.CollectMetrics(pollInterval, serverURL)

	// Горутина для отправки метрик на сервер с интервалом в 10 секунд
	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			sendDataToServer(metrics, serverURL)
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
