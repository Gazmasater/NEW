package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func HandleUpdate(storage *MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Разбиваем путь запроса на составл
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Извлекаем тип метрики, имя метрики и значение метрики из пути запроса
		metricType := parts[2]
		metricName := parts[3]
		metricValue := parts[4]

		// Преобразуем значение метрики в соответствующий тип
		var parsedValue interface{}
		var err error

		println("ТИП МЕТРИКИ", metricType)

		if metricType == "gauge" {
			parsedValue, err = strconv.ParseFloat(metricValue, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid metric value")
				return
			}
		} else if metricType == "counter" {
			parsedValue, err = strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid metric value")
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid metric type")
			return
		}
		// Обрабатываем полученные метрики

		storage.ProcessMetrics(metricType, metricName, parsedValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
