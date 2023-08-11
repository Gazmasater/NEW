package serverin

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "[MyApp] ", log.Ldate|log.Ltime)
}

type HandlerDependencies struct {
	Storage *MemStorage
	Logger  *log.Logger
}

func NewHandlerDependencies(storage *MemStorage, logger *log.Logger) *HandlerDependencies {
	return &HandlerDependencies{
		Storage: storage,
		Logger:  logger,
	}
}

func ParseAddr() (string, error) {
	// Определение и парсинг флага
	addr := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	flag.Parse()

	return *addr, nil
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

func HandleMetrics(deps *HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics := deps.Storage.GetAllMetrics()

		// Формируем JSON с данными о метриках
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Используем пакет encoding/json для преобразования данных в JSON и записи их в ResponseWriter.
		json.NewEncoder(w).Encode(allMetrics)
	}
}