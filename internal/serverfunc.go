package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
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
		fmt.Println("path[4]", metricValue)

		if lengpath != 5 {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return

		}

		if path[4] == "none" {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}

		num1, err := strconv.ParseInt(path[4], 10, 64)
		if err != nil {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return
		}

		if isInteger(path[4]) {
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

	if (len(metricName) > 0) && (path[4] == "") {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {

		num, err := strconv.ParseFloat(path[4], 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		if _, err1 := strconv.ParseFloat(path[4], 64); err1 == nil {
			fmt.Fprintf(w, "%v", num) // Возвращаем текущее значение метрики в текстовом виде
			mc.Storage.SaveMetric(path[2], metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)

		}

		if _, err1 := strconv.ParseInt(path[4], 10, 64); err1 == nil {
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
	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)
	fmt.Println("http.MethodGet", http.MethodGet)

	if lengpath != 4 {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if path[1] != "value" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}
	if path[2] != "gauge" && path[2] != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if path[2] == "counter" {
		num1, found := mc.Storage.counters[path[3]]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)

	}
	if path[2] == "gauge" {

		num1, found := mc.Storage.gauges[path[3]]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)

	}

}

func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		recorder := newResponseRecorder(w)
		next.ServeHTTP(recorder, r)

		elapsed := time.Since(startTime)
		logger.Info("Request processed",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("elapsed_time", elapsed),
			zap.Int("status_code", recorder.Status()),
			zap.Int("response_size", recorder.Size()),
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w}
}
func (r *responseRecorder) Write(data []byte) (int, error) {
	size, err := r.ResponseWriter.Write(data)
	r.size += size
	return size, err
}

func (r *responseRecorder) Status() int {
	return r.status
}

func (r *responseRecorder) Size() int {
	return r.size
}
func Init() {
	// Инициализация логгера
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger")
	}
	defer logger.Sync() // flushes buffer, if any
}
