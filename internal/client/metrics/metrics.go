package metrics

import (
	"fmt"
	"net/url"
	"runtime"
)

type RuntimeMetrics struct {
	gauge   map[string]float64
	counter map[string]int64
}

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
		"TotalAlloc"}

	counterMetrics := []string{"Counter"}
	for _, metric := range gaugeMetrics {
		m.gauge[metric] = 0
	}
	for _, metric := range counterMetrics {
		m.counter[metric] = 0
	}
}
func (m *RuntimeMetrics) Monitor() {
	//set up runtime metrics
	var RtMetrics runtime.MemStats
	//get runtime metrics
	runtime.ReadMemStats(&RtMetrics)
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
	m.counter["Counter"]++
}
func (m *RuntimeMetrics) UrlMetrics(host string) []string {
	var urls []string
	_, err := url.ParseRequestURI(host)
	if err != nil {
		return nil
	}
	for metric, value := range m.gauge {
		generatedUrl := fmt.Sprintf("%s/update/gauge/%s/%f", host, metric, value)
		urls = append(urls, generatedUrl)
	}
	for metric, value := range m.counter {
		generatedUrl := fmt.Sprintf("%s/update/counter/%s/%d", host, metric, value)
		urls = append(urls, generatedUrl)
	}
	return urls
}
