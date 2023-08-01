package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"project.com/internal"
)

func sendDataToServer(metrics []*internal.Metric, serverURL string) {

	for _, metric := range metrics {
		serverURL := fmt.Sprintf("http://%s/update/%s/%s/%v", serverURL, metric.Type, metric.Name, metric.Value)
		println("serverURL sendDataToServer  ", serverURL)
		//Отправка POST-запроса
		resp, err := http.Post(serverURL, "text/plain", nil)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса:", err)
			return
		}
		defer resp.Body.Close()

	}
}

func main() {

	var (
		reportSeconds int
		pollSeconds   int
		addr          string
	)

	// Чтение переменных окружения или установка значений по умолчанию
	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		addr = addrEnv
	} else {
		flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")
		if _, err := url.Parse(addr); err != nil {
			fmt.Printf("Ошибка: неверный формат адреса сервера в модуле агента: %s\n", addr)
			return
		}
	}
	// Проверка валидности адреса

	reportSecondsEnv := os.Getenv("REPORT_INTERVAL")
	if reportSecondsEnv != "" {
		reportSeconds, _ = strconv.Atoi(reportSecondsEnv)
	} else {
		flag.IntVar(&reportSeconds, "r", 10, "Частота отправки метрик на сервер (в секундах)")
		if reportSeconds <= 0 {
			fmt.Println("Частота отправки метрик должна быть положительным числом.")
			flag.Usage()
			os.Exit(1)
		}

	}

	pollSecondsEnv := os.Getenv("POLL_INTERVAL")
	if pollSecondsEnv != "" {
		pollSeconds, _ = strconv.Atoi(pollSecondsEnv)
	} else {
		flag.IntVar(&pollSeconds, "p", 2, "Частота опроса метрик из пакета runtime (в секундах)")
		if pollSeconds <= 0 {
			fmt.Println("Частота опроса метрик должна быть положительным числом.")
			flag.Usage()
			os.Exit(1)
		}
	}

	flag.Parse()

	pollInterval := time.Duration(pollSeconds) * time.Second
	reportInterval := time.Duration(reportSeconds) * time.Second

	metricsChan := internal.CollectMetrics(pollInterval, addr)

	// Горутина  отправки метрик на сервер с интервалом в reportInterval секунд
	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			sendDataToServer(metrics, addr)
		}
	}()

	for range time.Tick(pollInterval) {
		fmt.Println("Сбор метрик...")
	}
}
