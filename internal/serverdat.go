package internal

import (
	"encoding/json"
	"flag"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ParseAddr() (string, error) {

	addr := flag.String("a", "localhost:8080", "Адрес HTTP-сервера")
	flag.Parse()

	return *addr, nil
}

type HandlerDependencies struct {
	Storage *MemStorage
	Logger  *zap.Logger
}

func NewHandlerDependencies(storage *MemStorage, logger *zap.Logger) *HandlerDependencies {
	return &HandlerDependencies{
		Storage: storage,
		Logger:  logger,
	}
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

func CreateLogger() *zap.Logger {
	// Настройки логгера
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, _ := config.Build()
	return logger
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

func HandleMetrics(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics := storage.GetAllMetrics()
		println("r *http.Request", r)
		// Формируем JSON с данными о метриках
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Используем пакет encoding/json для преобразования данных в JSON и записи их в ResponseWriter.
		json.NewEncoder(w).Encode(allMetrics)
	}
}
