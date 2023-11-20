package collector

import (
	"math/rand"
	"reflect"
	"runtime"
)

type MetricField func(stats *runtime.MemStats) float64

var MetricFieldMap = map[string]MetricField{}

// generateFieldFunc generates a function for the specified field in MemStats
func generateFieldFunc(fieldName string) MetricField {
	return func(stats *runtime.MemStats) float64 {
		val := reflect.ValueOf(*stats)
		field := val.FieldByName(fieldName)
		return float64(reflect.Indirect(field).Uint()) // Assuming the fields are uintegers, adjust as needed
	}
}

// initializeMetricFieldMap initializes the MetricFieldMap with field names and corresponding functions
func initializeMetricFieldMap() {
	fields := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
		"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys",
		"LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys",
		"Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs",
		"StackInuse", "StackSys", "Sys", "TotalAlloc",
	}

	for _, field := range fields {
		MetricFieldMap[field] = generateFieldFunc(field)
	}

	// Adding RandomValue to MetricFieldMap
	MetricFieldMap["RandomValue"] = func(_ *runtime.MemStats) float64 { return rand.Float64() }
}

func Init() {
	initializeMetricFieldMap()
}
