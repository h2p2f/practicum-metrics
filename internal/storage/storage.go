package storage

// Storage is an interface for storage
type Storage interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetGauge(name string) ([]float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGauges() map[string][]float64
	GetAllCounters() map[string]int64
}

// MemStorage is a storage in memory
// it is a struct with two maps - gauges and counters
// gauges is a map of gauge name and slice of gauge values
// counters is a map of counter name and counter value
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
