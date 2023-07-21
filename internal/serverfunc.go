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
		metricValue, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			// Обработка ошибки, если metricValueStr не может быть преобразовано в float64
			// Или вернуть сообщение об ошибке клиенту.
		}

		println("metricValueStr", metricValueStr)

		// Преобразуем значение метрики в соответствующий тип

		path := strings.Split(r.URL.Path, "/")
		lengpath := len(path)

		fmt.Println("PATH", r.URL.Path)
		fmt.Println("LENGTH PATH", lengpath)
		for i := 0; i < len(path); i++ {
			fmt.Println(i, path[i])
		}
		if path[1] != "update" {

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}

		if lengpath == 4 && path[3] == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}
		if path[2] != "gauge" && path[2] != "counter" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric type")
			return
		}

		if (len(path[3]) > 0) && (path[4] == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if path[2] == "counter" {

			if isInteger(path[4]) {
				w.WriteHeader(http.StatusOK)
				return

			} else {
				w.WriteHeader(http.StatusBadRequest)
				return

			}
		}

		if path[2] == "gauge" {

			if _, err := strconv.ParseFloat(path[4], 64); err == nil {
				w.WriteHeader(http.StatusOK)
				return

			} else {
				w.WriteHeader(http.StatusBadRequest)

			}

			if _, err := strconv.ParseInt(path[4], 10, 64); err == nil {
				w.WriteHeader(http.StatusOK)
				return

			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		}

		//metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)

		// Обрабатываем полученные метрики
		println("VALUE перед SAVEMETRIC", metricValue)
		storage.SaveMetric(metricType, metricName, metricValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
