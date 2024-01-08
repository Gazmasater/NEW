package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"project.com/internal/models"
)

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

func (mc *app) HandleGetRequestOptimiz(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	// Обработка GET-запроса
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	switch metricType {
	case "counter":
		num, found := mc.Storage.GetCounters()[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}
		if contentType == "application/json" {
			mc.createAndSendUpdatedMetricCounterJSON(w, metricName, metricType, int64(num))
			return
		} else {

			w.Write([]byte(strconv.FormatInt(num, 10)))
		}
		return

	case "gauge":

		num, found := mc.Storage.GetGauges()[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}
		if contentType == "application/json" {
			mc.createAndSendUpdatedMetricJSON(w, metricName, metricType, float64(num))
			return
		} else {

			w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))
		}

	}

}

func (mc *app) updateHandlerJSONOptimiz(w http.ResponseWriter, r *http.Request) {

	var metric models.Metrics

	metricsFromFile := make(map[string]models.Metrics)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&metric); err != nil {
		_ = fmt.Errorf("ошибка при разборе JSON: %w", err)
		return
	}

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

		mc.Storage.GetCounters()[metric.ID] = *currentValue.Delta
		metricsFromFile[metric.ID] = currentValue

	}

	if metric.MType == "gauge" && metric.Value != nil {
		// Обновляем или создаем метрику в слайсе
		metricsFromFile[metric.ID] = metric

		// Сохраняем обновленные метрики в хранилище
		mc.Storage.GetGauges()[metric.ID] = *metric.Value
	}

	// Запись обновленных метрик в базу
	for _, updatedMetric := range metricsFromFile {
		dbErr := mc.WriteMetricToDatabaseOptimiz(updatedMetric)
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
			mc.createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, *updatedMetric.Value)

		default:
			http.Error(w, "Метрика не найдена", http.StatusNotFound)

		}
	}

}

func (mc *app) updateHandlerJSONValueOptimiz(w http.ResponseWriter, r *http.Request) {
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
				mc.createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, value)
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
		mc.createAndSendUpdatedMetricJSON(w, metric.ID, metric.MType, *metricFromFile.Value)
	case "counter":
		mc.createAndSendUpdatedMetricCounterJSON(w, metric.ID, metric.MType, *metricFromFile.Delta)

	}
}

func (mc *app) WriteMetricToDatabaseOptimiz(metric models.Metrics) error {
	var query string
	var args []interface{}

	switch metric.MType {
	case "gauge":
		query = "INSERT INTO metrics (name, type, value) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = excluded.type, value = excluded.value"
		args = []interface{}{metric.ID, metric.MType, metric.Value}
	case "counter":
		query = "INSERT INTO metrics (name, type, delta) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET type = excluded.type, delta = excluded.delta"
		args = []interface{}{metric.ID, metric.MType, metric.Delta}
	default:
		log.Printf("Неизвестный тип метрики: %s", metric.MType)
		return fmt.Errorf("неизвестный тип метрики")
	}

	if mc.DB == nil {
		log.Println("Ошибка: mc.DB не инициализирован.")
		return fmt.Errorf("mc.DB не инициализирован")
	}

	// Теперь выполняем вставку новой метрики или обновление существующей
	_, err := mc.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Ошибка при записи метрики в базу данных: %s", err)
		return err
	}
	return nil
}
