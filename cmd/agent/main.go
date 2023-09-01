package main

import (
	"fmt"

	"time"

	"project.com/internal"
)

func main() {
	config := internal.InitAgentConfig()
	if config == nil {
		return // Если возникли ошибки при инициализации конфигурации, выходим
	}
	internal.Init()
	// Используем параметры из конфигурации
	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	metricsChan := internal.CollectMetrics(pollInterval, config.Address)

	// Горутина отправки метрик на сервер с интервалом в reportInterval секунд
	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			internal.SendDataToServer(metrics, config.Address)
			internal.SendServerValue(metrics, config.Address)
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
