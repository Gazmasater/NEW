package internal

import (
	"flag"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func ParseAddr() (string, error) {
	// Определение и парсинг флага
	addr := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	flag.Parse()

	return *addr, nil
}

type Metric struct {
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

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

func (ms *MemStorage) PrbocessMetrics(metricType, metricName string, metricValue interface{}) {
	// Сохраняем метрики в хранил.
	ms.SaveMetric(metricType, metricName, metricValue)
}

// GetAllMetrics retrieves all the metr and their values from the storage.
func (ms *MemStorage) GetAllMetrics() map[string]map[string]interface{} {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	allMetrics := make(map[string]map[string]interface{})
	for name, value := range ms.gauges {
		allMetrics[name] = map[string]interface{}{
			"type":  "gauge",
			"value": value,
		}
	}
	for name, value := range ms.counters {
		allMetrics[name] = map[string]interface{}{
			"type":  "counter",
			"value": value,
		}
	}
	return allMetrics
}

func HandleMetrics(storage *MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		allMetrics := storage.GetAllMetrics()

		// Формируем JSON с данными о метриках
		c.JSON(http.StatusOK, allMetrics)
	}
}
