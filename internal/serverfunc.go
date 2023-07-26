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

				if _, err := strconv.ParseFloat(path[4], 64); err == nil {

					c.JSON(http.StatusOK, gin.H{"message": "StatusOK"})
					c.String(http.StatusOK, fmt.Sprintf("%v", num)) // Возвращаем текущее значение метрики в текстовом виде

					storage.SaveMetric(path[2], path[3], num)

					return

				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "StatusBadRequest"})

				}

				if _, err := strconv.ParseInt(path[4], 10, 64); err == nil {
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
				println("lengpath finish", lengpath)
				if lengpath != 4 {
					c.JSON(http.StatusNotFound, gin.H{"error": "StatusNotFound"})

					return

				}

				c.JSON(http.StatusOK, gin.H{"message finish": "StatusOK"})
				storage.SaveMetric(path[2], path[3], num)
				println("path3 SaveMetric", path[3])
				v1 := storage.counters[path[3]]

				c.String(http.StatusOK, fmt.Sprintf("%v", v1))

				return
			}

		}
	}
}
