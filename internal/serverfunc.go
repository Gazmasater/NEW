package internal

import (
	"compress/gzip"

	"net/http"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

func Init() {
	// Инициализация логгера
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger")
	}
	defer logger.Sync() // flushes buffer, if any
}

type Metric struct {
	Name  string
	Value interface{}
}

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (grw GzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}
