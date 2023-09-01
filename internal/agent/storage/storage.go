// Package storage implements a metric storage in memory.
package storage

import (
	"sync"
)

// MetricStorage stores metrics in memory.
type MetricStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mut     sync.RWMutex
}

// NewAgentStorage creates a new metric storage.
func NewAgentStorage() *MetricStorage {
	return &MetricStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}
