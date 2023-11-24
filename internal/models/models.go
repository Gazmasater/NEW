package models

type Metric struct {
	Name  string
	Value any
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type SysMetrics struct {
	TotalMemory    float64
	FreeMemory     float64
	CPUUtilization []float64
}
