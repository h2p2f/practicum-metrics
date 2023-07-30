// Package storage реализует хранилище метрик в памяти.
//
// Package storage implements a metric storage in memory.
package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/h2p2f/practicum-metrics/internal/agent/models"
)

// MetricStorage хранит метрики в памяти.
//
// MetricStorage stores metrics in memory.
type MetricStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mut     sync.RWMutex
}

// NewAgentStorage создает новое хранилище метрик.
//
// NewAgentStorage creates a new metric storage.
func NewAgentStorage() *MetricStorage {
	return &MetricStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// RuntimeMetricsMonitor собирает метрики из runtime.
//
// RuntimeMetricsMonitor collects metrics from runtime.
func (m *MetricStorage) RuntimeMetricsMonitor() {
	m.mut.RLock()
	defer m.mut.RUnlock()

	var rtMetrics runtime.MemStats

	runtime.ReadMemStats(&rtMetrics)

	m.gauge["Alloc"] = float64(rtMetrics.Alloc)
	m.gauge["BuckHashSys"] = float64(rtMetrics.BuckHashSys)
	m.gauge["Frees"] = float64(rtMetrics.Frees)
	m.gauge["GCCPUFraction"] = float64(rtMetrics.GCCPUFraction)
	m.gauge["GCSys"] = float64(rtMetrics.GCSys)
	m.gauge["HeapAlloc"] = float64(rtMetrics.HeapAlloc)
	m.gauge["HeapIdle"] = float64(rtMetrics.HeapIdle)
	m.gauge["HeapInuse"] = float64(rtMetrics.HeapInuse)
	m.gauge["HeapObjects"] = float64(rtMetrics.HeapObjects)
	m.gauge["HeapReleased"] = float64(rtMetrics.HeapReleased)
	m.gauge["HeapSys"] = float64(rtMetrics.HeapSys)
	m.gauge["LastGC"] = float64(rtMetrics.LastGC)
	m.gauge["Lookups"] = float64(rtMetrics.Lookups)
	m.gauge["MCacheInuse"] = float64(rtMetrics.MCacheInuse)
	m.gauge["MCacheSys"] = float64(rtMetrics.MCacheSys)
	m.gauge["MSpanInuse"] = float64(rtMetrics.MSpanInuse)
	m.gauge["MSpanSys"] = float64(rtMetrics.MSpanSys)
	m.gauge["Mallocs"] = float64(rtMetrics.Mallocs)
	m.gauge["NextGC"] = float64(rtMetrics.NextGC)
	m.gauge["NumForcedGC"] = float64(rtMetrics.NumForcedGC)
	m.gauge["NumGC"] = float64(rtMetrics.NumGC)
	m.gauge["OtherSys"] = float64(rtMetrics.OtherSys)
	m.gauge["PauseTotalNs"] = float64(rtMetrics.PauseTotalNs)
	m.gauge["StackInuse"] = float64(rtMetrics.StackInuse)
	m.gauge["StackSys"] = float64(rtMetrics.StackSys)
	m.gauge["Sys"] = float64(rtMetrics.Sys)
	m.gauge["TotalAlloc"] = float64(rtMetrics.TotalAlloc)
	m.gauge["RandomValue"] = rand.Float64() * 10000
	m.counter["PollCount"]++
}

// GopsUtilizationMonitor собирает метрики из gopsutil.
//
// GopsUtilizationMonitor collects metrics from gopsutil.
func (m *MetricStorage) GopsUtilizationMonitor() {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}
	cp, err := cpu.Percent(0, false)
	if err != nil {
		log.Println(err)
	}
	m.mut.Lock()
	defer m.mut.Unlock()
	m.gauge["TotalMemory"] = float64(memory.Total)
	m.gauge["FreeMemory"] = float64(memory.Free)
	m.gauge["CPUtilization1"] = cp[0]

}

// URLMetrics генерирует URL для отправки метрик на сервер.
// Необходима для обратной совместимости. В данный момент не используется.
func (m *MetricStorage) URLMetrics(host string) []string {

	m.mut.Lock()
	defer m.mut.Unlock()

	var urls []string

	for metric, value := range m.gauge {
		generatedURL := fmt.Sprintf("%s/update/gauge/%s/%f", host, metric, value)
		urls = append(urls, generatedURL)
	}
	for metric, value := range m.counter {
		generatedURL := fmt.Sprintf("%s/update/counter/%s/%d", host, metric, value)
		urls = append(urls, generatedURL)
	}
	m.counter["PollCount"] = 0
	return urls
}

// JSONMetrics генерирует слайс JSON-объектов для отправки метрик на сервер.
// Необходима для обратной совместимости. В данный момент не используется.
//
// JSONMetrics generates a slice of JSON objects to send metrics to the server.
// Required for backward compatibility. Currently not used.
func (m *MetricStorage) JSONMetrics() [][]byte {
	var res [][]byte
	var model models.Metric
	m.mut.RLock()
	defer m.mut.RUnlock()

	for metric, value := range m.gauge {
		value := value
		model = models.Metric{
			ID:    metric,
			Value: &value,
			MType: "gauge",
		}
		mj, err := json.Marshal(model)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, mj)
	}
	for metric, value := range m.counter {
		value := value
		model = models.Metric{
			ID:    metric,
			Delta: &value,
			MType: "counter",
		}
		mj, err := json.Marshal(model)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, mj)
	}
	m.counter["PollCount"] = 0
	return res
}

// BatchJSONMetrics генерирует JSON-объект для отправки метрик на сервер.
//
// BatchJSONMetrics generates a JSON object to send metrics to the server.
func (m *MetricStorage) BatchJSONMetrics() []byte {
	var res []byte
	var modelSlice []models.Metric
	m.mut.RLock()
	defer m.mut.RUnlock()

	for metric, value := range m.gauge {
		value := value
		mt := models.Metric{
			ID:    metric,
			Value: &value,
			MType: "gauge",
		}
		modelSlice = append(modelSlice, mt)
	}
	for metric, value := range m.counter {
		value := value
		mt := models.Metric{
			ID:    metric,
			Delta: &value,
			MType: "counter",
		}
		modelSlice = append(modelSlice, mt)
	}
	res, err := json.Marshal(modelSlice)
	if err != nil {
		fmt.Println(err)
	}
	m.counter["PollCount"] = 0
	return res
}
