package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func isNotInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err != nil
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
		path := strings.Split(r.URL.Path, "/")

		// Проверяем, что имя метрики не пустое
		fmt.Println("PATH", r.URL.Path)
		fmt.Println("LENGTH PATH", len(path))
		for i := 0; i < len(path); i++ {
			fmt.Println(path[i])
		}
		if path[0] != "/update/" {

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}

		if path[2] == "" && path[3] == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Metric name not provided", http.StatusNotFound)
			return
		}
		if path[1] != "gauge" && path[1] != "counter" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric type")
			return
		}

		if (len(metricName) > 0) && (metricValueStr == "") {
			w.WriteHeader(http.StatusBadRequest)
		}

		if path[1] == "counter" {

			if isInteger(path[3]) {
				w.WriteHeader(http.StatusOK)

			} else {
				w.WriteHeader(http.StatusBadRequest)

			}
		}

		if isNotInteger(path[4]) {
			w.WriteHeader(http.StatusBadRequest)

		}

		//metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)

		// Обрабатываем полученные метрики
		storage.SaveMetric(metricType, metricName, metricValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
