package serverin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type MyController struct {
	deps *HandlerDependencies
}

func NewMyController(deps *HandlerDependencies) *MyController {
	return &MyController{
		deps: deps,
	}
}

func (mc *MyController) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка POST-запроса
	// Используйте mc.deps для доступа к зависимостям
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)
	fmt.Println("http.MethodPost:=", http.MethodPost)

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {
		fmt.Println("lengpath path2=counter", lengpath)
		fmt.Println("metricValue", metricValue)

		if lengpath != 5 {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return

		}

		if metricValue == "none" {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}

		num1, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return
		}

		if isInteger(metricValue) {
			fmt.Println("Num1 в ветке POST ", num1)

			fmt.Fprintf(w, "%v", num1)

			mc.deps.Storage.SaveMetric(metricType, metricName, num1)

			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}
	}
	if lengpath == 4 && metricName == "" {
		http.Error(w, "Metric name not provided", http.StatusBadRequest)
		return
	}

	if (len(metricName) > 0) && (metricValue == "") {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {

		num, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		if _, err1 := strconv.ParseFloat(metricValue, 64); err1 == nil {
			fmt.Fprintf(w, "%v", num) // Возвращаем текущее значение метрики в текстовом виде
			mc.deps.Storage.SaveMetric(path[2], metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)

		}

		if _, err1 := strconv.ParseInt(metricValue, 10, 64); err1 == nil {
			fmt.Fprintf(w, "%v", num) // Возвращаем текущее значение метрики в текстовом виде
			mc.deps.Storage.SaveMetric(metricType, metricName, num)
			return

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

	}

}

func (mc *MyController) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка GET-запроса
	// Используйте mc.deps для доступа к зависимостям
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	path := strings.Split(r.URL.Path, "/")
	lengpath := len(path)
	mc.deps.Logger.Println("http.MethodGet:", http.MethodGet)

	if lengpath != 4 {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	if metricType == "counter" {
		num1, found := mc.deps.Storage.counters[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)

	}
	if metricType == "gauge" {

		num1, found := mc.deps.Storage.gauges[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)

		}

		fmt.Fprintf(w, "%v", num1)
		mc.deps.Logger.Printf("Значение измерителя %s: %v", metricName, num1)

	}

}

func (mc *MyController) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/update", mc.handlePostRequest)
	r.Get("/value", mc.handleGetRequest)

	// Добавьте другие обработчики URL здесь

	return r
}

func StartServer(address string, handler http.Handler) {
	// Создаем HTTP-сервер с настройками
	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	// Запуск HTTP-сервера через http.ListenAndServe()
	fmt.Printf("Запуск HTTP-сервера на адресе: %s\n", address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Ошибка при запуске HTTP-сервера: %s", err)
	}
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s) // Преобразуем строку в целое число, игнорируя результат
	return err == nil         // Если ошибки нет, то строка является целым числом
}
