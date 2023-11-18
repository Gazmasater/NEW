package collector

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"project.com/internal/models"
)

func CollectMetrics_old(pollInterval time.Duration, serverURL string) <-chan []*models.Metrics {
	metricsChan := make(chan []*models.Metrics)
	var pollCount int64 = 0
	println("CollectMetrics serverURL", serverURL)
	var memStats runtime.MemStats

	go func() {

		for {
			metrics := make([]*models.Metrics, 0)

			runtime.ReadMemStats(&memStats)

			addMetric(&metrics, "Alloc", float64(memStats.Alloc))

			addMetric(&metrics, "BuckHashSys", float64(memStats.BuckHashSys))

			addMetric(&metrics, "Frees", float64(memStats.Frees))

			addMetric(&metrics, "GCCPUFraction", float64(memStats.GCCPUFraction))

			addMetric(&metrics, "GCSys", float64(memStats.GCSys))

			addMetric(&metrics, "HeapAlloc", float64(memStats.HeapAlloc))

			addMetric(&metrics, "HeapIdle", float64(memStats.HeapIdle))

			addMetric(&metrics, "HeapInuse", float64(memStats.HeapInuse))

			addMetric(&metrics, "HeapObjects", float64(memStats.HeapObjects))

			addMetric(&metrics, "HeapReleased", float64(memStats.HeapReleased))

			addMetric(&metrics, "HeapSys", float64(memStats.HeapSys))

			addMetric(&metrics, "LastGC", float64(memStats.LastGC))

			addMetric(&metrics, "Lookups", float64(memStats.Lookups))

			addMetric(&metrics, "MCacheInuse", float64(memStats.MCacheInuse))

			addMetric(&metrics, "MCacheSys", float64(memStats.MCacheSys))

			addMetric(&metrics, "MSpanInuse", float64(memStats.MSpanInuse))

			addMetric(&metrics, "MSpanSys", float64(memStats.MSpanSys))

			addMetric(&metrics, "Mallocs", float64(memStats.Mallocs))

			addMetric(&metrics, "NextGC", float64(memStats.NextGC))

			addMetric(&metrics, "NumForcedGC", float64(memStats.NumForcedGC))

			addMetric(&metrics, "NumGC", float64(memStats.NumGC))

			addMetric(&metrics, "OtherSys", float64(memStats.OtherSys))

			addMetric(&metrics, "PauseTotalNs", float64(memStats.PauseTotalNs))

			addMetric(&metrics, "StackInuse", float64(memStats.StackInuse))

			addMetric(&metrics, "StackSys", float64(memStats.StackSys))

			addMetric(&metrics, "Sys", float64(memStats.Sys))

			addMetric(&metrics, "TotalAlloc", float64(memStats.TotalAlloc))

			randomValue := rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &models.Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			//  Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)

		}
	}()

	return metricsChan
}

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []*models.Metrics {
	metricsChan := make(chan []*models.Metrics)
	var pollCount int64 = 0
	println("CollectMetrics serverURL", serverURL)
	var memStats runtime.MemStats

	go func() {
		for {
			metrics := make([]*models.Metrics, 0)

			runtime.ReadMemStats(&memStats)

			for id, getField := range MetricFieldMap {
				addMetric(&metrics, id, getField(&memStats))
			}

			randomValue := rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			metrics = append(metrics, &models.Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

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
		fmt.Println("Ошибка при получении информации о памяти:", err)

		return 0, 0, nil
	}

	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Ошибка при получении информации о CPU:", err)
		return 0, 0, nil
	}

	// Преобразуйте vmStat.Total в float64
	totalMemory := float64(vmStat.Total)

	return totalMemory, float64(vmStat.Free), cpuInfo
}

func addMetric(metrics *[]*models.Metrics, id string, value float64) {
	*metrics = append(*metrics, &models.Metrics{
		MType: "gauge",
		ID:    id,
		Value: &value,
	})
}

type MetricField func(stats *runtime.MemStats) float64

var MetricFieldMap = map[string]MetricField{
	"Alloc":         func(stats *runtime.MemStats) float64 { return float64(stats.Alloc) },
	"BuckHashSys":   func(stats *runtime.MemStats) float64 { return float64(stats.BuckHashSys) },
	"Frees":         func(stats *runtime.MemStats) float64 { return float64(stats.Frees) },
	"GCCPUFraction": func(stats *runtime.MemStats) float64 { return float64(stats.GCCPUFraction) },

	"GCSys":        func(stats *runtime.MemStats) float64 { return float64(stats.GCSys) },
	"HeapAlloc":    func(stats *runtime.MemStats) float64 { return float64(stats.HeapAlloc) },
	"HeapIdle":     func(stats *runtime.MemStats) float64 { return float64(stats.HeapIdle) },
	"HeapInuse":    func(stats *runtime.MemStats) float64 { return float64(stats.HeapInuse) },
	"HeapObjects":  func(stats *runtime.MemStats) float64 { return float64(stats.HeapObjects) },
	"HeapReleased": func(stats *runtime.MemStats) float64 { return float64(stats.HeapReleased) },
	"HeapSys":      func(stats *runtime.MemStats) float64 { return float64(stats.HeapSys) },
	"LastGC":       func(stats *runtime.MemStats) float64 { return float64(stats.LastGC) },
	"Lookups":      func(stats *runtime.MemStats) float64 { return float64(stats.Lookups) },
	"MCacheInuse":  func(stats *runtime.MemStats) float64 { return float64(stats.MCacheInuse) },
	"MCacheSys":    func(stats *runtime.MemStats) float64 { return float64(stats.MCacheSys) },
	"MSpanInuse":   func(stats *runtime.MemStats) float64 { return float64(stats.MSpanInuse) },
	"MSpanSys":     func(stats *runtime.MemStats) float64 { return float64(stats.MSpanSys) },
	"Mallocs":      func(stats *runtime.MemStats) float64 { return float64(stats.Mallocs) },
	"NextGC":       func(stats *runtime.MemStats) float64 { return float64(stats.NextGC) },
	"NumForcedGC":  func(stats *runtime.MemStats) float64 { return float64(stats.NumForcedGC) },
	"NumGC":        func(stats *runtime.MemStats) float64 { return float64(stats.NumGC) },
	"OtherSys":     func(stats *runtime.MemStats) float64 { return float64(stats.OtherSys) },
	"PauseTotalNs": func(stats *runtime.MemStats) float64 { return float64(stats.PauseTotalNs) },
	"StackInuse":   func(stats *runtime.MemStats) float64 { return float64(stats.StackInuse) },
	"StackSys":     func(stats *runtime.MemStats) float64 { return float64(stats.StackSys) },
	"Sys":          func(stats *runtime.MemStats) float64 { return float64(stats.Sys) },
	"TotalAlloc":   func(stats *runtime.MemStats) float64 { return float64(stats.TotalAlloc) },
}
