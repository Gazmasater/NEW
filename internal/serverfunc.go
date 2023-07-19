package internal

import (
	"fmt"
	"net/http"
	"strconv"
)

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// ResponseRecorderWithLog is a custom implementation of http.ResponseWriter
// that records the response data and writes it to the log.

func HandleUpdate(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Извлекаем параметры из формы запроса
		metricType := r.FormValue("metricType")
		metricName := r.FormValue("metricName")
		metricValueStr := r.FormValue("metricValue")
		//	newpath := r.URL.Path + metricType

		// Преобразуем значение метрики в соответствующий тип
		var metricValue interface{}

		// Проверяем, что имя метрики не пустое
		if r.URL.Path != "/update/" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}

		if metricName == "" && metricValueStr == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}
		if metricType != "gauge" && metricType != "counter" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric type")
			return
		}

		if (len(metricName) > 0) && (metricValueStr == "") {
			w.WriteHeader(http.StatusBadRequest)
		}

		if metricType == "counter" {

			if isInteger(metricValueStr) {
				w.WriteHeader(http.StatusOK)

			} else {
				w.WriteHeader(http.StatusBadRequest)

			}
		}

		//metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)

		// Обрабатываем полученные метрики
		storage.SaveMetric(metricType, metricName, metricValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
