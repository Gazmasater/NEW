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

			// // // // Добавляем метрику RandomValue типа gauge с произвольным значением
			randomValue := rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			//  Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)

		}
	}()

	return metricsChan
}

func CollectMetricsJSON(pollInterval time.Duration, serverURL string) <-chan []*Metrics {
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

			// // // // Добавляем метрику RandomValue типа gauge с произвольным значением
			randomValue := rand.Float64()
			metrics = append(metrics, &Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			//  Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)

		}
	}()

	return metricsChan
}

func SendDataToServer(metrics []*Metrics, serverURL string) error {
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
			return fmt.Errorf("ошибка при сериализации данных в JSON:%w", err)
		}

		fmt.Println("Сериализированные данные в JSON:", string(jsonData))
		logger.Info("SendDataToServer Сериализированные данные в JSON", zap.String("json_data", string(jsonData)))

		serverURL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverURL, metric.MType, metric.ID, metricValue)
		//	println("SendDataToServer serverURL", serverURL)
		req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			return fmt.Errorf("ошибка при создании запроса:%w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		//logger.Info("SendDataToServer Запрос:", zap.Any("request", req))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return fmt.Errorf("ошибка при отправке запроса:%w", err)
		}
		defer resp.Body.Close()

		var responseBody []byte
		buf := make([]byte, 1024) // Размер буфера для чтения

		for {
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				return fmt.Errorf("ошибка при чтении тела ответа:%w", err)
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
				return fmt.Errorf("ошибка при декодировании ответа:%w", err)
			} else {
				// Обновление значения метрики
				if metric.MType == "counter" {
					*metric.Delta = *responseMetrics.Delta
				} else {
					*metric.Value = *responseMetrics.Value
				}
			}
		} else {
			return fmt.Errorf("ошибка при отправке запроса. Код статуса:%w", err)
		}
	}
	return nil
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

		} else {
			fmt.Println("Ошибка при отправке запроса. Код статуса:", resp.StatusCode)
		}
	}
}

func SendMetricsJSONToServer(url string, data []byte) error {
	println("!!!!!!!!!!SendMetricsJsonToServer  url ", url)
	// Распаковываем JSON-тело запроса в структуру Metrics
	var metricData Metrics
	if err := json.Unmarshal(data, &metricData); err != nil {
		return err
	}

	// fmt.Printf("SendMetricsJsonToServer ID: %s\n", metricData.ID)
	// fmt.Printf("SendMetricsJsonToServer Type: %s\n", metricData.MType)

	if metricData.Delta != nil {
		fmt.Printf("Delta: %d\n", *metricData.Delta)
		//	url = fmt.Sprintf("%s/%s/%s/%v", url, metricData.MType, metricData.ID, *metricData.Delta)

	} else {
		fmt.Println("Delta: nil")
	}

	if metricData.Value != nil {
		fmt.Printf("Value: %f\n", *metricData.Value)
		//	url = fmt.Sprintf("%s/%s/%s/%v", url, metricData.MType, metricData.ID, *metricData.Value)

	} else {
		fmt.Println("Value: nil")
	}

	// Используем значения из JSON-тела для создания URL

	println("????????????????????????sendMetricsJsonToServer url", url)

	//
	// Создаем запрос POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Устанавливаем заголовки (если необходимо)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неправильный статус код: %d", resp.StatusCode)
	}

	return nil
}

func SendMetricsJSONToServerValue(url string, data []byte) error {

	var metricData Metrics
	if err := json.Unmarshal(data, &metricData); err != nil {
		return err
	}

	if metricData.Delta != nil {
		fmt.Printf("SendMetricsJsonToServerValue Delta: %d\n", *metricData.Delta)
		//	url = fmt.Sprintf("%s/%s/%s/%v", url, metricData.MType, metricData.ID, *metricData.Delta)

	} else {
		fmt.Println("SendMetricsJsonToServerValue Delta: nil")
	}

	println("################sendMetricsJsonToServerValue", url)
	// Создаем запрос POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	fmt.Println("sendMetricsJsonToServerValue req", req)

	// Устанавливаем заголовки (Content-Type: application/json)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неправильный статус код: %d", resp.StatusCode)
	}

	// Читаем JSON-ответ с заполненными значениями метрик
	var metricResponse Metrics
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&metricResponse); err != nil {
		return err
	}

	fmt.Printf("Received Metrics:\n")
	fmt.Printf("ID: %s\n", metricResponse.ID)
	fmt.Printf("Type: %s\n", metricResponse.MType)
	if metricResponse.Delta != nil {
		fmt.Printf("Delta: %d\n", *metricResponse.Delta)
	} else {
		fmt.Println("Delta: nil")
	}
	if metricResponse.Value != nil {
		fmt.Printf("Value: %f\n", *metricResponse.Value)
	} else {
		fmt.Println("Value: nil")
	}

	return nil
}
