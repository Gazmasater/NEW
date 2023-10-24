package internal

import (
	"compress/gzip"

	"fmt"

	"net/http"
	"strings"

	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func (mc *HandlerDependencies) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(GzipMiddleware)

	r.Get("/", mc.HandleGetRequestHTML)

	return r
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
