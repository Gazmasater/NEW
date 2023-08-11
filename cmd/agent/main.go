package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"project.com/internal/agentin"

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
	config := agentin.InitAgentConfig()
	if config == nil {
		log.Println("Ошибка при инициализации конфигурации")
		return
	}

	// Используем параметры из конфигурации
	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	metricsChan := internal.CollectMetrics(pollInterval, config.Address)

	// Горутина отправки метрик на сервер с интервалом в reportInterval секунд
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		ticker := time.NewTicker(reportInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return // Завершаем горутину при протухании контекста
			case <-ticker.C:
				metrics := <-metricsChan
				sendDataToServer(metrics, config.Address)
			}
		}
	}(ctx)

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
