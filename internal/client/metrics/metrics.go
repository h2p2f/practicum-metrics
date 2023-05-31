package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math/rand"
	"runtime"
	"sync"
)

// JsonMetrics is a struct that contains all the metrics that are being monitored
type JSONMetrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// RuntimeMetrics is a struct that contains all the metrics that are being monitored
type RuntimeMetrics struct {
	mut     sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

// NewMetrics is a function that returns a map of metrics and their values
func (m *RuntimeMetrics) NewMetrics() {
	m.gauge = make(map[string]float64)
	m.counter = make(map[string]int64)
	gaugeMetrics := []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
		"TotalMemory",
		"FreeMemory",
		"CPUutilization1",
	}

	counterMetrics := []string{"PollCount"}
	//initialize metrics
	for _, metric := range gaugeMetrics {
		m.gauge[metric] = 0
	}
	for _, metric := range counterMetrics {
		m.counter[metric] = 0
	}
}

// RuntimeMonitor is a function that monitors the metrics
func (m *RuntimeMetrics) RuntimeMonitor() {
	//set up runtime metrics
	var RtMetrics runtime.MemStats

	//get runtime metrics
	runtime.ReadMemStats(&RtMetrics)
	//lock the mutex, update the metrics and unlock the mutex
	m.mut.Lock()
	defer m.mut.Unlock()
	m.gauge["Alloc"] = float64(RtMetrics.Alloc)
	m.gauge["BuckHashSys"] = float64(RtMetrics.BuckHashSys)
	m.gauge["Frees"] = float64(RtMetrics.Frees)
	m.gauge["GCCPUFraction"] = float64(RtMetrics.GCCPUFraction)
	m.gauge["GCSys"] = float64(RtMetrics.GCSys)
	m.gauge["HeapAlloc"] = float64(RtMetrics.HeapAlloc)
	m.gauge["HeapIdle"] = float64(RtMetrics.HeapIdle)
	m.gauge["HeapInuse"] = float64(RtMetrics.HeapInuse)
	m.gauge["HeapObjects"] = float64(RtMetrics.HeapObjects)
	m.gauge["HeapReleased"] = float64(RtMetrics.HeapReleased)
	m.gauge["HeapSys"] = float64(RtMetrics.HeapSys)
	m.gauge["LastGC"] = float64(RtMetrics.LastGC)
	m.gauge["Lookups"] = float64(RtMetrics.Lookups)
	m.gauge["MCacheInuse"] = float64(RtMetrics.MCacheInuse)
	m.gauge["MCacheSys"] = float64(RtMetrics.MCacheSys)
	m.gauge["MSpanInuse"] = float64(RtMetrics.MSpanInuse)
	m.gauge["MSpanSys"] = float64(RtMetrics.MSpanSys)
	m.gauge["Mallocs"] = float64(RtMetrics.Mallocs)
	m.gauge["NextGC"] = float64(RtMetrics.NextGC)
	m.gauge["NumForcedGC"] = float64(RtMetrics.NumForcedGC)
	m.gauge["NumGC"] = float64(RtMetrics.NumGC)
	m.gauge["OtherSys"] = float64(RtMetrics.OtherSys)
	m.gauge["PauseTotalNs"] = float64(RtMetrics.PauseTotalNs)
	m.gauge["StackInuse"] = float64(RtMetrics.StackInuse)
	m.gauge["StackSys"] = float64(RtMetrics.StackSys)
	m.gauge["Sys"] = float64(RtMetrics.Sys)
	m.gauge["TotalAlloc"] = float64(RtMetrics.TotalAlloc)
	m.counter["PollCount"]++
	m.gauge["RandomValue"] = rand.Float64() * 10000

}

// GopsUtilMonitor is a function that monitors the metrics
func (m *RuntimeMetrics) GopsUtilMonitor() {
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
	m.gauge["CPUutilization1"] = cp[0]

}

// URLMetrics  is a function that returns a slice of urls that are generated from the metrics and their values
func (m *RuntimeMetrics) URLMetrics(host string) []string {
	//lock the mutex
	m.mut.Lock()
	defer m.mut.Unlock()
	//create a slice of urls
	var urls []string
	//generate urls
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

// JSONMetrics  is a function that returns a slice of marshalled all data
// that are generated from the metrics and their values
func (m *RuntimeMetrics) JSONMetrics() []byte {
	//lock the mutex
	m.mut.Lock()
	defer m.mut.Unlock()

	var metrics []JSONMetrics
	for metric, value := range m.gauge {
		value := value
		jsonMetric := JSONMetrics{ID: metric, MType: "gauge", Value: &value}
		metrics = append(metrics, jsonMetric)
	}
	for metric, value := range m.counter {
		value := value
		jsonMetric := JSONMetrics{ID: metric, MType: "counter", Delta: &value}
		metrics = append(metrics, jsonMetric)
	}
	out, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	m.counter["PollCount"] = 0
	return out
}

// JSONMetricsForSingleSending  is a function that returns a slice of
// marshalled data that are generated from the metrics and their values
func (m *RuntimeMetrics) JSONMetricsForSingleSending() [][]byte {
	//lock the mutex
	m.mut.Lock()
	defer m.mut.Unlock()
	//create a slice of urls
	var result [][]byte
	//generate urls
	for metric, value := range m.gauge {
		value := value
		jsonMetric := JSONMetrics{ID: metric, MType: "gauge", Value: &value}
		out, err := json.Marshal(jsonMetric)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}
	for metric, value := range m.counter {
		value := value
		jsonMetric := JSONMetrics{ID: metric, MType: "counter", Delta: &value}
		out, err := json.Marshal(jsonMetric)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}
	m.counter["PollCount"] = 0
	return result
}
