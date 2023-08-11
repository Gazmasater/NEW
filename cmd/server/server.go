package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"project.com/internal/serverin"
)

func main() {
	// Инициализируем конфигурацию сервера
	serverCfg := serverin.InitServerConfig()

	storage := serverin.NewMemStorage()
	logger := serverin.NewLogger()

	deps := serverin.NewHandlerDependencies(storage, logger)

	r := newRouter(deps)

	serverin.StartServer(serverCfg.Address, r)
}

func UpdateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/update/") {
			http.Error(w, "StatusBadRequest no update", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValueMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/value/") {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newRouter(deps *serverin.HandlerDependencies) http.Handler {
	r := chi.NewRouter()

	r.Get("/metrics", serverin.HandleMetrics(deps))

	// Монтирование подмаршрута /update
	updateRouter := chi.NewRouter()
	updateRouter.Use(UpdateMiddleware)
	updateRouter.Post("/{metricType}/{metricName}/{metricValue}", func(w http.ResponseWriter, r *http.Request) {
		serverin.HandlePostRequest(w, r, deps)
	})
	r.Mount("/update", updateRouter)

	// Монтирование подмаршрута /value
	valueRouter := chi.NewRouter()
	valueRouter.Use(ValueMiddleware)
	valueRouter.Get("/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
		serverin.HandleGetRequest(w, r, deps)
	})
	r.Mount("/value", valueRouter)

	return r
}
