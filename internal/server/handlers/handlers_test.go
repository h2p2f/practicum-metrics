package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http/httptest"
	"practicum-metrics/internal/storage"
	"testing"
)

func TestMetricHandler_UpdatePage(t *testing.T) {
	type want struct {
		statusCode   int
		conttentType string
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
				statusCode:   200,
				conttentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Positive test 2",
			metric:      "gauge",
			metricName:  "test",
			metricValue: "1.0000000000001",
			want: want{
				statusCode:   200,
				conttentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 1",
			metric:      "someMetric",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:   400,
				conttentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 2",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1/1",
			want: want{
				statusCode:   404,
				conttentType: "text/plain; charset=utf-8",
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
		statusCode   int
		value        string
		conttentType string
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
				statusCode:   200,
				value:        "1",
				conttentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Positive test 2",
			metric:      "gauge",
			metricName:  "test",
			metricValue: "1.0000000000001",
			want: want{
				statusCode:   200,
				value:        "1.0000000000001",
				conttentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 1",
			metric:      "someMetric",
			metricName:  "test",
			metricValue: "1",
			want: want{
				statusCode:   404,
				conttentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Negative test 2",
			metric:      "counter",
			metricName:  "test",
			metricValue: "1/1",
			want: want{
				statusCode:   404,
				conttentType: "text/plain; charset=utf-8",
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
