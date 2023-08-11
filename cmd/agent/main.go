package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"project.com/internal/agentin"
)

func main() {
	config := agentin.InitAgentConfig()
	if config == nil {
		log.Println("Ошибка при инициализации конфигурации")
		return
	}

	// Используем параметры из конфигурации
	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	metricsChan := agentin.CollectMetrics(pollInterval, config.Address)

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
				agentin.SendDataToServer(metrics, config.Address)
			}
		}
	}(ctx)

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
