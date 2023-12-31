package server

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"project.com/internal/models"
)

func (mc *app) HandlePostRequest(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")

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
			if contentType == "application/json" {

				mc.Storage.SaveMetric(metricType, metricName, num1)
				mc.createAndSendUpdatedMetricCounterTEXT(w, metricName, metricType, int64(num1))
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

func (mc *app) HandlePostRequestOptimiz(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {

		num, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		if isInteger(metricValue) {

			w.Write([]byte(strconv.FormatInt(num, 10)))

			mc.Storage.SaveCounter(metricType, metricName, num)

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}
	}
	if metricName == "" || (len(metricName) > 0 && metricValue == "") {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {
		num, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		mc.Storage.SaveGauge(metricType, metricName, num)

		responseData := []byte(strconv.FormatFloat(num, 'f', -1, 64))
		w.Write(responseData)

	}

}

// Хендлер для Get запроса
func (mc *app) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
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
		println("HandleGetRequest  counter", mc.Storage.GetCounters()[metricName])

		num1, found := mc.Storage.GetCounters()[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}
		if contentType == "application/json" {
			mc.createAndSendUpdatedMetricCounterJSON(w, metricName, metricType, int64(num1))
			return
		} else {

			w.Write([]byte(strconv.FormatInt(num1, 10)))
		}
		return
	}
	if metricType == "gauge" {

		num, found := mc.Storage.GetGauges()[metricName]
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

func (mc *app) updateHandlerJSON(w http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	metricsFromFile := make(map[string]models.Metrics)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&metric); err != nil {
		_ = fmt.Errorf("ошибка при разборе JSON: %w", err)
		return

	}

	fmt.Printf("Metric: %+v\n", metric)

	// Прочитать тело запроса

	if mc.Config.Restore {
		var err error
		metricsFromFile, err = mc.ReadMetricsFromFile()
		if err != nil {
			_ = fmt.Errorf("ошибка чтения метрик из файла:%w", err)
			return
		}
	}

	// Обработка "counter"
	if metric.MType == "counter" && metric.Delta != nil {
		currentValue, ok := metricsFromFile[metric.ID]

		if !ok {
			// Если метрики нет в файле, проверяем в хранилище
			if value, exists := mc.Storage.GetCounters()[metric.ID]; exists {
				currentValue = models.Metrics{
					MType: metric.MType,

					ID:    metric.ID,
					Delta: new(int64),
				}
				*currentValue.Delta = value
			} else {
				// Если метрики нет ни в файле, ни в хранилище, инициализируем ее с нулевым значением
				currentValue = models.Metrics{
					MType: metric.MType,
					ID:    metric.ID,
					Delta: new(int64),
				}
				*currentValue.Delta = 0

			}
		}

		*currentValue.Delta += *metric.Delta
		println("currentValue   currentValueИМЯ", currentValue.MType, currentValue.ID, *currentValue.Delta)

		mc.Storage.GetCounters()[metric.ID] = *currentValue.Delta
		metricsFromFile[metric.ID] = currentValue

	}

	// Обработка "gauge"
	if metric.MType == "gauge" && metric.Value != nil {
		// Обновляем или создаем метрику в слайсе
		metricsFromFile[metric.ID] = metric

		// Сохраняем обновленные метрики в хранилище
		mc.Storage.GetGauges()[metric.ID] = *metric.Value
	}

	// Запись обновленных метрик в файл
	for _, updatedMetric := range metricsFromFile {
		println("updatedMetric", updatedMetric.MType, updatedMetric.ID)
		dbErr := mc.WriteMetricToDatabase(updatedMetric)
		if dbErr != nil {
			log.Printf("Ошибка при записи метрики в базу данных: %s", dbErr)

		}

		if err := mc.WriteMetricToFile(&updatedMetric); err != nil {
			_ = fmt.Errorf("ошибка записи метрик в файл:%w", err)
			return
		}

	}

	// Отправляем значение метрики
	if updatedMetric, ok := metricsFromFile[metric.ID]; ok {
		switch metric.MType {
		case "counter":
			mc.createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, *updatedMetric.Delta)
		case "gauge":
			createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, *updatedMetric.Value)

		default:
			http.Error(w, "Метрика не найдена", http.StatusNotFound)

		}
	}

}

func (mc *app) updateHandlerJSONValue(w http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	// Проверка заголовка Content-Encoding на предмет GZIP
	if r.Header.Get("Content-Encoding") == "gzip" {
		// Если данные приходят в GZIP, создаем Reader для распаковки
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Ошибка при создании GZIP Reader", http.StatusBadRequest)
			return
		}
		defer reader.Close()
		r.Body = reader
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&metric); err != nil {
		http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что поля "id" и "type" заполнены
	if metric.ID == "" || metric.MType == "" {
		http.Error(w, "Поля 'id' и 'type' обязательны для заполнения", http.StatusBadRequest)
		return
	}
	// Прочитать метрики из файла

	metricsFromFile, err := mc.ReadMetricsFromFile()
	if err != nil {
		http.Error(w, "Ошибка чтения метрик из файла", http.StatusInternalServerError)
		return
	}
	// Проверить наличие нужной метрики в файле
	metricFromFile, exists := metricsFromFile[metric.ID]

	// Если метрика отсутствует в файле, проверьте хранилище
	if !exists {
		switch metric.MType {
		case "gauge":
			value, ok := mc.Storage.GetGauges()[metric.ID]
			if ok {
				// Метрика существует в хранилище, используйте значение из хранилища
				createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, value)
				return
			}
		case "counter":
			value, ok := mc.Storage.GetCounters()[metric.ID]
			if ok {
				// Метрика существует в хранилище, используйте значение из хранилища
				mc.createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, value)
				return

			}
		}

		// Если метрика отсутствует и в файле, и в хранилище, отправьте статус "Not Found"
		http.Error(w, "Метрика не найдена", http.StatusNotFound)
		return
	}

	switch metric.MType {
	case "gauge":
		createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, *metricFromFile.Value)
	case "counter":
		mc.createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, *metricFromFile.Delta)

	}
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

func createAndSendUpdatedMetricJSON(w http.ResponseWriter, metricName, metricType string, num float64) {
	updatedMetric := &models.Metrics{
		ID:    metricName,
		MType: metricType,
		Value: &num,
	}
	responseData, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(w, "Ошибка при сериализации данных в JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Отправьте JSON в теле ответа
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)
	_, _ = w.Write([]byte("\n"))

}

func (mc *app) createAndSendUpdatedMetricCounterJSON(w http.ResponseWriter, metricName, metricType string, num int64) {
	// Создайте экземпляр структуры с обновленным значением Value
	updatedMetric := &models.Metrics{
		ID:    metricName,
		MType: metricType,
		Delta: &num,
	}

	responseData, err := json.Marshal(updatedMetric)
	if err != nil {
		http.Error(w, "Ошибка при сериализации данных в JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)
	_, _ = w.Write([]byte("\n"))

}

func (mc *app) createAndSendUpdatedMetricCounterTEXT(w http.ResponseWriter, metricName, metricType string, num int64) {
	// Создайте экземпляр структуры с обновленным значением Value
	updatedMetric := &models.Metrics{
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
	w.Header().Set("Content-Type", "text/plain")

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseData)
	_, _ = w.Write([]byte("\n"))

}

func (mc *app) HandleGetRequestHTML(w http.ResponseWriter, r *http.Request) {
	println("HandleGetRequestHTML")

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

func (mc *app) getKnownMetrics() []models.Metric {
	// Собрать список известных метрик
	var metrics []models.Metric

	for name, counter := range mc.Storage.GetCounters() {
		metrics = append(metrics, models.Metric{
			Name:  name,
			Value: int64(counter),
		})
	}

	for name, gauge := range mc.Storage.GetGauges() {
		metrics = append(metrics, models.Metric{
			Name:  name,
			Value: float64(gauge),
		})
	}

	return metrics
}

func (mc *app) WriteMetricToFile(metric *models.Metrics) error {

	file, err := os.OpenFile(mc.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		mc.Logger.Error("Ошибка при открытии файла для записи", zap.Error(err))
		return err
	}
	defer file.Close()

	// Читаем метрики из файла
	var metrics []models.Metrics
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var existingMetric models.Metrics
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

func (mc *app) WriteJSONToFile(fileStoragePath string, jsonData string) error {
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

func (mc *app) ReadMetricsFromFile() (map[string]models.Metrics, error) {
	metricsMap := make(map[string]models.Metrics)

	file, err := os.OpenFile(mc.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		mc.Logger.Error("Ошибка при открытии файла для чтения", zap.Error(err))
		return metricsMap, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var metric models.Metrics
		if err := json.Unmarshal([]byte(line), &metric); err == nil {
			metricsMap[metric.ID] = metric
		} else {
			mc.Logger.Error("Ошибка при разборе JSON", zap.Error(err))
		}
	}

	if err := scanner.Err(); err != nil {
		mc.Logger.Error("Ошибка при чтении файла", zap.Error(err))
		return metricsMap, err
	}

	return metricsMap, nil
}

func (mc *app) Ping(w http.ResponseWriter, r *http.Request) {
	println("PING")
	println("DatabaseDSN", mc.Config.DatabaseDSN)
	db, err := sql.Open("postgres", mc.Config.DatabaseDSN)
	if err != nil {
		mc.Logger.Error("Ping ошибка", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Проверка на ошибку открытия базы данных
	if err := db.Ping(); err != nil {
		mc.Logger.Error("Ping ошибка", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError) // Меняем статус на 503 Service Unavailable
		return
	}

	// Если успешно, возвращаем HTTP-статус 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database is working\n")
}

func (mc *app) WriteMetricToDatabase(metric models.Metrics) error {
	var query string
	var args []interface{}

	switch metric.MType {
	case "gauge":
		query = "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)"
		args = []interface{}{metric.ID, metric.MType, metric.Value}
	case "counter":
		query = "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)"
		args = []interface{}{metric.ID, metric.MType, metric.Delta}
	default:
		log.Printf("Неизвестный тип метрики: %s", metric.MType)
		return fmt.Errorf("неизвестный тип метрики")

	}
	if mc.DB == nil {
		log.Println("Ошибка: mc.DB не инициализирован.")
		return fmt.Errorf("mc.DB не инициализирован")
	}

	// Проверяем, существует ли метрика с такими же значениями name и type
	var count int
	err := mc.DB.QueryRow("SELECT COUNT(*) FROM metrics WHERE name = $1 AND type = $2", metric.ID, metric.MType).Scan(&count)
	if err != nil {
		log.Printf("Ошибка при проверке наличия метрики: %s", err)
		return err
	}

	if count > 0 {
		// Метрика с такими значениями name и type существует, удаляем ее
		_, err := mc.DB.Exec("DELETE FROM metrics WHERE name = $1 AND type = $2", metric.ID, metric.MType)
		if err != nil {
			log.Printf("Ошибка при удалении метрики: %s", err)
			return err
		}
	}

	// Теперь выполняем вставку новой метрики
	_, err = mc.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Ошибка при записи метрики в базу данных: %s", err)
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

func (mc *app) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	println("MetricsHandler")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем, был ли запрос сжат с использованием gzip
	isGzip := false
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		isGzip = true
	}

	// Чтение тела запроса
	var bodyBuffer bytes.Buffer
	_, err := bodyBuffer.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Если запрос сжат с использованием gzip, распаковываем его
	if isGzip {
		gzipReader, err := gzip.NewReader(&bodyBuffer)
		if err != nil {
			http.Error(w, "Failed to unpack gzip data", http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()

		var unpackedBuffer bytes.Buffer
		_, err = unpackedBuffer.ReadFrom(gzipReader)
		if err != nil {
			http.Error(w, "Failed to unpack gzip data", http.StatusBadRequest)
			return
		}

		bodyBuffer = unpackedBuffer
	}

	fmt.Printf("Body Data: %s\n", bodyBuffer.String())

	// Создаем переменную для хранения распакованных метрик
	var batch []models.Metrics
	// Распаковка JSON-тела запроса в объект batch
	err = json.Unmarshal(bodyBuffer.Bytes(), &batch)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Обработка и сохранение полученных метрик
	if err := mc.updateHandlerJSONforBatch(batch); err != nil {
		http.Error(w, fmt.Sprintf("Error processing metrics: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Batch contents: %+v\n", batch)

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Metrics received and processed successfully")
}

func (mc *app) updateHandlerJSONforBatch(metrics []models.Metrics) error {
	println("updateHandlerJSONforBatch")
	// var metricsFromFile map[string]Metrics
	var err error
	metricsFromFile := make(map[string]models.Metrics)
	if mc.Config.Restore {
		metricsFromFile, err = mc.ReadMetricsFromFile()
		if err != nil {
			return fmt.Errorf("ошибка чтения метрик из файла: %w", err)
		}
	}

	for _, metric := range metrics {

		if metric.MType == "counter" && metric.Delta != nil {
			currentValue, ok := metricsFromFile[metric.ID]

			if !ok {
				// Если метрики нет в файле, проверяем в хранилище
				if value, exists := mc.Storage.GetCounters()[metric.ID]; exists {
					currentValue = models.Metrics{
						MType: metric.MType,
						ID:    metric.ID,
						Delta: new(int64),
					}
					*currentValue.Delta = value
				} else {
					// Если метрики нет ни в файле, ни в хранилище, инициализируем ее с нулевым значением
					currentValue = models.Metrics{
						MType: metric.MType,
						ID:    metric.ID,
						Delta: new(int64),
					}
					*currentValue.Delta = 0
				}
			}

			*currentValue.Delta += *metric.Delta

			// Сохраняем обновленные метрики в хранилище
			mc.Storage.GetCounters()[metric.ID] = *currentValue.Delta
			metricsFromFile[metric.ID] = currentValue
		}

		// Обработка "gauge"
		if metric.MType == "gauge" && metric.Value != nil {
			// Обновляем или создаем метрику в слайсе
			metricsFromFile[metric.ID] = metric

			// Сохраняем обновленные метрики в хранилище
			mc.Storage.GetGauges()[metric.ID] = *metric.Value
		}

		// Запись обновленных метрик в базу
		for _, updatedMetric := range metricsFromFile {
			if mc.Config.DatabaseDSN != "" {
				dbErr := mc.WriteMetricToDatabase(updatedMetric)
				if dbErr != nil {
					log.Printf("Ошибка при записи метрики в базу данных: %s", dbErr)
				}
			}
			// Запись обновленных метрик в файл
			if err := mc.WriteMetricToFile(&updatedMetric); err != nil {
				return fmt.Errorf("ошибка записи метрик в файл:%w", err)

			}
		}
	}
	return nil
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Обоащение к базе с интервалами
func (mc *app) SetupDatabase() error {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", mc.Config.DatabaseDSN)
	if err != nil {
		log.Printf("Ошибка при открытии базы данных: %v", err)
		return err
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Printf("Ошибка при проверке соединения с базой данных: %v", err)
		return err
	}

	// Запрос для создания таблицы
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS metrics (
            name VARCHAR(255) NOT NULL,
            type VARCHAR(50) NOT NULL,
            value DOUBLE PRECISION,
            delta BIGINT
        )
    `

	// Выполняем запрос для создания таблицы
	_, err = db.Exec(createTableQuery)
	if err != nil {
		pqErr, isPQError := err.(*pq.Error)
		if isPQError && pqErr.Code == "23505" {
			// Код "23505" соответствует ошибке уникального нарушения.
			log.Printf("Ошибка при создании таблицы: %v (UniqueViolation)", err)
			return err
		}
		log.Printf("Ошибка при создании таблицы: %v", err)
		return err
	}

	return nil
}

func KeyHashMiddleware(expectedKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка наличия ключа
		key := r.URL.Query().Get("k")

		if key != "" {
			if key != expectedKey {
				http.Error(w, "Несоответствие ключа", http.StatusBadRequest)
				return
			}

			// Если ключ существует, выполните хеширование
			h := sha256.New()
			h.Write([]byte(r.URL.Path))
			h.Write([]byte(r.Method))
			h.Write([]byte(r.Header.Get("Content-Type")))

			// Чтение и хеширование тела запроса
			body := r.Body
			defer body.Close()
			tee := io.TeeReader(body, h)

			// Прочитать JSON из тела запроса, если это JSON
			var jsonData interface{}
			decoder := json.NewDecoder(tee)
			_ = decoder.Decode(&jsonData)

			// Преобразование хеша в строку в шестнадцатеричном формате
			hash := hex.EncodeToString(h.Sum(nil))

			// Устанавливаем заголовок HashSHA256 в запросе
			r.Header.Set("HashSHA256", hash)
		}

		// Продолжаем выполнение следующего обработчика
		next.ServeHTTP(w, r)
	})
}
