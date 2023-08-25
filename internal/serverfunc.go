package internal

import (
	"fmt"
	"net/http"

	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

func (mc *HandlerDependencies) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		mc.HandleGetRequest(w, r)
	})

	r.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		mc.HandlePostRequest(w, r)
	})

	r.Get("/value/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		mc.HandleGetRequest(w, r)
	})

	return r
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (mc *HandlerDependencies) HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка POST-запроса

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)
	fmt.Println("http.MethodPost:=", http.MethodPost)

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {
		fmt.Println("lengpath path2=counter", lengpath)
		fmt.Println("metricValue", metricValue)

		if lengpath != 5 {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return

		}

		if metricValue == "none" {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}

		num1, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return
		}

		if isInteger(metricValue) {
			fmt.Println("Num1 в ветке POST ", num1)

			fmt.Fprintf(w, "%v", num1)

			mc.Storage.SaveMetric(metricType, metricName, num1)

			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}
	}
	if lengpath == 4 && metricName == "" {
		http.Error(w, "Metric name not provided", http.StatusBadRequest)
		return
	}

	if (len(metricName) > 0) && (metricValue == "") {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {

		num, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		if _, err1 := strconv.ParseFloat(metricValue, 64); err1 == nil {
			fmt.Fprintf(w, "%v", num) // Возвращаем текущее значение метрики в текстовом виде
			mc.Storage.SaveMetric(path[2], metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)

		}

		if _, err1 := strconv.ParseInt(metricValue, 10, 64); err1 == nil {
			fmt.Fprintf(w, "%v", num) // Возвращаем текущее значение метрики в текстовом виде
			mc.Storage.SaveMetric(path[2], path[3], num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

	}

}

func (mc *HandlerDependencies) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса
	metricType := chi.URLParam(r, "metricType")
	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)
	fmt.Println("http.MethodGet", http.MethodGet)

	if lengpath != 4 {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		num1, found := mc.Storage.counters[path[3]]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)

	}
	if metricType == "gauge" {

		num1, found := mc.Storage.gauges[path[3]]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)

	}

}
