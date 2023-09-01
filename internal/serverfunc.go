package internal

import (
	"encoding/json"

	"net/http"
	"strconv"
	"strings"

	//	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (mc *HandlerDependencies) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return LoggingMiddleware(mc.Logger, next)
	})

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		mc.HandleGetRequest(w, r)
	})

	r.Post("/update/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		mc.HandlePostRequest(w, r)
	})

	r.Post("/value/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		mc.HandleGetRequest(w, r)
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
	contentType := r.Header.Get("Content-Type")

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {

		if lengpath != 5 {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return

		}

		if metricValue == "none" {
			println("metricValuenone")
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}

		num1, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return
		}

		if isInteger(metricValue) {
			if contentType == "application/json" {

				mc.Storage.SaveMetric(metricType, metricName, num1)
				createAndSendUpdatedMetricCounter(w, metricName, metricType, int64(num1))
				return
			} else {
				w.Write([]byte(strconv.FormatInt(num1, 10)))

				mc.Storage.SaveMetric(metricType, metricName, num1)
				return
			}

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
			println("strconv.ParseFloat")
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		if _, err1 := strconv.ParseFloat(metricValue, 64); err1 == nil {

			if contentType == "application/json" {
				mc.Storage.SaveMetric(metricType, metricName, num)

			} else {
				w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))
				mc.Storage.SaveMetric(metricType, metricName, num)
				return
			}

		} else {

			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		// if _, err1 := strconv.ParseInt(metricValue, 10, 64); err1 == nil {
		// 	if contentType == "application/json" {
		// 		mc.Storage.SaveMetric(metricType, metricName, num)

		// 	} else {
		// 		w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))
		// 		mc.Storage.SaveMetric(metricType, metricName, num)
		// 		return
		// 	}

		// } else {
		// 	println("strconv.ParseInt GAUGE")
		// 	http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		// 	return
		// }
		createAndSendUpdatedMetric(w, metricName, metricType, float64(num))

	}

}

func (mc *HandlerDependencies) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	// Обработка GET-запроса
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)
	//fmt.Println("http.MethodGet", http.MethodGet)

	if lengpath != 4 {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		num1, found := mc.Storage.counters[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}
		if contentType == "application/json" {
			createAndSendUpdatedMetricCounter(w, metricName, metricType, int64(num1))
			return
		} else {

			w.Write([]byte(strconv.FormatInt(num1, 10)))
		}
		return
	}
	if metricType == "gauge" {

		num, found := mc.Storage.gauges[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}
		if contentType == "application/json" {
			createAndSendUpdatedMetric(w, metricName, metricType, float64(num))
			return
		} else {

			w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))
		}

	}

}

func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateLogger()
		//startTime := time.Now()

		recorder := newResponseRecorder(w)
		next.ServeHTTP(recorder, r)

		// elapsed := time.Since(startTime)
		// logger.Info("Request processed",
		// 	zap.String("uri", r.RequestURI),
		// 	zap.String("method", r.Method),
		// 	zap.Duration("elapsed_time", elapsed),
		// 	zap.Int("status_code", recorder.Status()),
		// 	zap.Int("response_size", recorder.Size()),
		// )
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

func createAndSendUpdatedMetric(w http.ResponseWriter, metricName, metricType string, num float64) {
	// Создайте экземпляр структуры с обновленным значением Value
	updatedMetric := &Metrics{
		ID:    metricName,
		MType: metricType,
		Value: &num,
	}
	Init()
	// Сериализуйте структуру в JSON
	responseData, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(w, "Ошибка при сериализации данных в JSON", http.StatusInternalServerError)
		return
	}
	//logger.Info("Сериализированные данные в JSON responseData GAUGE", zap.String("json_data", string(responseData)))
	// Установите Content-Type и статус код для ответа
	w.Header().Set("Content-Type", "application/json")

	// Отправьте JSON в теле ответа
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)

}

func createAndSendUpdatedMetricCounter(w http.ResponseWriter, metricName, metricType string, num int64) {
	// Создайте экземпляр структуры с обновленным значением Value
	Init()
	updatedMetric := &Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: &num,
	}

	// Сериализуйте структуру в JSON
	responseData, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(w, "Ошибка при сериализации данных в JSON", http.StatusInternalServerError)
		return
	}

	//	logger.Info("Сериализированные данные в JSON responseData COUNTER", zap.String("json_data", string(responseData)))
	// Установите Content-Type и статус код для ответа
	w.Header().Set("Content-Type", "application/json")

	// Отправьте JSON в теле ответа
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)
	//	fmt.Println("createAndSendUpdatedMetricCounter Тело ответа:&&&&&&&&&&", string(responseData))

}
