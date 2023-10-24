package internal

import (
	"bufio"

	"compress/gzip"

	"encoding/json"
	"fmt"

	"net/http"
	"strings"

	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (mc *HandlerDependencies) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(GzipMiddleware)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", mc.HandlePostRequest)

	r.Post("/value/{metricType}/{metricName}", mc.HandleGetRequest)

	r.Get("/value/{metricType}/{metricName}", mc.HandleGetRequest)

	r.Get("/metrics", mc.HandleGetRequest)

	r.Get("/", mc.HandleGetRequestHTML)

	return r
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// хендлер post путь /update/
func (mc *HandlerDependencies) HandlePostRequest(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	println("HandlePostRequest")

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {

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
				createAndSendUpdatedMetricCounterTEXT(w, metricName, metricType, int64(num1))
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

		w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))

	}

}

func (mc *HandlerDependencies) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	println("HandleGetRequest")
	contentType := r.Header.Get("Content-Type")
	// Обработка GET-запроса
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		println("HandleGetRequest  counter", mc.Storage.counters[metricName])

		num1, found := mc.Storage.counters[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}
		if contentType == "application/json" {
			createAndSendUpdatedMetricCounterJSON(w, metricName, metricType, int64(num1))
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
			createAndSendUpdatedMetricJSON(w, metricName, metricType, float64(num))
			return
		} else {

			w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))
		}

	}

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

func createAndSendUpdatedMetricJSON(w http.ResponseWriter, metricName, metricType string, num float64) {
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
	_, _ = w.Write([]byte("\n"))

}

func createAndSendUpdatedMetricCounterJSON(w http.ResponseWriter, metricName, metricType string, num int64) {
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

	// Установите Content-Type и статус код для ответа
	w.Header().Set("Content-Type", "application/json")

	// Отправьте JSON в теле ответа
	logger.Info("createAndSendUpdatedMetric Тело ответа", zap.String("response_body", string(responseData)))

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)
	_, _ = w.Write([]byte("\n"))

}

func createAndSendUpdatedMetricCounterTEXT(w http.ResponseWriter, metricName, metricType string, num int64) {
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
	w.Header().Set("Content-Type", "text/plain")

	// Отправьте JSON в теле ответа
	logger.Info("createAndSendUpdatedMetric Тело ответа", zap.String("response_body", string(responseData)))

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)
	_, _ = w.Write([]byte("\n"))
	fmt.Println("createAndSendUpdatedMetricCounter Тело ответа:&&&&&&&&&&", string(responseData))

}

func (mc *HandlerDependencies) HandleGetRequestHTML(w http.ResponseWriter, r *http.Request) {
	println("HandleGetRequestHTML")
	//contentType := r.Header.Get("Content-Type")

	// Получить список известных метрик
	metrics := mc.getKnownMetrics()

	// Генерировать HTML-страницу
	var htmlPage string
	for _, metric := range metrics {
		htmlPage += fmt.Sprintf("<p>%s: %v</p>", metric.Name, metric.Value)
	}

	// Отправить HTML-страницу как ответ
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlPage))
}

func (mc *HandlerDependencies) getKnownMetrics() []Metric {
	// Собрать список известных метрик
	var metrics []Metric

	for name, counter := range mc.Storage.counters {
		metrics = append(metrics, Metric{
			Name:  name,
			Value: int64(counter),
		})
	}

	for name, gauge := range mc.Storage.gauges {
		metrics = append(metrics, Metric{
			Name:  name,
			Value: float64(gauge),
		})
	}

	return metrics
}

type Metric struct {
	Name  string
	Value interface{}
}

func (mc *HandlerDependencies) WriteMetricToFile(metric *Metrics) error {
	// Открываем файл для чтения и записи
	file, err := os.OpenFile(mc.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		mc.Logger.Error("Ошибка при открытии файла для записи", zap.Error(err))
		return err
	}
	defer file.Close()

	// Читаем метрики из файла
	var metrics []Metrics
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var existingMetric Metrics
		if err := json.Unmarshal([]byte(line), &existingMetric); err == nil {
			if existingMetric.ID != metric.ID {
				// Если ID метрики не совпадает, добавляем ее в список метрик
				metrics = append(metrics, existingMetric)
			} else {
				// Если ID совпадает, устанавливаем тип из существующей метрики
				metric.MType = existingMetric.MType
			}
		} else {
			mc.Logger.Error("Ошибка при разборе JSON", zap.Error(err))
		}
	}

	if err := scanner.Err(); err != nil {
		mc.Logger.Error("Ошибка при чтении файла", zap.Error(err))
		return err
	}

	// Добавляем новую метрику к уже существующим метрикам
	metrics = append(metrics, *metric)

	// Перезаписываем файл с обновленными метриками
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	for _, updatedMetric := range metrics {
		if err := encoder.Encode(updatedMetric); err != nil {
			mc.Logger.Error("Ошибка при записи метрики в файл", zap.Error(err))
			return err
		}
	}
	return err
}

func WriteJSONToFile(fileStoragePath string, jsonData string) error {
	// Попробуем открыть файл для записи, или создадим его, если он не существует
	file, err := os.OpenFile(fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Определяем, поддерживает ли клиент сжатие Gzip
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gzWriter := GzipResponseWriter{Writer: gz, ResponseWriter: w}
			next.ServeHTTP(gzWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (grw GzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}
