package internal

import (
	"database/sql"
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

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type HandlerDependencies struct {
	Storage *MemStorage
	Logger  *zap.Logger
	Config  *ServerConfig
	DB      *sql.DB // Добавлено поле для базы данных
}

func NewHandlerDependencies(storage *MemStorage, logger *zap.Logger, config *ServerConfig, db *sql.DB) *HandlerDependencies {
	return &HandlerDependencies{
		Storage: storage,
		Logger:  logger,
		Config:  config,
		DB:      db, // Инициализировано поле для базы данных
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
			println("SaveMetric", ms.counters[metricName])
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
func (ms *MemStorage) GetAllMetrics() []Metrics {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var allMetrics []Metrics
	for name, value := range ms.gauges {
		allMetrics = append(allMetrics, Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}
	for name, delta := range ms.counters {
		allMetrics = append(allMetrics, Metrics{
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

	var allMetrics []Metrics
	for name, value := range ms.gauges {
		allMetrics = append(allMetrics, Metrics{
			ID:    name,
			MType: "gauge",
			Value: &value,
		})
	}
	for name, delta := range ms.counters {
		allMetrics = append(allMetrics, Metrics{
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
