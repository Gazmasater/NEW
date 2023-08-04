package main

import (
	"fmt"
	"net/http"

	"time"

	"project.com/internal"
)

func sendDataToServer(metrics []*internal.Metric, serverURL string) {

	for _, metric := range metrics {
		serverURL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverURL, metric.Type, metric.Name, metric.Value)
		//Отправка POST-запроса
		resp, err := http.Post(serverURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

	}
}

func main() {
	config := internal.InitAgentConfig()
	if config == nil {
		return // Если возникли ошибки при инициализации конфигурации, выходим
	}

	// Используем параметры из конфигурации
	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	metricsChan := internal.CollectMetrics(pollInterval, config.Address)

	// Горутина отправки метрик на сервер с интервалом в reportInterval секунд
	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			sendDataToServer(metrics, config.Address)
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
