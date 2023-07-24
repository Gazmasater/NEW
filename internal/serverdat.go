package internal

import "sync"

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

func (ms *MemStorage) SaveMetric(metricType, metricName string, metricValue interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	switch metricType {

	case "gauge":
		if v, ok := metricValue.(float64); ok {
			ms.gauges[metricName] = v
		}
	case "counter":
		if v, ok := metricValue.(int64); ok {
			ms.counters[metricName] += v
		}
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

func (ms *MemStorage) ProcessMetrics(metricType, metricName string, metricValue interface{}) {
	// Сохраняем метрику в хранилище
	ms.SaveMetric(metricType, metricName, metricValue)
}

// GetAllMetrics retrieves all the metrics and their values from the storage.
func (ms *MemStorage) GetAllMetrics() map[string]interface{} {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	allMetrics := make(map[string]interface{})
	for name, value := range ms.gauges {
		allMetrics[name] = value
	}
	for name, value := range ms.counters {
		allMetrics[name] = value
	}
	return allMetrics
}
