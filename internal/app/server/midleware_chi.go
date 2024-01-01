package server

import (
	"log"
	"net/http"
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
