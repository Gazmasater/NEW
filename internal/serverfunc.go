package internal

import (
	"fmt"
	"net/http"
	"strconv"
)

func HandleUpdate(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("METHOD PATH", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/update/" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid URL path")
			return
		}

		// Извлекаем параметры из формы запроса
		metricType := r.FormValue("metricType")
		metricName := r.FormValue("metricName")
		metricValueStr := r.FormValue("metricValue")

		// Проверяем, что имя метрики и значение не пустые
		if metricName == "" && metricValueStr == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Metric name and value not provided")
			return
		}

		// Проверяем, что имя метрики не пустое
		if metricName == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}

		// Проверяем, что значение метрики не пустое
		if metricValueStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Metric value not provided")
			return
		}

		var metricValue interface{}
		var err error

		if metricType == "gauge" {
			// Проверяем, является ли значение действительным числом (int64 или float64)
			if value, err := strconv.ParseInt(metricValueStr, 10, 64); err == nil {
				metricValue = value
			} else if valueFloat, err := strconv.ParseFloat(metricValueStr, 64); err == nil {
				metricValue = valueFloat
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid metric value for gauge, expected a floating-point or integer number")
				return
			}
		} else if metricType == "counter" {
			// Проверяем, является ли значение действительным числом (int64)
			metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid metric value for counter, expected an integer number")
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric type")
			return
		}

		// Обрабатываем полученные метрики
		storage.SaveMetric(metricType, metricName, metricValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
