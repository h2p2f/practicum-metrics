package storage

// MemStorage is a storage in memory
// it is a struct with two maps - gauges and counters

type MemStorage struct {
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
	m.Gauges[name] = append(m.Gauges[name], value)
}

// SetCounter sets a counter value
func (m *MemStorage) SetCounter(name string, value int64) {
	m.Counters[name] = value
}

// GetGauge gets a gauge value
func (m *MemStorage) GetGauge(name string) ([]float64, bool) {
	value, ok := m.Gauges[name]
	return value, ok
}

// GetCounter gets a counter value
func (m *MemStorage) GetCounter(name string) (int64, bool) {
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

func (m *MemStorage) GetAllMetricsSliced() []Metrics {
	var metrics []Metrics
	for key, value := range m.Gauges {
		metrics = append(metrics, Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value[len(value)-1],
		})
	}
	for key, value := range m.Counters {
		metrics = append(metrics, Metrics{
			ID:    key,
			MType: "counter",
			Delta: &value,
		})
	}
	return metrics
}

// RestoreMetrics restores metrics from slice
func (m *MemStorage) RestoreMetrics(metrics []Metrics) {
	for _, metric := range metrics {
		switch metric.MType {
		case "gauge":
			{
				m.SetGauge(metric.ID, *metric.Value)
			}
		case "counter":
			{
				m.SetCounter(metric.ID, *metric.Delta)
			}
		}
	}
}
