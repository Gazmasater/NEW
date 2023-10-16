package internal

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"
	"strconv"

	"time"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func (mc *HandlerDependencies) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5)) // Здесь 5 - это уровень сжатия (0-9), где 0 - без сжатия, а 9 - максимальное сжатие.

	r.Use(func(next http.Handler) http.Handler {
		return LoggingMiddleware(mc.Logger, next)
	})

	r.Post("/update/", mc.updateHandlerJSON)

	r.Post("/value/", mc.updateHandlerJSONValue)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", mc.HandlePostRequest)

	r.Post("/value/{metricType}/{metricName}", mc.HandleGetRequest)

	r.Get("/value/{metricType}/{metricName}", mc.HandleGetRequest)

	r.Get("/metrics", mc.HandleGetRequest)

	r.Get("/", mc.HandleGetRequestHTML)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		mc.Ping(w, r)
	})

	return r
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

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

func (mc *HandlerDependencies) updateHandlerJSON(w http.ResponseWriter, r *http.Request) {
	var metric Metrics

	metricsFromFile := make(map[string]Metrics)

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
	if metric.MType == "counter" {
		currentValue, ok := metricsFromFile[metric.ID]

		if !ok {
			// Если метрики нет в файле, проверяем в хранилище
			if value, exists := mc.Storage.counters[metric.ID]; exists {
				currentValue = Metrics{
					ID:    metric.ID,
					Delta: new(int64),
				}
				*currentValue.Delta = value
			} else {
				// Если метрики нет ни в файле, ни в хранилище, инициализируем ее с нулевым значением
				currentValue = Metrics{
					ID:    metric.ID,
					Delta: new(int64),
				}
				*currentValue.Delta = 0

			}
		}

		*currentValue.Delta += *metric.Delta

		// Обновляем или создаем метрику в слайсе

		metricsFromFile[metric.ID] = currentValue
		// Сохраняем обн}овленные метрики в хранилище
		mc.Storage.SaveMetric(metric.MType, metric.ID, *currentValue.Delta)

	}

	// Обработка "gauge"
	if metric.MType == "gauge" && metric.Value != nil {
		// Обновляем или создаем метрику в слайсе
		metricsFromFile[metric.ID] = metric

		// Сохраняем обновленные метрики в хранилище
		mc.Storage.gauges[metric.ID] = *metric.Value
	}

	// Запись обновленных метрик в файл
	for _, updatedMetric := range metricsFromFile {

		if err := mc.WriteMetricToFile(&updatedMetric); err != nil {
			_ = fmt.Errorf("ошибка записи метрик в файл:%w", err)
			return
		}

	}

	if mc.Config.DatabaseDSN != "" {
		mc.SetupDatabase()
		dbErr := mc.WriteMetricToDatabase(metric)
		if dbErr != nil {
			log.Printf("Ошибка при записи метрики в базу данных: %s", dbErr)
		}
	}

	// Отправляем значение метрики
	if updatedMetric, ok := metricsFromFile[metric.ID]; ok {
		if metric.MType == "counter" {
			createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, *updatedMetric.Delta)
		} else if metric.MType == "gauge" {
			createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, *updatedMetric.Value)
		}
	} else {
		http.Error(w, "Метрика не найдена", http.StatusNotFound)

	}

}

func (mc *HandlerDependencies) updateHandlerJSONValue(w http.ResponseWriter, r *http.Request) {

	var metric Metrics

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
		if metric.MType == "gauge" {
			value, ok := mc.Storage.gauges[metric.ID]
			if ok {
				// Метрика существует в хранилище, используйте значение из хранилища
				createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, value)
				return
			}
		} else if metric.MType == "counter" {
			value, ok := mc.Storage.counters[metric.ID]
			if ok {
				// Метрика существует в хранилище, используйте значение из хранилища
				createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, value)
				return
			}
		}

		// Если метрика отсутствует и в файле, и в хранилище, отправьте статус "Not Found"
		http.Error(w, "Метрика не найдена", http.StatusNotFound)
		return
	}

	// Отправить значение метрики в ответ
	if metric.MType == "gauge" {
		createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, *metricFromFile.Value)
	} else if metric.MType == "counter" {
		createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, *metricFromFile.Delta)
	}
}

func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		CreateLogger()
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

	// Записываем данные JSON в файл
	_, err = file.WriteString(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (mc *HandlerDependencies) ReadMetricsFromFile() (map[string]Metrics, error) {
	metricsMap := make(map[string]Metrics)

	file, err := os.OpenFile(mc.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		mc.Logger.Error("Ошибка при открытии файла для чтения", zap.Error(err))
		return metricsMap, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var metric Metrics
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

func (mc *HandlerDependencies) Ping(w http.ResponseWriter, r *http.Request) {

	// Попытка открыть соединение с базой данных
	db, err := sql.Open("postgres", mc.Config.DatabaseDSN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Если успешно, возвращаем HTTP-статус 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database is working\n")

}

func (mc *HandlerDependencies) SetupDatabase() error {
	println("SetupDatabase")
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", mc.Config.DatabaseDSN)
	if err != nil {
		log.Println("Ошибка при открытии базы данных", err.Error())
		return err
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Println("Ошибка при проверке соединения с базой данных", err.Error())
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
		log.Println("Ошибка при создании таблицы", err.Error())
		return err
	}

	// Возвращаем nil, чтобы показать, что всё прошло успешно

	return nil
}

func (mc *HandlerDependencies) WriteMetricToDatabase(metric Metrics) error {
	var query string
	var args []interface{}

	if metric.MType == "gauge" {
		query = "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3)"
		args = []interface{}{metric.ID, metric.MType, metric.Value}
	} else if metric.MType == "counter" {
		query = "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3)"
		args = []interface{}{metric.ID, metric.MType, metric.Delta}
	} else {
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
