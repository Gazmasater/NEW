package main

import (
	"fmt"

	"time"

	"project.com/internal/collector"
	"project.com/internal/config"
	"project.com/internal/logger"
	"project.com/internal/models"
	"project.com/internal/sender"
)

func main() {
	cfg := config.Must(config.New())

	log, err := logger.New()
	if err != nil {

		fmt.Printf("Ошибка при создании логгера: %s\n", err)
		return
	}

	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	metricsChan := collector.CollectMetrics(pollInterval, cfg.Address)
	metricsChanTotal := make(chan models.SysMetrics)

	requestChannel := make(chan struct{}, cfg.RateLimit)

	// Горутина отправки метрик на сервер
	go func() {
		for range time.Tick(reportInterval) {
			select {
			case metrics := <-metricsChan:

				select {
				case requestChannel <- struct{}{}:
					go func() {
						sender.SendDataToServer(metrics, cfg.Address) // POST-запрос по пути /update/
						<-requestChannel
					}()
				default:

					log.Info("Превышен лимит одновременных запросов, метрика будет отправлена в следующем цикле.")
				}
			case metrics := <-metricsChanTotal:
				sender.SendSysMetricsToServer(metrics, cfg.Address)
			}
		}
	}()

	go func() {
		for range time.Tick(pollInterval) {
			totalMemory, freeMemory, cpuUtilization := collector.CollectAdditionalMetrics()
			metrics := models.SysMetrics{
				TotalMemory:    totalMemory,
				FreeMemory:     freeMemory,
				CPUUtilization: cpuUtilization,
			}

			metricsChanTotal <- metrics
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
