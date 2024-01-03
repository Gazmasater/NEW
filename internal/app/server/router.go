package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (mc *app) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return KeyHashMiddleware("MyKey", next)
	})
	r.Use(GzipMiddleware)

	r.Use(LoggerMiddleware)
	r.Post("/update/", mc.updateHandlerJSON)
	r.Post("/updates/", mc.MetricsHandler)

	r.Post("/value/", mc.updateHandlerJSONValue)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", mc.HandlePostRequest)

	r.Post("/value/{metricType}/{metricName}", mc.HandleGetRequestOptimiz)

	r.Get("/value/{metricType}/{metricName}", mc.HandleGetRequestOptimiz)

	r.Get("/metrics", mc.HandleGetRequest)

	r.Get("/", mc.HandleGetRequestHTML)

	r.Get("/ping", mc.Ping)

	return r
}
