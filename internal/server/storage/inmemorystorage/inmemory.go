package inmemorystorage

import (
	"encoding/json"
	"log"
	"sync"

	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/server/models"
	"github.com/h2p2f/practicum-metrics/internal/server/servererrors"
)

type memStorage struct {
	mut      sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
	logger   *zap.Logger
}

// NewMemStorage creates a new instance of memStorage.
func NewMemStorage(log *zap.Logger) *memStorage {
	return &memStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
		logger:   log,
	}
}

// SetGauge sets the value of a gauge in the memStorage.
func (m *memStorage) SetGauge(name string, value float64) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.gauges[name] = value
}

// SetCounter sets the counter value for the given name.
func (m *memStorage) SetCounter(name string, value int64) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.counters[name] = m.counters[name] + value
}

// GetGauge returns the value of the gauge with the given name.
// If the gauge is not found, it returns 0 and an error.
func (m *memStorage) GetGauge(name string) (float64, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	value, ok := m.gauges[name]
	if !ok {
		return 0, servererrors.ErrNotFound
	}
	return value, nil
}

// GetCounter returns the counter value for the given name.
// If the counter does not exist, it returns 0 and an error.
func (m *memStorage) GetCounter(name string) (int64, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()
	value, ok := m.counters[name]
	if !ok {
		return 0, servererrors.ErrNotFound
	}
	return value, nil
}

func (m *memStorage) GetCounters() map[string]int64 {
	return m.counters
}

func (m *memStorage) GetGauges() map[string]float64 {
	return m.gauges
}

func (m *memStorage) GetAllSerialized() [][]byte {
	var result [][]byte
	var met models.Metric
	m.mut.RLock()
	defer m.mut.RUnlock()
	for metric, value := range m.gauges {
		met.ID = metric
		met.MType = "gauge"
		met.Value = &value
		out, err := json.Marshal(met)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}
	for metric, value := range m.counters {
		met.ID = metric
		met.MType = "counter"
		met.Delta = &value
		out, err := json.Marshal(met)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, out)
	}
	return result
}

func (m *memStorage) RestoreFromSerialized(data [][]byte) error {

	var met models.Metric
	for _, value := range data {
		err := json.Unmarshal(value, &met)
		if err != nil {
			return err
		}
		switch met.MType {
		case "counter":
			m.SetCounter(met.ID, *met.Delta)
		case "gauge":
			m.SetGauge(met.ID, *met.Value)
		}
	}
	return nil
}

func (m *memStorage) Ping() error {
	return servererrors.ErrNotImplemented
}
