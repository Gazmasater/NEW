package serverin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (mc *HandlerDependencies) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка POST-запроса

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {

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
			mc.Logger.Printf("Num1 в ветке POST: %v", num1)

			mc.Storage.SaveMetric(metricType, metricName, num1)

			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}
	}
	if metricName == "" {
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
			mc.Storage.SaveMetric(metricType, metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)

		}

		if _, err1 := strconv.ParseInt(metricValue, 10, 64); err1 == nil {
			fmt.Fprintf(w, "%v", num) // Возвращаем текущее значение метрики в текстовом виде
			mc.Storage.SaveMetric(metricType, metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

	}

}

func (mc *HandlerDependencies) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса

	log.Println("handleGetRequest")
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	mc.Logger.Println("http.MethodGet:", http.MethodGet)

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		num1, found := mc.Storage.counters[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)

	}
	if metricType == "gauge" {

		num1, found := mc.Storage.gauges[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)
		mc.Logger.Printf("Значение измерителя %s: %v", metricName, num1)

	}

}

func handleMetrics(mc *HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allMetrics := mc.Storage.GetAllMetrics()

		// Формируем JSON с данными о метриках
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Используем пакет encoding/json для преобразования данных в JSON и записи их в ResponseWriter.
		json.NewEncoder(w).Encode(allMetrics)
	}
}

func (mc *HandlerDependencies) handleMetrics(w http.ResponseWriter, r *http.Request) {
	handleMetrics(mc)(w, r)
}

func (mc *HandlerDependencies) Route() *chi.Mux {

	r := chi.NewRouter()

	// Создаем отдельные роутеры для каждого типа запросов
	getRouter := chi.NewRouter()
	getRouter.Get("/{metricType}/{metricName}", mc.handleGetRequest)

	postRouter := chi.NewRouter()
	postRouter.Post("/{metricType}/{metricName}/{metricValue}", mc.handlePostRequest)

	metricsRouter := chi.NewRouter()
	metricsRouter.Get("/", mc.handleMetrics)

	// Монтируем роутеры
	r.Mount("/value", getRouter)       // Монтирование роутера для GET-запросов
	r.Mount("/update", postRouter)     // Монтирование роутера для POST-запросов
	r.Mount("/metrics", metricsRouter) // Монтирование роутера для /metrics

	// Обработчик для случая, когда путь не соответствует заданному шаблону
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return r
}

func StartServer(address string, handler http.Handler) {
	// Создаем HTTP-сервер с настройками
	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	// Запуск HTTP-сервера через http.ListenAndServe()
	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s) // Преобразуем строку в целое число, игнорируя результат
	return err == nil         // Если ошибки нет, то строка является целым числом
}
