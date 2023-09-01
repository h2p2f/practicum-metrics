package storage

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math/rand"
	"runtime"
)

// RuntimeMetricsMonitor is a method of the MetricStorage structure that collects metrics from runtime.
func (m *MetricStorage) RuntimeMetricsMonitor() {
	m.mut.Lock()
	defer m.mut.Unlock()

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

// GopsUtilizationMonitor is a method of the MetricStorage structure that collects metrics from gopsutil.
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
