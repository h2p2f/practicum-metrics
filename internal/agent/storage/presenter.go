package storage

import (
	"encoding/json"
	"fmt"

	"github.com/h2p2f/practicum-metrics/internal/agent/models"
)

// JSONMetrics is a method of the MetricStorage structure that generates a slice of
// JSON objects to send metrics to the server.
// Required for backward compatibility. Currently not used.
func (m *MetricStorage) JSONMetrics() [][]byte {

	var res [][]byte
	var model models.Metric
	m.mut.RLock()

	for metric, value := range m.gauge {
		value := value
		model = models.Metric{
			ID:    metric,
			Value: &value,
			MType: "gauge",
		}
		mj, err := json.Marshal(model)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, mj)
	}
	for metric, value := range m.counter {
		value := value
		model = models.Metric{
			ID:    metric,
			Delta: &value,
			MType: "counter",
		}
		mj, err := json.Marshal(model)
		if err != nil {
			fmt.Println(err)
		}
		res = append(res, mj)
	}
	m.mut.RUnlock()
	m.mut.Lock()
	defer m.mut.Unlock()
	m.counter["PollCount"] = 0
	return res
}

// BatchJSONMetrics is a method of the MetricStorage structure that generates
// a batch JSON object to send metrics to the server.
func (m *MetricStorage) BatchJSONMetrics() []byte {
	var res []byte
	var modelSlice []models.Metric
	m.mut.RLock()

	for metric, value := range m.gauge {
		value := value
		mt := models.Metric{
			ID:    metric,
			Value: &value,
			MType: "gauge",
		}
		modelSlice = append(modelSlice, mt)
	}
	for metric, value := range m.counter {
		value := value
		mt := models.Metric{
			ID:    metric,
			Delta: &value,
			MType: "counter",
		}
		modelSlice = append(modelSlice, mt)
	}
	res, err := json.Marshal(modelSlice)
	if err != nil {
		fmt.Println(err)
	}
	m.mut.RUnlock()
	m.mut.Lock()
	defer m.mut.Unlock()
	m.counter["PollCount"] = 0
	return res
}

func (m *MetricStorage) GetAllGauge() map[string]float64 {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.gauge
}

func (m *MetricStorage) GetAllCounter() map[string]int64 {
	m.mut.RLock()
	defer m.mut.RUnlock()
	return m.counter
}
