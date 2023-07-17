package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

var pollCount uint64 // Счётчик типа uint64 (беззнаковый int)

type Metric struct {
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func collectMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics["Alloc"] = memStats.Alloc
	// Добавьте остальные метрики памяти, которые вам интересны

	// Увеличиваем счётчик PollCount при каждом вызове collectMetrics
	atomic.AddUint64(&pollCount, 1)

	return metrics
}

func sendMetrics(metrics map[string]interface{}) {
	for name, value := range metrics {
		url := fmt.Sprintf("%s/update/gauge/%s/%v", serverAddress, name, value)
		// Здесь мы форматируем URL, чтобы соответствовать указанному формату данных.

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при выполнении запроса:", err)
			continue
		}
		defer resp.Body.Close()

		// Опционально: можно обработать ответ от сервера, если это необходимо.
		// respBody, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	fmt.Println("Ошибка при чтении ответа:", err)
		// 	continue
		// }
		// fmt.Println("Ответ сервера:", string(respBody))
	}
}

func main() {
	for {
		startTime := time.Now()

		// Собираем метрики
		metrics := collectMetrics()

		// Добавляем PollCount и RandomValue
		metrics["PollCount"] = atomic.LoadUint64(&pollCount)
		metrics["RandomValue"] = 42 // Произвольное значение, замените на свою логику

		// Отправляем метрики на сервер
		sendMetrics(metrics)

		// Ожидаем до следующего опроса и отправки
		sleepTime := reportInterval - time.Since(startTime)
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
}
