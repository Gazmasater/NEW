package internal

import "sync"

type MemStorage struct {
	metrics map[string]interface{}
	mu      sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]interface{}),
	}
}

func (ms *MemStorage) SaveMetric(metricType, metricName string, metricValue interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Проверяем, есть ли уже запись для данного типа метрики
	key := metricType + ":" + metricName
	prevValue, ok := ms.metrics[key]

	switch metricValue.(type) {
	case float64:
		// Если тип метрики - gauge (float64), замещаем предыдущее знач
		ms.metrics[key] = metricValue
	case int64:
		// Если тип метрики - counter (int64), добавляем новое значение к предыдущему
		if ok {
			if prevCounter, ok := prevValue.(int64); ok {
				newCounter := prevCounter + metricValue.(int64)
				ms.metrics[key] = newCounter
			} else {
				// Если предыдущее значение не является типом int64, замещаем его новым значением
				ms.metrics[key] = metricValue
			}
		} else {
			ms.metrics[key] = metricValue
		}
	}
}

func (ms *MemStorage) GetMetric(metricType, metricName string) (interface{}, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Проверяем, есть ли запись для данного типа метрики
	key := metricType + ":" + metricName
	value, ok := ms.metrics[key]
	return value, ok
}
func (ms *MemStorage) ProcessMetrics(metricType, metricName string, metricValue interface{}) {
	// Сохраняем метрику в хранилище
	ms.SaveMetric(metricType, metricName, metricValue)
}
