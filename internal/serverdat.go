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

func (s *MemStorage) SaveMetric(metricType, metricName string, metricValue interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch metricType {
	case "gauge":
		if v, ok := metricValue.(float64); ok {
			s.gauges[metricName] = v
		}
	case "counter":
		if v, ok := metricValue.(int64); ok {
			s.counters[metricName] += v
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
