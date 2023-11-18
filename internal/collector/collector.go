package collector

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"project.com/internal/models"
)

func CollectMetrics(pollInterval time.Duration, serverURL string) <-chan []*models.Metrics {
	metricsChan := make(chan []*models.Metrics)
	var pollCount int64 = 0
	println("CollectMetrics serverURL", serverURL)
	var memStats runtime.MemStats

	go func() {

		for {
			metrics := make([]*models.Metrics, 0)

			runtime.ReadMemStats(&memStats)

			addMetric(&metrics, "Alloc", float64(memStats.Alloc))

			addMetric(&metrics, "BuckHashSys", float64(memStats.Alloc))

			addMetric(&metrics, "Frees", float64(memStats.Alloc))

			addMetric(&metrics, "GCCPUFraction", float64(memStats.Alloc))
			gCSysValue := float64(memStats.GCSys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "GCSys", Value: &gCSysValue})

			heapAllocValue := float64(memStats.HeapAlloc)
			heapAllocValue += rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "HeapAlloc", Value: &heapAllocValue})

			heapIdleValue := float64(memStats.HeapIdle)
			heapIdleValue += rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "HeapIdle", Value: &heapIdleValue})

			heapInuseValue := float64(memStats.HeapInuse)
			heapInuseValue += rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "HeapInuse", Value: &heapInuseValue})

			heapObjectsValue := float64(memStats.HeapObjects)
			heapObjectsValue += rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "HeapObjects", Value: &heapObjectsValue})

			heapReleasedValue := float64(memStats.HeapReleased)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "HeapReleased", Value: &heapReleasedValue})

			heapSysValue := float64(memStats.HeapSys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "HeapSys", Value: &heapSysValue})

			lastGCValue := float64(memStats.LastGC)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "LastGC", Value: &lastGCValue})

			lookupsValue := float64(memStats.Lookups)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "Lookups", Value: &lookupsValue})

			mCacheInuseValue := float64(memStats.MCacheInuse)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "MCacheInuse", Value: &mCacheInuseValue})

			mCacheSysValue := float64(memStats.MCacheSys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "MCacheSys", Value: &mCacheSysValue})

			mSpanInuseValue := float64(memStats.MSpanInuse)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "MSpanInuse", Value: &mSpanInuseValue})

			mSpanSysValue := float64(memStats.MSpanSys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "MSpanSys", Value: &mSpanSysValue})

			mallocsValue := float64(memStats.Mallocs)
			mallocsValue += rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "Mallocs", Value: &mallocsValue})

			nextGCValue := float64(memStats.NextGC)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "NextGC", Value: &nextGCValue})

			numForcedGCValue := float64(memStats.NumForcedGC)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "NumForcedGC", Value: &numForcedGCValue})

			numGCValue := float64(memStats.NumGC)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "NumGC", Value: &numGCValue})

			otherSysValue := float64(memStats.OtherSys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "OtherSys", Value: &otherSysValue})

			pauseTotalNsValue := float64(memStats.PauseTotalNs)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "PauseTotalNs", Value: &pauseTotalNsValue})

			stackInuseValue := float64(memStats.StackInuse)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "StackInuse", Value: &stackInuseValue})

			stackSysValue := float64(memStats.StackSys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "StackSys", Value: &stackSysValue})

			sysValue := float64(memStats.Sys)
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "Sys", Value: &sysValue})

			totalAllocValue := float64(memStats.TotalAlloc)
			totalAllocValue += rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "TotalAlloc", Value: &totalAllocValue})

			randomValue := rand.Float64()
			metrics = append(metrics, &models.Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValue})

			// Добавляем метрику PollCount типа counter!!
			metrics = append(metrics, &models.Metrics{MType: "counter", ID: "PollCount", Delta: &pollCount})

			//  Увеличиваем счетчик обновлений метр!!!
			pollCount++

			metricsChan <- metrics
			time.Sleep(pollInterval)

		}
	}()

	return metricsChan
}

func CollectAdditionalMetrics() (float64, float64, []float64) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Ошибка при получении информации о памяти:", err)

		return 0, 0, nil
	}

	cpuInfo, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Ошибка при получении информации о CPU:", err)
		return 0, 0, nil
	}

	// Преобразуйте vmStat.Total в float64
	totalMemory := float64(vmStat.Total)

	return totalMemory, float64(vmStat.Free), cpuInfo
}

func addMetric(metrics *[]*models.Metrics, id string, value float64) {
	*metrics = append(*metrics, &models.Metrics{
		MType: "gauge",
		ID:    id,
		Value: &value,
	})
}
