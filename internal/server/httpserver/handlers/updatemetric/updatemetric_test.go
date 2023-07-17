package updatemetric

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatemetric/mocks"
)

func TestUpdateMetric(t *testing.T) {

	tests := []struct {
		name   string
		metric string
		key    string
		value  string
		want   int
	}{
		{
			name:   "Test 1",
			metric: "gauge",
			key:    "testKey",
			value:  "10.01",
			want:   http.StatusOK,
		},
		{
			name:   "Test 2",
			metric: "counter",
			key:    "testKey",
			value:  "-10",
			want:   http.StatusBadRequest,
		},
		{
			name:   "Test 3",
			metric: "",
			key:    "testKey",
			value:  "10",
			want:   http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			updaterMock := mocks.NewUpdater(t)
			if tt.want == http.StatusOK && tt.metric == "gauge" {
				updaterMock.On("SetGauge", tt.key, mock.Anything).Return(nil)
			}
			if tt.want == http.StatusOK && tt.metric == "counter" {
				updaterMock.On("SetCounter", tt.key, mock.Anything).Return(nil)
			}
			logger := zaptest.NewLogger(t)
			handler := Handler(logger, updaterMock)

			link := fmt.Sprintf("/update/%s/%s/%s", tt.metric, tt.key, tt.value)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metric", tt.metric)
			rctx.URLParams.Add("key", tt.key)
			rctx.URLParams.Add("value", tt.value)

			request := httptest.NewRequest(http.MethodPost, link, nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("Handler() = %v, want %v", response.Code, tt.want)
			}
		})
	}
}
