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
		println("LENGTH", lengpath)
		// Обрабатываем полученные метрики
		// Преобразование строки во float64

		switch r.Method {
		case http.MethodPost:

			if path[1] != "update" {

				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Metric name not provided")
				return
			}

			if path[2] != "gauge" && path[2] != "counter" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "StatusBadRequest")
				return
			}

			if path[2] == "counter" {

				if lengpath != 5 {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintln(w, "StatusNotFound")
					return

				}

				if path[4] == "none" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, "StatusBadRequest")
					return

				}

				num1, err := strconv.ParseInt(path[4], 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintln(w, "StatusNotFound")
					return
				}

				if isInteger(path[4]) {

					w.WriteHeader(http.StatusOK)
					fmt.Fprintln(w, "StatusOK")

					storage.SaveMetric(path[2], path[3], num1)

					return

				} else {
					w.WriteHeader(http.StatusBadRequest)
					return

				}
			}
			if lengpath == 4 && path[3] == "" {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "Metric name not provided")
				return
			}

			if (len(path[3]) > 0) && (path[4] == "") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "StatusBadRequest")

				return
			}

			if path[2] == "gauge" {

				num, err := strconv.ParseFloat(path[4], 64)
				if err != nil {
					fmt.Println("Ошибка при преобразовании строки во float64:", err)
					return
				}

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

			num1, err := strconv.ParseFloat(path[1], 64)

			if (err != nil) && (lengpath == 2) {

				allMetrics := storage.GetAllMetrics()

				// Form an HTML page with the list of all metrics and their values
				html := "<html><head><title>Metrics</title></head><body><h1>Metrics List</h1><ul>"
				for name, value := range allMetrics {
					html += fmt.Sprintf("<li>%s: %v</li>", name, value)
				}
				html += "</ul></body></html>"

				// Set the Content-Type header to indicate that the response is HTML
				w.Header().Set("Content-Type", "text/html; charset=utf-8")

				// Write the HTML page to the response
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, html)

				return
			}

			if err != nil {

				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "StatusNotFound")

				//fmt.Println("Ошибка при преобразовании строки во float64:", err)
				return
			}
			if (path[2] != "gauge") && (path[2] != "counter") {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "StatusNotFound")
				return
			}

			if path[2] == "counter" {

				num, err := strconv.ParseInt(path[1], 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintln(w, "StatusNotFound")
					return
				}
				storage.SaveMetric(path[2], path[3], num)
				return
			}
			if lengpath != 4 {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "StatusNotFound")
				return

			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "StatusOK")
			storage.SaveMetric(path[2], path[3], num1)

		}
	}
}
