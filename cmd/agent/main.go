package main

import (
	"fmt"

	"time"

	"project.com/internal/collector"
	"project.com/internal/config"
	"project.com/internal/logger"
	"project.com/internal/sender"
)

func main() {
	cfg := config.Must(config.New())

	collector.Init()

	_, err := logger.New()
	if err != nil {

		fmt.Printf("Ошибка при создании логгера: %s\n", err)
		return
	}
	// Используем параметры из конфигурации
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	metricsChan := collector.CollectMetrics(pollInterval, cfg.Address)

	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			sender.SendDataToServer(metrics, cfg.Address) //post запрос по пути /update/
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
