package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (mc *app) HandlePostRequestOptimiz(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "counter" {

		num, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		if isInteger(metricValue) {

			w.Write([]byte(strconv.FormatInt(num, 10)))

			mc.Storage.SaveCounter(metricType, metricName, num)

		} else {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return

		}
	}
	if metricName == "" || (len(metricName) > 0 && metricValue == "") {
		http.Error(w, "StatusBadRequest", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {
		num, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "StatusBadRequest", http.StatusBadRequest)
			return
		}

		mc.Storage.SaveGauge(metricType, metricName, num)

		responseData := []byte(strconv.FormatFloat(num, 'f', -1, 64))
		w.Write(responseData)

	}

}

func (mc *app) HandleGetRequestOptimiz(w http.ResponseWriter, r *http.Request) {

	// Обработка GET-запроса
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "StatusNotFound", http.StatusNotFound)
		return
	}

	switch metricType {

	case "counter":

		_, found := mc.Storage.GetCounters()[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return // Add return statement here

		}
		// Handle counter case response here

	case "gauge":

		num, found := mc.Storage.GetGauges()[metricName]
		if !found {
			http.Error(w, "StatusNotFound", http.StatusNotFound)
			return // Add return statement here

		} else {
			// Handle gauge case response here
			w.Write([]byte(strconv.FormatFloat(num, 'f', -1, 64)))
		}
	}
}
