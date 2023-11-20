package collector

import (
	"math/rand"
	"reflect"
	"runtime"
)

type MetricField func(stats *runtime.MemStats) float64

func getMetricField(fieldName string) MetricField {
	switch fieldName {
	case "RandomValue":
		return func(_ *runtime.MemStats) float64 { return rand.Float64() }
	case "GCCPUFraction":
		return func(stats *runtime.MemStats) float64 { return stats.GCCPUFraction }
	default:
		return func(stats *runtime.MemStats) float64 {
			val := reflect.ValueOf(stats).Elem().FieldByName(fieldName)
			if val.IsValid() {
				return float64(val.Uint())
			}
			return 0
		}
	}
}

func createMetricFieldMap(fieldNames ...string) map[string]MetricField {
	fieldMap := make(map[string]MetricField)
	for _, fieldName := range fieldNames {
		fieldMap[fieldName] = getMetricField(fieldName)
	}
	return fieldMap
}

var MetricFieldMap = createMetricFieldMap(
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
	"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
	"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "RandomValue",
)
