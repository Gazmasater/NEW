package collector

import (
	"log"

	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"project.com/internal/models"
)

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []models.Metrics {
	metricsChan := make(chan []models.Metrics)
	var pollCount int64 = 0
	println("CollectMetrics serverURL", serverURL)
	var memStats runtime.MemStats

	go func() {
		for {
			metrics := make([]models.Metrics, 0)

			runtime.ReadMemStats(&memStats)

			for id, getField := range MetricFieldMap {
				addMetric(&metrics, id, getField(&memStats))
			}

			metrics = append(metrics, models.Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)
		}
	}()

	return metricsChan
}

func CollectAdditionalMetrics() (float64, float64, []float64) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Ошибка при получении информации о памяти:", err)

		return 0, 0, nil
	}

	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Println("Ошибка при получении информации о CPU:", err)
		return 0, 0, nil
	}

	totalMemory := float64(vmStat.Total)

	return totalMemory, float64(vmStat.Free), cpuInfo
}

func addMetric(metrics *[]models.Metrics, id string, value float64) {
	newValue := &value // Получение указателя на значение float64
	*metrics = append(*metrics, models.Metrics{
		MType: "gauge",
		ID:    id,
		Value: newValue,
	})
}
