// В файле main_test.go

package main

import (
	"testing"
	"time"

	"project.com/internal"
)

func TestCollectMetrics(t *testing.T) {
	// Создаем канал для получения метрик
	metricsChan := internal.CollectMetrics(1*time.Second, "http://internal.Addr/update/gauge/test1/100")

	// Ждем 2 секунды, чтобы метрики собрались
	time.Sleep(2 * time.Second)

	// Получаем метрики из канала
	metrics := <-metricsChan

	// Проверяем, что полученные метрики не пустые
	if len(metrics) == 0 {
		t.Errorf("Expected non-empty metrics, but got empty")
	}
}
