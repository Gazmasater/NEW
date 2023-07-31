package main

import (
	"flag"
	"fmt"
	"net/http"
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

//func parseAddr() (string, error) {
// Определение и парсинг флага
//	var addr string

//	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")

//	fmt.Println("here is address server", addr)

//	return addr, nil
//}

func main() {
	// Определение флагов -a, -r и -p с значениями по умолчанию
	// Вызываем новую функцию для парсинга флага и получения адреса сервера
	//addr, err := parseAddr()
	//if err != nil {
	//	fmt.Println("Ошибка парсинга адреса сервера:", err)
	//	return
	//}

	var (
		reportSeconds int
		pollSeconds   int
		addr          string
	)

	flag.StringVar(&addr, "a", "localhost:8080", "Адрес HTTP-сервера")

	fmt.Println("here is address agent", addr)

	flag.IntVar(&reportSeconds, "r", 10, "Частота отправки метрик на сервер (в секундах)")
	flag.IntVar(&pollSeconds, "p", 2, "Частота опроса метрик из пакета runtime (в секундах)")

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
