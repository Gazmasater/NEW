package server

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK
	}
	size, err := sr.ResponseWriter.Write(b)
	sr.size += size
	return size, err
}

func (sr *statusRecorder) Status() int {
	return sr.status
}

func (sr *statusRecorder) Size() int {
	return sr.size
}

func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Создаем экземпляр statusRecorder
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // Устанавливаем начальный статус по умолчанию
		}

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
