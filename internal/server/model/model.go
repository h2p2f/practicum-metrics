package model

import (
	"encoding/json"
	"log"
	"sync"
)

// MemStorage is a model in memory
// it is a struct with two maps - gauges and counters

type MemStorage struct {
	//mutex is suspended work with files TODO: fix it
	mut      sync.RWMutex
	Gauges   map[string]float64
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
		mut:      sync.RWMutex{},
		Gauges:   make(map[string]float64),
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
	m.mut.Lock()
	defer m.mut.Unlock()
	m.Gauges[name] = value
}

// SetCounter sets a counter value
func (m *MemStorage) SetCounter(name string, value int64) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.Counters[name] = value
}

// GetAllGauges gets all gauges
func (m *MemStorage) GetAllGauges() map[string]float64 {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.Gauges
}

// GetAllCounters gets all counters
func (m *MemStorage) GetAllCounters() map[string]int64 {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.Counters
}

// GetAllMetricsSliced gets all metrics
func (m *MemStorage) GetAllMetricsSliced() []metrics {
	m.mut.RLock()
	defer m.mut.RUnlock()
	var metrics []metrics
	for key, value := range m.GetAllCounters() {
		met := NewMetricsCounter(key, "counter", value)
		metrics = append(metrics, *met)
	}
	for key, value := range m.Gauges {
		met := NewMetricsGauge(key, "gauge", value)
		metrics = append(metrics, *met)
	}
	return metrics
}

// GetAllInBytesSliced gets all metrics in bytes
func (m *MemStorage) GetAllInBytesSliced() [][]byte {
	var result [][]byte
	m.mut.RLock()
	defer m.mut.RUnlock()
	for metric, value := range m.GetAllGauges() {
		met := metrics{ID: metric, MType: "gauge", Value: &value}
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

// RestoreMetric restores metric from bytes
func (m *MemStorage) RestoreMetric(data [][]byte) {
	var met metrics
	//m.mut.Lock()
	//defer m.mut.Unlock()
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

// GetGauge gets a gauge value
func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	value, ok := m.Gauges[name]
	return value, ok
}

// GetCounter gets a counter value
func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	value, ok := m.Counters[name]
	return value, ok
}
