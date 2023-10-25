package main

import (
	"fmt"
	"sync"
	"time"

	"project.com/internal"
)

func main() {
	config := internal.InitAgentConfig()
	if config == nil {
		return
	}
	internal.Init()

	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	metricsChan := internal.CollectMetrics(pollInterval, config.Address)

	var mu sync.Mutex
	var bufferedMetrics []*internal.Metrics

	immediateMu := sync.Mutex{} // Используем {} для создания экземпляра мьютекса

	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan

			// Блокируем доступ к bufferedMetrics при немедленной отправке
			immediateMu.Lock()
			internal.SendDataToServer(metrics, config.Address)
			immediateMu.Unlock()
		}
	}()

	go func() {
		for range time.Tick(reportInterval) {
			mu.Lock()
			metricsToSend := make([]*internal.Metrics, len(bufferedMetrics))
			copy(metricsToSend, bufferedMetrics)
			bufferedMetrics = nil
			mu.Unlock()

			if len(metricsToSend) > 0 { // Проверка на пустой срез

				err := internal.SendDataToServerBatch(metricsToSend, config.Address)
				if err != nil {
					fmt.Println("Ошибка при отправке метрик на сервер:", err)
				}
			}
		}
	}()

	for {
		select {
		case metrics := <-metricsChan:
			fmt.Println("Сбор метрик...")

			mu.Lock()
			bufferedMetrics = append(bufferedMetrics, metrics...)
			mu.Unlock()
		case <-time.After(pollInterval):
			// Выполняйте другие задачи, если не получены новые метрики в течение pollInterval
		}
	}
}
