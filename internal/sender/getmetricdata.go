package sender

import "project.com/internal/models"

func getMetricData(metric models.Metrics) any {
	data := map[string]any{
		"type": metric.MType,
		"id":   metric.ID,
	}

	if metric.MType == "counter" {
		data["delta"] = *metric.Delta
	} else {
		data["value"] = *metric.Value
	}

	return data
}
