package main

import (
	"fmt"
	"net/http"
	"time"

	"project.com/internal"
)

func sendDataToServer(metrics []*internal.Metric) {

	for _, metric := range metrics {
		serverURL := fmt.Sprintf("http://*internal.Addr/update/%s/%s/%v", metric.Type, metric.Name, metric.Value)
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
	serverURL := "http://*internal.Addr/update/gauge/test1/100"

	metricsChan := internal.CollectMetrics(pollInterval, serverURL)

	// Горутина для отправки метрик на сервер с интервалом в 10 секунд
	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			sendDataToServer(metrics)
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
