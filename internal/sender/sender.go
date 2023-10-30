package sender

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"project.com/internal/logger"
	"project.com/internal/models"
	//"project.com/internal/sender"
)

func SendDataToServer(metrics []*models.Metrics, serverURL string) error {
	for _, metric := range metrics {
		var metricValue any
		if metric.MType == "counter" {
			metricValue = *metric.Delta
		} else {
			metricValue = *metric.Value
		}

		data := map[string]any{
			"type": metric.MType,
			"id":   metric.ID,
		}

		if metric.MType == "gauge" {
			data["value"] = metricValue
		} else if metric.MType == "counter" {
			data["delta"] = metricValue
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("ошибка при сериализации данных в JSON:%w", err)
		}

		// Вычисление хеша данных с использованием ключа
		hash := ComputeHash(jsonData, "MyKey")
		log, err := logger.Create()
		if err != nil {
			// Обработка ошибки
			return fmt.Errorf("ошибка при создании логгера: %v", err)
		}

		// Теперь вы можете использовать переменную log (ваш логгер)
		log.Info("Это информационное сообщение")

		serverURL := fmt.Sprintf("http://%s/update/", serverURL)
		req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("ошибка при создании запроса:%w", err)
		}

		// Добавление хеша в заголовок запроса
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("HashSHA256", hash) // Добавление хеша в заголовок запроса

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
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

		if resp.StatusCode == http.StatusOK {
			// Чтение и обработка ответа
			var responseMetrics models.Metrics

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
			return fmt.Errorf("ошибка при отправке запроса. Код статуса: %d", resp.StatusCode)
		}
	}
	return nil
}

func ComputeHash(data []byte, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(data)
	hash := mac.Sum(nil)
	return hex.EncodeToString(hash)
}

func SendDataToServerBatch(metrics []*models.Metrics, serverURL string) error {
	// Создаем общий словарь data для всех метрик
	data := make([]map[string]interface{}, len(metrics))

	for i, metric := range metrics {
		metricData := make(map[string]interface{})
		metricData["id"] = metric.ID
		metricData["type"] = metric.MType

		if metric.MType == "counter" {
			metricData["delta"] = *metric.Delta
		} else {
			metricData["value"] = *metric.Value
		}

		data[i] = metricData
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации данных в JSON:%w", err)
	}

	log := logger.CustomLogger{}
	log.Info("BATCH SendDataToServer Сериализированные данные в JSON", zap.String("json_data", string(jsonData)))

	var gzippedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzippedData)
	_, err = gzipWriter.Write(jsonData)
	if err != nil {
		return fmt.Errorf("ошибка при сжатии данных:%w", err)
	}
	if err = gzipWriter.Close(); err != nil {
		return fmt.Errorf("ошибка при завершении сжатия:%w", err)
	}

	// Теперь данные хранятся в gzippedData

	serverURL = fmt.Sprintf("http://%s/updates/", serverURL)
	req, err := http.NewRequest("POST", serverURL, &gzippedData)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса:%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip") // Указываем кодирование gzip

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при отправке запроса3:%w", err)
	}
	defer resp.Body.Close()

	//var responseBody []byte
	buf := make([]byte, 1024) // Размер буфера для чтения

	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("ошибка при чтении тела ответа:%w", err)
		}
		if n == 0 {
			break
		}
		//responseBody = append(responseBody, buf[:n]...)
	}

	return nil
}