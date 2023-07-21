// В файле main_test.go

package main

import (
	"testing"
	"time"
)

func TestCollectMetrics(t *testing.T) {
	// Создаем канал для получения метрик
	metricsChan := collectMetrics(1*time.Second, "http://localhost:8080/update/gauge/test1/100")

	// Ждем 2 секунды, чтобы метрики собрались
	time.Sleep(2 * time.Second)

	// Получаем метрики из канала
	metrics := <-metricsChan

	// Проверяем, что полученные метрики не пустые
	if len(metrics) == 0 {
		t.Errorf("Expected non-empty metrics, but got empty")
	}
}
