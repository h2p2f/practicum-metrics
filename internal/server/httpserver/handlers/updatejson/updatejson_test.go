package updatejson

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatejson/mocks"
	"github.com/h2p2f/practicum-metrics/internal/server/models"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		key    string
		value  float64
		delta  int64
		want   int
	}{
		{
			name:   "Test 1",
			metric: "gauge",
			key:    "testKey",
			value:  10.01,
			want:   http.StatusOK,
		},
		{
			name:   "Test 2",
			metric: "counter",
			key:    "testKey",
			value:  -10,
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 3",
			metric: "",
			key:    "testKey",
			value:  10,
			want:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updaterMock := mocks.NewUpdater(t)

			if tt.want == http.StatusOK {
				switch tt.metric {
				case "gauge":
					updaterMock.On("SetGauge", tt.key, tt.value).Return(nil)
				case "counter":
					updaterMock.On("SetCounter", tt.key, tt.value).Return(nil)
					updaterMock.On("GetCounter", tt.key).Return(tt.value, nil)
				}
			}

			logger := zaptest.NewLogger(t)
			handler := Handler(logger, updaterMock)

			metric := models.Metric{
				MType: tt.metric,
				ID:    tt.key,
				Value: &tt.value,
				Delta: &tt.delta,
			}

			body, _ := json.Marshal(metric)

			request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(body))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("Handler() got = %v, want %v", response.Code, tt.want)
			}
		})
	}
}
