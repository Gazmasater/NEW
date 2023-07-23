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

		// Преобразуем значение метрики в соответствующий тип
		println("r.URL.Path", r.URL.Path)
		println("METHOD", r.Method)
		path := strings.Split(r.URL.Path, "/")
		lengpath := len(path)
		// Обрабатываем полученные метрики
		// Преобразование строки во float64
		num, err := strconv.ParseFloat(path[4], 64)
		if err != nil {
			fmt.Println("Ошибка при преобразовании строки во float64:", err)
			return
		}
		num1, err := strconv.ParseFloat(path[4], 64)
		if err != nil {
			fmt.Println("Ошибка при преобразовании строки во float64:", err)
			return
		}
		println("path[4] перед storage.SaveMetric", num)

		switch r.Method {
		case http.MethodPost:

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
					storage.SaveMetric(path[2], path[3], num1)

					return

				} else {
					w.WriteHeader(http.StatusBadRequest)
					return

				}
			}

			if path[2] == "gauge" {

				if _, err := strconv.ParseFloat(path[4], 64); err == nil {

					w.WriteHeader(http.StatusOK)
					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					w.WriteHeader(http.StatusBadRequest)

				}

				if _, err := strconv.ParseInt(path[4], 10, 64); err == nil {
					w.WriteHeader(http.StatusOK)
					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

			}

			//metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Metric received")
		case http.MethodGet:
			if lengpath != 5 {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "Invalid request format")
				return
			}

			metricType := path[2]
			metricName := path[3]
			storage.SaveMetric(path[2], path[3], num)

			// Получаем метрику из хранилища
			value, found := storage.GetMetric(metricType, metricName)
			if !found {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "Metric not found")
				return
			}

			// Преобразуем значение метрики в текст и отправляем в ответ
			responseValue := fmt.Sprintf("%v", value)

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, responseValue)
		}
	}

}
