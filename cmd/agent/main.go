package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"project.com/internal/agentin"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	config := agentin.InitAgentConfig(logger)
	if config == nil {
		logger.Error("Ошибка при инициализации конфигурации")
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
