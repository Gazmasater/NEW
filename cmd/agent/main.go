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
			bufferedMetrics := <-metricsChan
			err := internal.SendDataToServer(bufferedMetrics, config.Address)
			if err != nil {
				fmt.Println("Ошибка при отправке метрик на сервер:", err)
			}
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
