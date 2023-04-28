package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"math/rand"
	"net/http/httptest"
	"strings"
	"testing"
)

// test for Handlers
func TestMetricHandler_UpdatePage(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name        string
		metric      string
		metricName  string
		metricValue string
		want        want
	}{
		{
			name:        "Positive test 1",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:  200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Positive test 2",
			metric:      "gauge",
			metricName:  "test",
			metricValue: "1.0000000000001",
			want: want{
				statusCode:  200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 1",
			metric:      "someMetric",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:  501,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 2",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1/1",
			want: want{
				statusCode:  404,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/update/"+tt.metric+"/"+tt.metricName+"/"+tt.metricValue, nil)
			r := chi.NewRouter()
			testStorage := storage.NewMemStorage()
			handler := NewMetricHandler(testStorage)
			r.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.want.statusCode {
				t.Errorf("MetricHandler.UpdatePage() = %v, want %v", w.Code, tt.want.statusCode)
				//TODO: add content type check
			}
		})
	}
}

func TestMetricHandler_GetMetricValue(t *testing.T) {
	type want struct {
		statusCode  int
		value       string
		contentType string
	}
	tests := []struct {
		name        string
		metric      string
		metricName  string
		metricValue string
		want        want
	}{
		{
			name:        "Positive test 1",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:  200,
				value:       "1",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Positive test 2",
			metric:      "gauge",
			metricName:  "test",
			metricValue: "1.0000000000001",
			want: want{
				statusCode:  200,
				value:       "1.0000000000001",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 1",
			metric:      "someMetric",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:  404,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 2",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1/1",
			want: want{
				statusCode:  404,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqPost := httptest.NewRequest("POST", "/update/"+tt.metric+"/"+tt.metricName+"/"+tt.metricValue, nil)
			reqGet := httptest.NewRequest("GET", "/value/"+tt.metric+"/"+tt.metricName, nil)
			r := chi.NewRouter()
			testStorage := storage.NewMemStorage()
			handler := NewMetricHandler(testStorage)
			r.Get("/value/{metric}/{key}", handler.GetMetricValue)
			r.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, reqPost)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, reqGet)
			if w.Code != tt.want.statusCode {
				t.Errorf("MetricHandler.GetMetricValue() = %v, want %v", w.Code, tt.want.statusCode)
			}
		})
	}
}

func TestMetricHandler_GetMetricCounterSum(t *testing.T) {
	type want struct {
		statusCode  int
		value       string
		contentType string
	}
	tests := []struct {
		name        string
		metric      string
		metricName  string
		metricValue string
		want        want
	}{
		{
			name:        "Positive test 1",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:  200,
				value:       "1",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Positive test 2",
			metric:      "counter",
			metricName:  "test",
			metricValue: "2",
			want: want{
				statusCode:  200,
				value:       "3",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Positive test 3",
			metric:      "counter",
			metricName:  "test",
			metricValue: "3",
			want: want{
				statusCode:  200,
				value:       "6",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	testStorage := storage.NewMemStorage()
	handler := NewMetricHandler(testStorage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqPost := httptest.NewRequest("POST", "/update/"+tt.metric+"/"+tt.metricName+"/"+tt.metricValue, nil)
			reqGet := httptest.NewRequest("GET", "/sum/"+tt.metric+"/"+tt.metricName, nil)
			r := chi.NewRouter()

			r.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
			r.Get("/sum/{metric}/{key}", handler.GetMetricValue)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, reqPost)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, reqGet)
			if w.Body.String() != tt.want.value {
				t.Errorf("MetricHandler.GetMetricCounterSum() = %v, want %v", w.Body.String(), tt.want.value)
			}
		})
	}
}

func TestMetricHandler_MainPage(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		header      string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Positive test 1",
			want: want{
				statusCode:  200,
				contentType: "text/html",
				header:      "<h1>Metrics</h1>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			r := chi.NewRouter()
			testStorage := storage.NewMemStorage()
			handler := NewMetricHandler(testStorage)
			r.Get("/", handler.MainPage)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tt.want.statusCode {
				t.Errorf("MetricHandler.MainPage() = %v, want %v", w.Code, tt.want.statusCode)
			}
			if w.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("MetricHandler.MainPage() = %v, want %v", w.Header().Get("Content-Type"), tt.want.contentType)
			}
			if !strings.Contains(w.Body.String(), tt.want.header) {
				t.Errorf("MetricHandler.MainPage() = %v, want %v", w.Body.String(), tt.want.header)
			}
		})
	}
}

func TestMetricHandler_ValueJSON(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name        string
		metric      string
		metricName  string
		metricDelta *int64
		metricValue *float64
		want        want
	}{
		{
			name:       "Positive test 1",
			metric:     "counter",
			metricName: "test",

			want: want{
				statusCode:  200,
				contentType: "application/json",
			},
		},
		{
			name:       "Positive test 2",
			metric:     "gauge",
			metricName: "test-gauge",

			want: want{
				statusCode:  200,
				contentType: "application/json",
			},
		},
		{
			name:       "Positive test 3",
			metric:     "gauge",
			metricName: "test-another-gauge",

			want: want{
				statusCode:  200,
				contentType: "application/json",
			},
		},
	}
	testStorage := storage.NewMemStorage()
	handler := NewMetricHandler(testStorage)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := rand.Int63()
			value := rand.Float64()
			metricUpd := Metrics{}
			metricsVal := Metrics{
				ID:    tt.metricName,
				MType: tt.metric,
			}
			if tt.metric == "counter" {
				metricUpd = Metrics{
					ID:    tt.metricName,
					MType: tt.metric,
					Delta: &delta,
				}
			} else {
				metricUpd = Metrics{
					ID:    tt.metricName,
					MType: tt.metric,
					Value: &value,
				}
			}
			dataUpd, err := json.Marshal(metricUpd)
			if err != nil {
				t.Errorf("MetricHandler.ValueJSON() = %v", err)
			}
			dataVal, err := json.Marshal(metricsVal)
			if err != nil {
				t.Errorf("MetricHandler.ValueJSON() = %v", err)
			}
			reqPostUpdate := httptest.NewRequest("POST", "/update/", bytes.NewBuffer(dataUpd))
			reqPostUpdate.Header.Set("Content-Type", "application/json")
			reqPostValue := httptest.NewRequest("POST", "/value/", bytes.NewBuffer(dataVal))
			reqPostValue.Header.Set("Content-Type", "application/json")
			r := chi.NewRouter()
			r.Post("/update/", handler.UpdateJSON)
			r.Post("/value/", handler.ValueJSON)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, reqPostUpdate)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, reqPostValue)
			if w.Code != tt.want.statusCode {
				t.Errorf("MetricHandler.ValueJSON() = %v, want %v", w.Code, tt.want.statusCode)
			}
			if w.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("MetricHandler.ValueJSON() = %v, want %v", w.Header().Get("Content-Type"), tt.want.contentType)
			}
			//TODO: check and match response body
		})
	}
}
