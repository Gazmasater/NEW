package server

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Создаем обертку для записи статуса ответа и размера
		recorder := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Передаем запрос следующему обработчику
		next.ServeHTTP(recorder, r)

		// Записываем информацию о запросе
		log.Printf(
			"[REQUEST] Method: %s, URI: %s, Duration: %v",
			r.Method,
			r.RequestURI,
			time.Since(startTime),
		)

		// Записываем информацию о ответе
		log.Printf(
			"[RESPONSE] Status: %d, Size: %d bytes",
			recorder.Status(),
			recorder.BytesWritten(),
		)
	})
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
