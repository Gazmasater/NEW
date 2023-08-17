package serverin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

// InitLogger инициализирует логгер для использования в WithLogging.
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar = logger.Sugar()
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// создаем ResponseRecorder, чтобы получить доступ к коду статуса и размеру ответа
		rr := httptest.NewRecorder()

		// вызываем оригинальный хендлер
		h.ServeHTTP(rr, r)

		// получаем код статуса и размер ответа
		statusCode := rr.Code
		responseSize := rr.Body.Len()

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		logger, err := zap.NewDevelopment()
		if err != nil {
			// Обработка ошибки
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		sugar.Infow(
			"uri", uri,
			"method", method,
			"duration", duration,
			"statusCode", strconv.Itoa(statusCode), // Преобразование числового значения в строку
			"responseSize", strconv.Itoa(responseSize), // Преобразование числового значения в строку
		)

		// Закрываем логгер после использования
		defer logger.Sync()

		// копируем данные из ResponseRecorder в http.ResponseWriter
		for k, v := range rr.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(statusCode)
		w.Write(rr.Body.Bytes())
	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}

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
			mc.Logger.Info("Num1 в ветке POST:", zap.Int64("value", num1))

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

			mc.Logger.Info("Текущее значение метрики num:", zap.Float64("value", num))

			mc.Storage.SaveMetric(metricType, metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)

		}

		if _, err1 := strconv.ParseInt(metricValue, 10, 64); err1 == nil {

			mc.Logger.Info("Возвращаем текущее значение метрики в текстовом виде:", zap.String("value", fmt.Sprintf("%f", num)))

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
	mc.Logger.Info("HTTP Method Get:", zap.String("method", http.MethodGet))

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		num1, found := mc.Storage.counters[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		mc.Logger.Info("Значение num1:", zap.Int64("value", num1))

		fmt.Fprintf(w, "%v", num1)

	}
	if metricType == "gauge" {

		num1, found := mc.Storage.gauges[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)
		mc.Logger.Info("Значение измерителя", zap.String("metricName", metricName), zap.Float64("value", num1))

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
	InitLogger()
	r := chi.NewRouter()
	r.Use(WithLogging)

	r.Get("/value/{metricType}/{metricName}", mc.handleGetRequest)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", mc.handlePostRequest)

	r.Get("/metrics/", mc.handleMetrics)

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
