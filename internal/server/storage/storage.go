package storage

import "fmt"

// MemStorage is a storage in memory
// it is a struct with two maps - gauges and counters

type MemStorage struct {
	//mutex is suspended work with files TODO: fix it
	//mut      sync.RWMutex
	Gauges   map[string][]float64
	Counters map[string]int64
	//TODO: deal with scopes
}

// NewMemStorage creates a new MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string][]float64),
		Counters: make(map[string]int64),
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
func (m *MemStorage) GetAllMetricsSliced() []Metrics {
	//m.mut.Lock()
	//defer m.mut.Unlock()
	var metrics []Metrics
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

// RestoreMetrics restores metrics from slice
func (m *MemStorage) RestoreMetrics(metrics []Metrics) {
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
