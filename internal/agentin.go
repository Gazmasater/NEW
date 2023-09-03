package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"math/rand"
	"net/http"
	"runtime"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger

type Metrics struct {
	ID    string   `json:"id"`    // имя метрики
	MType string   `json:"type"`  // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value"` // значение метрики в случае передачи gauge
}

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []*Metrics {
	metricsChan := make(chan []*Metrics)
	println("CollectMetrics serverURL string", serverURL)

	var pollCount int64 = 0

	var memStats runtime.MemStats
	go func() {

		for {
			metrics := make([]*Metrics, 0) // Инициализируем срез

			runtime.ReadMemStats(&memStats)

			allocValue := float64(memStats.Alloc)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "Alloc", Value: &allocValue})

			buckHashSysValue := float64(memStats.BuckHashSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "BuckHashSys", Value: &buckHashSysValue})

			freesValue := float64(memStats.Frees)
			freesValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "Frees", Value: &freesValue})

			gCCPUFractionValue := float64(memStats.GCCPUFraction)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "GCCPUFraction", Value: &gCCPUFractionValue})

			gCSysValue := float64(memStats.GCSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "GCSys", Value: &gCSysValue})

			heapAllocValue := float64(memStats.HeapAlloc)
			heapAllocValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapAlloc", Value: &heapAllocValue})

			heapIdleValue := float64(memStats.HeapIdle)
			heapIdleValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapIdle", Value: &heapIdleValue})

			heapInuseValue := float64(memStats.HeapInuse)
			heapInuseValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapInuse", Value: &heapInuseValue})

			heapObjectsValue := float64(memStats.HeapObjects)
			heapObjectsValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapObjects", Value: &heapObjectsValue})

			heapReleasedValue := float64(memStats.HeapReleased)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapReleased", Value: &heapReleasedValue})

			heapSysValue := float64(memStats.HeapSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "HeapSys", Value: &heapSysValue})

			lastGCValue := float64(memStats.LastGC)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "LastGC", Value: &lastGCValue})

			lookupsValue := float64(memStats.Lookups)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "Lookups", Value: &lookupsValue})

			mCacheInuseValue := float64(memStats.MCacheInuse)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "MCacheInuse", Value: &mCacheInuseValue})

			mCacheSysValue := float64(memStats.MCacheSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "MCacheSys", Value: &mCacheSysValue})

			mSpanInuseValue := float64(memStats.MSpanInuse)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "MSpanInuse", Value: &mSpanInuseValue})

			mSpanSysValue := float64(memStats.MSpanSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "MSpanSys", Value: &mSpanSysValue})

			mallocsValue := float64(memStats.Mallocs)
			mallocsValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "Mallocs", Value: &mallocsValue})

			nextGCValue := float64(memStats.NextGC)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "NextGC", Value: &nextGCValue})

			numForcedGCValue := float64(memStats.NumForcedGC)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "NumForcedGC", Value: &numForcedGCValue})

			numGCValue := float64(memStats.NumGC)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "NumGC", Value: &numGCValue})

			otherSysValue := float64(memStats.OtherSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "OtherSys", Value: &otherSysValue})

			pauseTotalNsValue := float64(memStats.PauseTotalNs)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "PauseTotalNs", Value: &pauseTotalNsValue})

			stackInuseValue := float64(memStats.StackInuse)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "StackInuse", Value: &stackInuseValue})

			stackSysValue := float64(memStats.StackSys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "StackSys", Value: &stackSysValue})

			sysValue := float64(memStats.Sys)
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "Sys", Value: &sysValue})

			totalAllocValue := float64(memStats.TotalAlloc)
			totalAllocValue += rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "TotalAlloc", Value: &totalAllocValue})

			// // // Добавляем метрику RandomValue типа gauge с произвольным значением
			randomValue := rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			// // Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			// // Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)
			for _, metric := range metrics {
				var metricValue interface{}
				if metric.MType == "counter" {
					metricValue = *metric.Delta
				} else {
					metricValue = *metric.Value
				}

				fmt.Printf("CollectMetrics!!!!!!! MType: %s, ID: %s, Value: %v\n", metric.MType, metric.ID, metricValue)
			}

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

		data := map[string]interface{}{
			"type":  metric.MType,
			"id":    metric.ID,
			"value": metricValue,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Ошибка при сериализации данных в JSON", err)
			return
		}

		//fmt.Println("Сериализированные данные в JSON:", string(jsonData))
		//	logger.Info("SendDataToServer Сериализированные данные в JSON", zap.String("json_data", string(jsonData)))

		serverURL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverURL, metric.MType, metric.ID, metricValue)
		//	println("SendDataToServer serverURL", serverURL)
		req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		//logger.Info("SendDataToServer Запрос:", zap.Any("request", req))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

		var responseBody []byte
		buf := make([]byte, 1024) // Размер буфера для чтения

		for {
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("Ошибка при чтении тела ответа:", err)
				//	return
			}
			if n == 0 {
				break
			}
			responseBody = append(responseBody, buf[:n]...)
		}

		// Вывод тела ответа на экран
		fmt.Println("SendDataToServer Тело ответа:", string(responseBody))

		if resp.StatusCode == http.StatusOK {
			// Чтение и обработка ответа
			var responseMetrics Metrics

			err := json.Unmarshal(responseBody, &responseMetrics)
			if err != nil {
				fmt.Println("Ошибка при декодировании ответа:", err)
			} else {
				// Обновление значения метрики
				if metric.MType == "counter" {
					*metric.Delta = *responseMetrics.Delta
				} else {
					*metric.Value = *responseMetrics.Value
				}
			}
		} else {
			fmt.Println("Ошибка при отправке запроса. Код статуса:", resp.StatusCode)
		}
	}
}
func SendServerValue(metrics []*Metrics, serverURL string) {
	for _, metric := range metrics {
		data := map[string]interface{}{
			"type": metric.MType,
			"id":   metric.ID,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Ошибка при сериализации данных в JSON", err)
			return
		}

		serverURL := fmt.Sprintf("http://%s/value/%s/%s", serverURL, metric.MType, metric.ID)
		req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		logger.Info("SendServerValue Запрос:", zap.Any("request", req))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

		var responseBody []byte
		buf := make([]byte, 1024) // Размер буфера для чтения

		for {
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("Ошибка при чтении тела ответа:", err)
				return
			}
			if n == 0 {
				break
			}
			responseBody = append(responseBody, buf[:n]...)
		}

		// Вывод тела ответа на экран
		fmt.Println("SendServerValue Тело ответа:", string(responseBody))

		if resp.StatusCode == http.StatusOK {
			// Чтение и обработка ответа
			var responseMetrics Metrics

			err := json.Unmarshal(responseBody, &responseMetrics)
			if err != nil {
				fmt.Println("Ошибка при декодировании ответа:", err)
			} else {
				// Обновление значения метрики
				if metric.MType == "counter" {
					*metric.Delta = *responseMetrics.Delta
				} else {
					*metric.Value = *responseMetrics.Value
				}
			}
			//fmt.Printf("Response Metrics: %+v\n", responseMetrics)

			if metric.MType == "gauge" {
				fmt.Printf("Value: %f\n", *responseMetrics.Value)
			} else if metric.MType == "counter" {
				fmt.Printf("Delta: %d\n", *responseMetrics.Delta)
			}

		} else {
			fmt.Println("Ошибка при отправке запроса. Код статуса:", resp.StatusCode)
		}
	}
}
