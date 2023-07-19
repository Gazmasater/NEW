package internal

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"strconv"
	"strings"
)

// ResponseRecorderWithLog is a custom implementation of http.ResponseWriter
// that records the response data and writes it to the log.
type ResponseRecorderWithLog struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func NewResponseRecorderWithLog(w http.ResponseWriter) *ResponseRecorderWithLog {
	return &ResponseRecorderWithLog{
		ResponseWriter: w,
		buffer:         bytes.NewBuffer(nil),
	}
}

func (rw *ResponseRecorderWithLog) Write(data []byte) (int, error) {
	rw.buffer.Write(data)
	return rw.ResponseWriter.Write(data)
}

func (rw *ResponseRecorderWithLog) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseRecorderWithLog) Flush() {
	// Log the response data here
	responseDump := rw.buffer.String()
	log.Printf("Ответ:\n%s\n", responseDump)
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
		//	newpath := r.URL.Path + metricType

		// Преобразуем значение метрики в соответствующий тип
		var metricValue interface{}
		var err error

		// Проверяем, что имя метрики не пустое
		if metricName == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Metric name not provided")
			return
		}

		if metricType == "gauge" {
			// Проверяем, является ли значение действительным числом (не целым)
			if !strings.Contains(metricValueStr, ".") {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid metric value for gauge, expected a floating-point number")
				return
			}
			metricValue, err = strconv.ParseFloat(metricValueStr, 64)
		} else if metricType == "counter" {
			metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric type")
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric value")
			return
		}

		// Обрабатываем полученные метрики
		storage.SaveMetric(metricType, metricName, metricValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
