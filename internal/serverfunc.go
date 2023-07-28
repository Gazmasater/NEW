package internal

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []*Metric {
	metricsChan := make(chan []*Metric)

	// Переменная для счетчика обновлений метрик
	pollCount := 0

	var metrics []*Metric
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	go func() {
		for {

			metrics = append(metrics, &Metric{Type: "gauge", Name: "Alloc", Value: float64(memStats.Alloc)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "BuckHashSys", Value: float64(memStats.BuckHashSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Frees", Value: float64(memStats.Frees)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "GCCPUFraction", Value: float64(memStats.GCCPUFraction)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "GCSys", Value: float64(memStats.GCSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapAlloc", Value: float64(memStats.HeapAlloc)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapIdle", Value: float64(memStats.HeapIdle)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapInuse", Value: float64(memStats.HeapInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapObjects", Value: float64(memStats.HeapObjects)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapReleased", Value: float64(memStats.HeapReleased)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "HeapSys", Value: float64(memStats.HeapSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "LastGC", Value: float64(memStats.LastGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Lookups", Value: float64(memStats.Lookups)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MCacheInuse", Value: float64(memStats.MCacheInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MCacheSys", Value: float64(memStats.MCacheSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MSpanInuse", Value: float64(memStats.MSpanInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "MSpanSys", Value: float64(memStats.MSpanSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Mallocs", Value: float64(memStats.Mallocs)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "NextGC", Value: float64(memStats.NextGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "NumForcedGC", Value: float64(memStats.NumForcedGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "NumGC", Value: float64(memStats.NumGC)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "OtherSys", Value: float64(memStats.OtherSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "PauseTotalNs", Value: float64(memStats.PauseTotalNs)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "StackInuse", Value: float64(memStats.StackInuse)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "StackSys", Value: float64(memStats.StackSys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "Sys", Value: float64(memStats.Sys)})
			metrics = append(metrics, &Metric{Type: "gauge", Name: "TotalAlloc", Value: float64(memStats.TotalAlloc)})

			// Добавляем метрику RandomValue типа gauge с произвольным значением
			randomValue := rand.Float64()
			metrics = append(metrics, &Metric{Type: "gauge", Name: "RandomValue", Value: randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &Metric{Type: "counter", Name: "PollCount", Value: pollCount})

			// Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)
		}
	}()

	return metricsChan
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func HandleUpdate(storage *MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		println("http.Method:=", c.Request.Method)
		path := strings.Split(c.Request.URL.Path, "/")
		lengpath := len(path)
		println("LENGTH", lengpath)
		// Обрабатываем полученные метрики
		// Преобразование строки во float64

		switch c.Request.Method {
		//==========================================================================================
		case http.MethodPost:
			println("http.MethodPost:=", http.MethodPost)

			if path[1] != "update" {

				c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest no update"})

				return
			}

			if path[2] != "gauge" && path[2] != "counter" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

				return
			}

			if path[2] == "counter" {
				println("lengpath path2=counter", lengpath)
				println("path[4]", path[4])

				if lengpath != 5 {
					c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

					return

				}

				if path[4] == "none" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

					return

				}

				num1, err := strconv.ParseInt(path[4], 10, 64)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

					return
				}

				if isInteger(path[4]) {
					fmt.Println("Num1 в ветке POST ", num1)

					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
					c.String(http.StatusOK, fmt.Sprintf("%v", num1)) // Возвращаем текущее значение метрики в текстовом виде

					storage.SaveMetric(path[2], path[3], num1)

					return

				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})
					return

				}
			}
			if lengpath == 4 && path[3] == "" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Metric name not provided"})

				return
			}

			if (len(path[3]) > 0) && (path[4] == "") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

				return
			}

			if path[2] == "gauge" {

				num, err := strconv.ParseFloat(path[4], 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

					return
				}

				if _, err1 := strconv.ParseFloat(path[4], 64); err1 == nil {

					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
					c.String(http.StatusOK, fmt.Sprintf("%v", num)) // Возвращаем текущее значение метрики в текстовом виде

					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

				}

				if _, err1 := strconv.ParseInt(path[4], 10, 64); err1 == nil {
					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
					c.String(http.StatusOK, fmt.Sprintf("%v", num)) // Возвращаем текущее значение метрики в текстовом виде

					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})
					return
				}

			}

			//================================================================================
		case http.MethodGet:
			println("http.MethodGet", http.MethodGet)

			if lengpath != 4 {
				c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})
				return
			}

			if path[1] != "value" {
				c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})
				return
			}
			if path[2] != "gauge" && path[2] != "counter" {
				c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

				return
			}

			if path[2] == "counter" {
				num1 := storage.counters[path[3]]

				c.String(http.StatusOK, fmt.Sprintf("%v", num1))

			}
			if path[2] == "gauge" {
				num1 := storage.gauges[path[3]]

				c.String(http.StatusOK, fmt.Sprintf("%v", num1))

			}

		}
	}
}
