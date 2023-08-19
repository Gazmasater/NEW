package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"project.com/internal/serverin"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	serverin.InitLogger()

	defer logger.Sync()

	serverin.Sugar.Info("Initializing logger...") // Используем sugar для логирования

	// Инициализация конфигурации и хранилища
	config := serverin.InitServerConfig(logger)
	storage := serverin.NewMemStorage()

	serverin.Sugar.Info("Initializing configuration and storage...")
	// Создание контроллера
	controller := serverin.NewHandlerDependencies(storage, logger)

	// Создание маршрутизатора
	r := chi.NewRouter()

	// Middleware для логирования
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("HTTP Method:", zap.String("method", r.Method))
			next.ServeHTTP(w, r)
		})
	}

	// Добавляем middleware к маршруту
	r.Route("/", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return serverin.WithLogging(next, serverin.Sugar)
		})
		r.Use(loggingMiddleware) // Middleware для логирования
		r.Mount("/", controller.Route())
	})

	// Запуск сервера
	serverin.StartServer(config.Address, r)

	serverin.Sugar.Info("Server started on address:", config.Address)
}
