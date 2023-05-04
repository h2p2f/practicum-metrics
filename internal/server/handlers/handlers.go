package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Storager is an interface for storage
type Storager interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetGauge(name string) ([]float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGauges() map[string][]float64
	GetAllCounters() map[string]int64
	//GetAllMetricsSliced() []Metrics
}

type DataBaser interface {
	PingContext(ctx context.Context) error
}

// MetricHandler is a handler for metrics
type MetricHandler struct {
	Storage Storager
	DB      DataBaser
}

// Metrics is a struct for metrics with json tags
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewMetricHandler creates a new MetricHandler
func NewMetricHandler(s Storager, baser DataBaser) *MetricHandler {
	return &MetricHandler{Storage: s, DB: baser}
}

// UpdatePage is a handler for metrics (POST requests)
// now it works only with requests like this:
// POST http://localhost:8080/update/gauge/gaugeMetric/78
// port is variable, set it in main.go
func (m *MetricHandler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	//check method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//get metric, key and value from request with chi's URLParam
	metric := chi.URLParam(r, "metric")
	key := chi.URLParam(r, "key")
	value := chi.URLParam(r, "value")
	//prepare metric and set value
	switch strings.ToLower(metric) {
	case "counter":
		{
			n, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			currentValue, _ := m.Storage.GetCounter(key)
			m.Storage.SetCounter(key, n+currentValue)
			w.WriteHeader(http.StatusOK)
		}
	case "gauge":
		{
			n, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			m.Storage.SetGauge(key, n)
			w.WriteHeader(http.StatusOK)
		}
		//if metric is not counter or gauge, return bad request
	default:
		{
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
}

// GetMetricValue is a handler for metrics (GET requests)
// now it works only with requests like this:
// GET http://localhost:8080/value/gauge/gaugeMetric
// port is variable, set it in main.go
func (m *MetricHandler) GetMetricValue(w http.ResponseWriter, r *http.Request) {
	//check method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//get metric and key from request with chi's URLParam
	metric := chi.URLParam(r, "metric")
	key := chi.URLParam(r, "key")
	//get metric and return value
	switch strings.ToLower(metric) {
	case "counter":
		{
			n, err := m.Storage.GetCounter(key)
			//if there is no such key, return not found
			if !err {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			//if there is such key, return value
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, er := w.Write([]byte(strconv.FormatInt(n, 10)))
			if er != nil {
				fmt.Println(err)
			}
		}
	case "gauge":
		{
			n, err := m.Storage.GetGauge(key)
			//if there is no such key, return not found
			if !err {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			//if there is such key, return value
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, er := w.Write([]byte(strconv.FormatFloat(n[len(n)-1], 'f', -1, 64)))
			if er != nil {
				fmt.Println(err)
			}
			//if mentors will say to return all values, uncomment this
			//w.Write([]byte(fmt.Sprintf("%d", n)))
		}
		//if metric is not counter or gauge, return not found
	default:
		{
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
	}
}

// MainPage is a handler for metrics (GET requests to main page)
func (m *MetricHandler) MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//prepare response
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("<h1>Metrics</h1>"))
	if err != nil {
		fmt.Println(err)
	}
	//get all counters and gauges and write them to response
	for k, v := range m.Storage.GetAllCounters() {
		_, err := w.Write([]byte(fmt.Sprintf("<p> %s: %d</p>", k, v)))
		if err != nil {
			fmt.Println(err)
		}
	}
	for k, v := range m.Storage.GetAllGauges() {
		_, err := w.Write([]byte(fmt.Sprintf("<p> %s: %f</p>", k, v[len(v)-1])))
		if err != nil {
			fmt.Println(err)
		}
		//if mentors will say to return all values, uncomment this
		//w.Write([]byte(fmt.Sprintf("<p> %s: %v</p>", k, v)))
	}
}

// UpdateJSON is a handler for metrics (POST requests)
func (m *MetricHandler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	//check method and content-type
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//read request body and unmarshal it to MetricFromRequest
	var buf bytes.Buffer
	var MetricFromRequest Metrics

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &MetricFromRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch strings.ToLower(MetricFromRequest.MType) {
	case "counter":
		{
			currentValue, _ := m.Storage.GetCounter(MetricFromRequest.ID)
			m.Storage.SetCounter(MetricFromRequest.ID, *MetricFromRequest.Delta+currentValue)
			*MetricFromRequest.Delta = *MetricFromRequest.Delta + currentValue
		}
	case "gauge":
		{
			m.Storage.SetGauge(MetricFromRequest.ID, *MetricFromRequest.Value)
		}
	default:
		{
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
	//prepare response
	response, _ := json.Marshal(MetricFromRequest)
	w.Header().Add("Content-Type", "application/json")
	//this code below does not work with gzip middleware
	//so i hard nailed the header in the middleware code
	//TODO: fix it
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		fmt.Println(err)
	}

}

// ValueJSON is a handler for metrics (POST requests to /value)
func (m *MetricHandler) ValueJSON(w http.ResponseWriter, r *http.Request) {
	//check method and content-type
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//read request body and unmarshal it to MetricFromRequest
	var buf bytes.Buffer
	var MetricFromRequest Metrics
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &MetricFromRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	//prepare response
	switch strings.ToLower(MetricFromRequest.MType) {
	case "counter":
		{
			n, ok := m.Storage.GetCounter(MetricFromRequest.ID)
			if !ok {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			MetricFromRequest.Delta = &n
		}
	case "gauge":
		{
			value, ok := m.Storage.GetGauge(MetricFromRequest.ID)
			if !ok {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			MetricFromRequest.Value = &value[len(value)-1]
		}
	}
	//still prepare response
	response, err := json.Marshal(MetricFromRequest)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *MetricHandler) DbPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := m.DB.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return
	}
}
