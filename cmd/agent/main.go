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
	config := config.New()
	if config == nil {
		return // Если возникли ошибки при инициализации конфигурации, выходим
	}
	logger.Create()
	// Используем параметры из конфигурации
	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	metricsChan := collector.CollectMetrics(pollInterval, config.Address)

	// Горутина отправки метрик на сервер с интервалом в reportInterval секунд
	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			sender.SendDataToServer(metrics, config.Address) //post запрос по пути /update/
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
