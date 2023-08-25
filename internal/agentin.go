package internal

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []*Metric {
	metricsChan := make(chan []*Metric)
	println("CollectMetrics serverURL string", serverURL)
	// Переменная для счетчика обновлений метрик
	pollCount := 0

	var metrics []*Metric
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	go func() {
		for {

			metrics = append(metrics, &Metric{Type: "gauge", Name: "Alloc", Value: float64(memStats.Alloc)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "BuckHashSys", Value: float64(memStats.BuckHashSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Frees", Value: float64(memStats.Frees)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "GCCPUFraction", Value: float64(memStats.GCCPUFraction)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "GCSys", Value: float64(memStats.GCSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapAlloc", Value: float64(memStats.HeapAlloc)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapIdle", Value: float64(memStats.HeapIdle)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapInuse", Value: float64(memStats.HeapInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapObjects", Value: float64(memStats.HeapObjects)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapReleased", Value: float64(memStats.HeapReleased)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapSys", Value: float64(memStats.HeapSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "LastGC", Value: float64(memStats.LastGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Lookups", Value: float64(memStats.Lookups)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MCacheInuse", Value: float64(memStats.MCacheInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MCacheSys", Value: float64(memStats.MCacheSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MSpanInuse", Value: float64(memStats.MSpanInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MSpanSys", Value: float64(memStats.MSpanSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Mallocs", Value: float64(memStats.Mallocs)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "NextGC", Value: float64(memStats.NextGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "NumForcedGC", Value: float64(memStats.NumForcedGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "NumGC", Value: float64(memStats.NumGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "OtherSys", Value: float64(memStats.OtherSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "PauseTotalNs", Value: float64(memStats.PauseTotalNs)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "StackInuse", Value: float64(memStats.StackInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "StackSys", Value: float64(memStats.StackSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Sys", Value: float64(memStats.Sys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "TotalAlloc", Value: float64(memStats.TotalAlloc)})

			// Добавляем метрику RandomValue типа gauge с произвольным значением
			randomValue := rand.Float64()
			metrics = append(metrics, &Metric{Type: "gauge", Name: "RandomValue", Value: randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &Metric{Type: "counter", Name: "PollCount", Value: pollCount})

			// Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)
		}
	}()

	return metricsChan
}

func SendDataToServer(metrics []*Metric, serverURL string) {

	for _, metric := range metrics {
		serverURL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverURL, metric.Type, metric.Name, metric.Value)
		//Отправка POST-запроса
		resp, err := http.Post(serverURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

	}
}
