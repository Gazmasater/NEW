package internal

import (
	"compress/gzip"

	"net/http"

	_ "github.com/lib/pq"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (grw GzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}
