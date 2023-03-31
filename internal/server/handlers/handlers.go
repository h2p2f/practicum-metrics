package handlers

import (
	"net/http"
	"practicum-metrics/internal/storage"
	"strconv"
	"strings"
)

// MetricHandler is a handler for metrics
type MetricHandler struct {
	Storage storage.Storage
}

// NewMetricHandler creates a new MetricHandler
func NewMetricHandler(s storage.Storage) *MetricHandler {
	return &MetricHandler{Storage: s}
}

// MainPage is a handler for metrics (POST requests)
// now it works only with requests like this: POST http://localhost:8080/update/gauge/gaugeMetric/78
func (m *MetricHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	//check method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//check path, separate it and get values
	path := r.URL.Path
	metrics := strings.Split(path, "/")
	//check if path is correct
	if len(metrics) != 5 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//get values from path
	action, metric, key, value := metrics[1], metrics[2], metrics[3], metrics[4]
	//check if action is correct
	if action == "update" {
		//check if metric is correct
		switch strings.ToLower(metric) {
		//if metric is correct, set value
		case "counter":
			{
				if n, err := strconv.ParseInt(value, 10, 64); err == nil {
					m.Storage.SetCounter(key, n)
				}
			}
		case "gauge":
			{
				if n, err := strconv.ParseFloat(value, 64); err == nil {
					m.Storage.SetGauge(key, n)
				}
			}
		}
		//this code for debug
		//for k, v := range m.Storage.GetAllCounters() {
		//	fmt.Println("key", k, "value", v)
		//}
		//for k, v := range m.Storage.GetAllGauges() {
		//	fmt.Println("key", k, "value", v)
		//}
	}
}
