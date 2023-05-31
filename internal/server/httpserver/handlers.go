package httpserver

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// Storager is an interface for model
type Storager interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}

type DataBaseHandler interface {
	PingContext(ctx context.Context) error
}

// MetricHandler is a handler for metrics
type MetricHandler struct {
	Storage   Storager
	DBHandler DataBaseHandler
	Key       string
}

// metrics is a struct for metrics with json tags
type metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewMetricHandler creates a new MetricHandler
func NewMetricHandler(s Storager, baser DataBaseHandler, k string) *MetricHandler {
	return &MetricHandler{Storage: s, DBHandler: baser, Key: k}
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
	if metric == "" || key == "" || value == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch strings.ToLower(metric) {
	case "counter":
		{
			n, err := strconv.ParseInt(value, 0, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			if n < 0 {
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
			if n < 0 {
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
	if metric == "" || key == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
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
			_, er := w.Write([]byte(strconv.FormatFloat(n, 'f', -1, 64)))
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
	_, err := w.Write([]byte("<h1>metrics</h1>"))
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
		_, err := w.Write([]byte(fmt.Sprintf("<p> %s: %f</p>", k, v)))
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
	var MetricFromRequest metrics

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	checkSum := r.Header.Get("HashSHA256")
	if checkSum != "" && m.Key != "" {
		controlCheckSum := fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))

		if controlCheckSum != checkSum {
			fmt.Println("wrong checksum", controlCheckSum, checkSum)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}

	if err = json.Unmarshal(buf.Bytes(), &MetricFromRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if MetricFromRequest.ID == "" || MetricFromRequest.MType == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	switch strings.ToLower(MetricFromRequest.MType) {
	case "counter":
		{
			if *MetricFromRequest.Delta < 0 {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			currentValue, _ := m.Storage.GetCounter(MetricFromRequest.ID)
			m.Storage.SetCounter(MetricFromRequest.ID, *MetricFromRequest.Delta+currentValue)
			*MetricFromRequest.Delta = *MetricFromRequest.Delta + currentValue
		}
	case "gauge":
		{
			if *MetricFromRequest.Value < 0 {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			m.Storage.SetGauge(MetricFromRequest.ID, *MetricFromRequest.Value)
		}
	default:
		{
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
	//this code is for an incorrect test of 11 increments - to get metrics,
	//the test immediately accesses the database without waiting for a White to it
	//time.Sleep(1 * time.Second)
	//prepare response
	response, _ := json.Marshal(MetricFromRequest)

	if m.Key != "" {
		hash, err := GetHash(m.Key, response)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		hashHeader := fmt.Sprintf("%x", hash)
		w.Header().Set("HashSHA256", hashHeader)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	//this code below does not work with gzip middleware
	//so i hard nailed the header in the middleware code
	//TODO: fix it
	//w.WriteHeader(http.StatusOK)
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
	var MetricFromRequest metrics

	checkSum := r.Header.Get("HashSHA256")
	if checkSum != "" && m.Key != "" {
		ok, err := checkDataHash(checkSum, m.Key, buf.Bytes())
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &MetricFromRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if MetricFromRequest.ID == "" || MetricFromRequest.MType == "" {
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
			MetricFromRequest.Value = &value
		}
	}
	//still prepare response
	response, err := json.Marshal(MetricFromRequest)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if m.Key != "" {
		hash, err := GetHash(m.Key, response)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		hashHeader := fmt.Sprintf("%x", hash)
		w.Header().Set("HashSHA256", hashHeader)
	}
	w.Header().Add("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}

// DBPing is a handler for check connect to DB (GET requests to /ping)
func (m *MetricHandler) DBPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := m.DBHandler.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return
	}
}

// UpdatesBatch is a handler for metrics (POST requests to /updates)
func (m *MetricHandler) UpdatesBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var buf bytes.Buffer
	var MetricsFromRequest []metrics
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	checkSum := r.Header.Get("HashSHA256")
	if checkSum != "" && m.Key != "" {
		ok, err := checkDataHash(checkSum, m.Key, buf.Bytes())
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}

	err = json.Unmarshal(buf.Bytes(), &MetricsFromRequest)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	for _, metric := range MetricsFromRequest {
		switch strings.ToLower(metric.MType) {
		case "counter":
			{
				if *metric.Delta < 0 {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				currentValue, _ := m.Storage.GetCounter(metric.ID)
				m.Storage.SetCounter(metric.ID, *metric.Delta+currentValue)
				*metric.Delta = *metric.Delta + currentValue
			}
		case "gauge":
			{
				if *metric.Value < 0 {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				m.Storage.SetGauge(metric.ID, *metric.Value)
			}
		default:
			{
				http.Error(w, "Not implemented", http.StatusNotImplemented)
			}
		}
	}
	answer := []byte("OK")
	if m.Key != "" {
		hash, err := GetHash(m.Key, answer)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		hashHeader := fmt.Sprintf("%x", hash)
		w.Header().Set("HashSHA256", hashHeader)
	}
	w.WriteHeader(http.StatusOK)
}
