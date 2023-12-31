package storage

import (
	"encoding/json"
	"sync"

	"project.com/internal/models"
)

type MemStorage struct {
	mu       sync.RWMutex
	counters map[string]int64
	gauges   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (ms *MemStorage) GetCounters() map[string]int64 {
	return ms.counters
}

func (ms *MemStorage) GetGauges() map[string]float64 {
	return ms.gauges
}

func (ms *MemStorage) SaveGauge(metricType, metricName string, metricValue float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.gauges[metricName] = metricValue

}

func (ms *MemStorage) SaveCounter(metricType, metricName string, metricValue int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.counters[metricName] += metricValue

}
func (ms *MemStorage) SaveMetric(metricType, metricName string, metricValue any) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	switch v := metricValue.(type) {
	case float64:
		// Если metricValue - float64, сохраняем как Gauge
		ms.gauges[metricName] = v
	case int64:
		// Если metricValue - int64, сохраняем как Counter
		ms.counters[metricName] += v

	}
}

func (ms *MemStorage) GetMetric(metricType, metricName string) (interface{}, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	switch metricType {
	case "gauge":
		value, ok := ms.gauges[metricName]
		return value, ok
	case "counter":
		value, ok := ms.counters[metricName]
		return value, ok
	default:
		return nil, false
	}
}

func (ms *MemStorage) GetAllMetrics() []models.Metrics {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var allMetrics []models.Metrics
	for name, value := range ms.gauges {
		allMetrics = append(allMetrics, models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}
	for name, delta := range ms.counters {
		allMetrics = append(allMetrics, models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		})
	}
	return allMetrics
}

func (ms *MemStorage) GetAllMetricsJSON() string {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var allMetrics []models.Metrics
	for name, value := range ms.gauges {
		allMetrics = append(allMetrics, models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}
	for name, delta := range ms.counters {
		allMetrics = append(allMetrics, models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		})
	}

	// Преобразуем слайс метрик в JSON
	jsonData, err := json.Marshal(allMetrics)
	if err != nil {
		// Обработка ошибки, если не удалось преобразовать в JSON
		return ""
	}

	// Преобразуем []byte в строку с помощью string()
	return string(jsonData)
}
