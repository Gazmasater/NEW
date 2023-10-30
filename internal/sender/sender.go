package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"project.com/internal"
	"project.com/internal/logger"
	"project.com/internal/models"
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
		hash := internal.ComputeHash(jsonData, "MyKey")
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
