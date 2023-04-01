package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"practicum-metrics/internal/storage"
	"testing"
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name        string
		action      string
		metricType  string
		metricName  string
		metricValue string
		want        int
	}{
		{
			name:        "Positive gauge test",
			action:      "update",
			metricType:  "gauge",
			metricName:  "gaugeMetric",
			metricValue: "78",
			want:        200,
		},
		{
			name:        "Positive counter test",
			action:      "update",
			metricType:  "counter",
			metricName:  "counterMetric",
			metricValue: "78",
			want:        200,
		},
		{
			name:        "Negative action test",
			action:      "delete",
			metricType:  "gauge",
			metricName:  "gaugeMetric",
			metricValue: "78",
			want:        400,
		},
		{
			name:        "Negative metric type test",
			action:      "update",
			metricType:  "unknown",
			metricName:  "gaugeMetric",
			metricValue: "78",
			want:        400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/%s/%s/%s/%s", tt.action, tt.metricType, tt.metricName, tt.metricValue)
			request := httptest.NewRequest("POST", url, nil)
			w := httptest.NewRecorder()
			testStorage := storage.NewMemStorage()
			testHandler := NewMetricHandler(testStorage)
			h := http.HandlerFunc(testHandler.MainPage)
			h.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.want, res.StatusCode)
		})

	}
}
