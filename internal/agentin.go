package internal

import (
	"go.uber.org/zap"

	"fmt"
	"net/http"
	"runtime"
	"time"
)

var logger *zap.Logger

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []*Metrics {
	metricsChan := make(chan []*Metrics)
	println("CollectMetrics serverURL string", serverURL)
	// Переменная для счетчика обновлений метрик
	var pollCount int64 = 0
	var metrics []*Metrics
	var memStats runtime.MemStats

	go func() {
		for {
			runtime.ReadMemStats(&memStats)
			allocValue := float64(memStats.Alloc)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "Alloc", Value: &allocValue})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "BuckHashSys", Value: float64(memStats.BuckHashSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "Frees", Value: float64(memStats.Frees)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "GCCPUFraction", Value: float64(memStats.GCCPUFraction)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "GCSys", Value: float64(memStats.GCSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapAlloc", Value: float64(memStats.HeapAlloc)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapIdle", Value: float64(memStats.HeapIdle)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapInuse", Value: float64(memStats.HeapInuse)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapObjects", Value: float64(memStats.HeapObjects)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapReleased", Value: float64(memStats.HeapReleased)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapSys", Value: float64(memStats.HeapSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "LastGC", Value: float64(memStats.LastGC)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "Lookups", Value: float64(memStats.Lookups)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "MCacheInuse", Value: float64(memStats.MCacheInuse)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "MCacheSys", Value: float64(memStats.MCacheSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "MSpanInuse", Value: float64(memStats.MSpanInuse)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "MSpanSys", Value: float64(memStats.MSpanSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "Mallocs", Value: float64(memStats.Mallocs)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "NextGC", Value: float64(memStats.NextGC)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "NumForcedGC", Value: float64(memStats.NumForcedGC)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "NumGC", Value: float64(memStats.NumGC)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "OtherSys", Value: float64(memStats.OtherSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "PauseTotalNs", Value: float64(memStats.PauseTotalNs)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "StackInuse", Value: float64(memStats.StackInuse)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "StackSys", Value: float64(memStats.StackSys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "Sys", Value: float64(memStats.Sys)})
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "TotalAlloc", Value: float64(memStats.TotalAlloc)})

			// // Добавляем метрику RandomValue типа gauge с произвольным значением
			// randomValue := rand.Float64()
			// metrics = append(metrics, &Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			// Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)
		}
	}()

	return metricsChan
}

func SendDataToServer(metrics []*Metrics, serverURL string) {

	for _, metric := range metrics {
		var metricValue interface{}
		if metric.MType == "counter" {
			metricValue = *metric.Delta
		} else {
			metricValue = *metric.Value
		}

		serverURL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverURL, metric.MType, metric.ID, metricValue)
		logger.Info("Sending metric",
			zap.String("url", serverURL),
			zap.String("metric_type", metric.MType),
			zap.String("metric_id", metric.ID),
			zap.Any("metric_value", metricValue),
		)

		resp, err := http.Post(serverURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

	}
}
