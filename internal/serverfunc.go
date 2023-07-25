package internal

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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

					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})

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

				if _, err := strconv.ParseFloat(path[4], 64); err == nil {

					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

				}

				if _, err := strconv.ParseInt(path[4], 10, 64); err == nil {
					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})
					return
				}

			}

			c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
			//================================================================================
		case http.MethodGet:

			num1, err := strconv.ParseFloat(path[1], 64)

			if (err != nil) && (lengpath == 2) {

				allMetrics := storage.GetAllMetrics()

				// Form an HTML page with the list of all metrics and their values
				html := "<html><head><title>Metrics</title></head><body><h1>Metrics List</h1><ul>"
				for name, value := range allMetrics {
					html += fmt.Sprintf("<li>%s: %v</li>", name, value)
				}
				html += "</ul></body></html>"

				c.Header("Content-Type", "text/html; charset=utf-8")
				c.String(http.StatusOK, html)

				return
			}

			if err != nil {

				c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

				//fmt.Println("Ошибка при преобразовании строки во float64:", err)
				return
			}
			if (path[2] != "gauge") && (path[2] != "counter") {
				c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

				return
			}

			if path[2] == "counter" {
				println("path2==counter", path[2])

				num, err := strconv.ParseInt(path[1], 10, 64)
				println("NUM ERR", num, err)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

					return
				}
				println("path3 counter strconv.ParseFloat(path[3], 64)", path[3])
				_, err1 := strconv.ParseFloat(path[3], 64)
				if err1 == nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

				}

				if lengpath != 4 {
					c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

					return

				}

				c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})

				storage.SaveMetric(path[2], path[3], num1)

				return
			}

		}
	}
}
