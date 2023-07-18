package internal

import (
	"fmt"
	"log"
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

		newPath := r.URL.Path + r.FormValue("metricType") + "/" + r.FormValue("metricName") + "/" + r.FormValue("metricValue")
		log.Println("URL Path:", newPath)

		// Разбиваем путь запроса на части
		parts := strings.Split(newPath, "/")
		length := len(parts)
		log.Println("LENGTH PATH", length)
		if length != 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Извлекаем тип метрики, имя метрики и значение метрики из пути запроса
		metricType := parts[2]
		metricName := parts[3]
		metricValue := parts[4]
		log.Println("TYPE NAME VALUE", metricType, metricName, metricValue)

		// Преобразуем значение метрики в соответствующий тип
		var parsedValue interface{}
		var err error

		if metricType == "gauge" {
			parsedValue, err = strconv.ParseFloat(metricValue, 64)
		} else if metricType == "counter" {
			parsedValue, err = strconv.ParseInt(metricValue, 10, 64)
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
		storage.ProcessMetrics(metricType, metricName, parsedValue)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Metric received")
	}
}
