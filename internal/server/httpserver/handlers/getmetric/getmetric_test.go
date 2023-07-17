package getmetric

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap/zaptest"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/getmetric/mocks"
)

func TestGetMetric(t *testing.T) {
	tests := []struct {
		name   string
		metric string
		key    string
		want   int
	}{
		{
			name:   "Test 1",
			metric: "gauge",
			key:    "testKey",
			want:   http.StatusOK,
		},
		{
			name:   "Test 2",
			metric: "counter",
			key:    "testKey",
			want:   http.StatusOK,
		},
		{
			name:   "Test 3",
			metric: "",
			key:    "testKey",
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 4",
			metric: "gauge",
			key:    "",
			want:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			getterMock := mocks.NewGetter(t)
			if tt.want == http.StatusOK && tt.metric == "gauge" {
				getterMock.On("GetGauge", tt.key).Return(float64(10), nil)
			}
			if tt.want == http.StatusOK && tt.metric == "counter" {
				getterMock.On("GetCounter", tt.key).Return(int64(1), nil)
			}
			logger := zaptest.NewLogger(t)

			handler := Handler(logger, getterMock)

			link := fmt.Sprintf("/value/%s/%s", tt.metric, tt.key)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metric", tt.metric)
			rctx.URLParams.Add("key", tt.key)

			request := httptest.NewRequest(http.MethodGet, link, nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("Handler() = %v, want %v", response.Code, tt.want)
			}
		})
	}
}
