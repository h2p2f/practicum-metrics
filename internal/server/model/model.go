package model

import (
	"encoding/json"
	"fmt"
	"log"
)

// MemStorage is a model in memory
// it is a struct with two maps - gauges and counters

type MemStorage struct {
	//mutex is suspended work with files TODO: fix it
	//mut      sync.RWMutex
	Gauges   map[string][]float64
	Counters map[string]int64
}

type metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// NewMemStorage creates a new MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string][]float64),
		Counters: make(map[string]int64),
	}
}

// NewMetricsCounter is a constructor function that returns a new counter struct
func NewMetricsCounter(ID, MType string, delta int64) *metrics {
	return &metrics{
		ID:    ID,
		MType: MType,
		Delta: &delta,
	}
}

// NewMetricsGauge is a constructor function that returns a new gauge struct
func NewMetricsGauge(ID, MType string, value float64) *metrics {
	return &metrics{
		ID:    ID,
		MType: MType,
		Value: &value,
	}
}

// SetGauge sets or adds a gauge value
func (m *MemStorage) SetGauge(name string, value float64) {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	m.Gauges[name] = append(m.Gauges[name], value)
}

// SetCounter sets a counter value
func (m *MemStorage) SetCounter(name string, value int64) {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	m.Counters[name] = value
}

// GetGauge gets a gauge value
func (m *MemStorage) GetGauge(name string) ([]float64, bool) {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	value, ok := m.Gauges[name]
	return value, ok
}

// GetCounter gets a counter value
func (m *MemStorage) GetCounter(name string) (int64, bool) {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	value, ok := m.Counters[name]
	return value, ok
}

// GetAllGauges gets all gauges
func (m *MemStorage) GetAllGauges() map[string][]float64 {
	return m.Gauges
}

// GetAllCounters gets all counters
func (m *MemStorage) GetAllCounters() map[string]int64 {
	return m.Counters
}

// GetAllMetricsSliced gets all metrics
// this code is not beautiful
// but i received a some bug with pointers
// if i put counters value to Metrics struct directly
// i receive the same pointer for all counters
// so i implemented via constructor
func (m *MemStorage) GetAllMetricsSliced() []metrics {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	var metrics []metrics
	for key, value := range m.GetAllCounters() {
		met := NewMetricsCounter(key, "counter", value)
		//met.Delta = &value
		fmt.Println(key, value)
		fmt.Println(met)
		metrics = append(metrics, *met)

	}
	for key, value := range m.Gauges {
		met := NewMetricsGauge(key, "gauge", value[len(value)-1])
		//met.Value = &value[len(value)-1]
		metrics = append(metrics, *met)

	}
	return metrics
}

func (m *MemStorage) GetAllInBytesSliced() [][]byte {
	var result [][]byte

	for metric, value := range m.GetAllGauges() {
		met := metrics{ID: metric, MType: "gauge", Value: &value[len(value)-1]}
		out, err := json.Marshal(met)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}
	for metric, value := range m.GetAllCounters() {
		met := metrics{ID: metric, MType: "counter", Delta: &value}
		out, err := json.Marshal(met)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}
	return result
}

// RestoreMetrics restores metrics from slice
func (m *MemStorage) RestoreMetrics(metrics []metrics) {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	for _, metric := range metrics {
		switch metric.MType {
		case "counter":
			{
				m.SetCounter(metric.ID, *metric.Delta)
			}
		case "gauge":
			{
				m.SetGauge(metric.ID, *metric.Value)
			}
		}
	}
}

func (m *MemStorage) RestoreMetric(data [][]byte) {
	var met metrics
	for _, value := range data {
		err := json.Unmarshal(value, &met)
		if err != nil {
			log.Fatal(err)
		}
		switch met.MType {
		case "counter":
			{
				m.SetCounter(met.ID, *met.Delta)
			}
		case "gauge":
			{
				m.SetGauge(met.ID, *met.Value)
			}
		}
	}
}
