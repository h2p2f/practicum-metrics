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
		method string
		name   string
		metric string
		key    string
		value  string
		want   int
	}{
		{
			method: "POST",
			name:   "Test 1",
			metric: "gauge",
			key:    "testKey",
			value:  "10.01",
			want:   http.StatusOK,
		},
		{
			method: "POST",
			name:   "Test 2",
			metric: "counter",
			key:    "testKey",
			value:  "-10",
			want:   http.StatusBadRequest,
		},
		{
			method: "POST",
			name:   "Test 3",
			metric: "",
			key:    "testKey",
			value:  "10",
			want:   http.StatusNotFound,
		},
		{
			method: "GET",
			name:   "Test 4",
			metric: "counter",
			key:    "testKey",
			value:  "10",
			want:   http.StatusMethodNotAllowed,
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

			request := httptest.NewRequest(tt.method, link, nil)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("Handler() = %v, want %v", response.Code, tt.want)
			}
		})
	}
}

func Example() {
	//создаем тестовый объект
	//
	//create a test object
	t := &testing.T{}
	//создаем моковый объект базы данных
	//
	//create a mock database object
	updaterMock := mocks.NewUpdater(t)
	updaterMock.On("SetGauge", "testKey", mock.Anything).Return(nil)
	//создаем тестовый объект логгера
	//
	//create a test logger object
	logger := zaptest.NewLogger(t)
	//создаем объект запроса
	//
	//create a request object
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("metric", "gauge")
	rctx.URLParams.Add("key", "testKey")
	rctx.URLParams.Add("value", "10.01")
	request := httptest.NewRequest(http.MethodPost, "/update/gauge/testKey/10.01", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
	//создаем объект ответа
	//
	//create a response object
	response := httptest.NewRecorder()
	//вызываем обработчик
	//
	//call the handler
	Handler(logger, updaterMock).ServeHTTP(response, request)
	//выводим код ответа
	//
	//output response code
	fmt.Println(response.Code)

	// Output:
	// 200
}

func BenchmarkHandler(b *testing.B) {
	//создаем тестовый объект
	//
	//create a test object
	t := &testing.T{}
	//создаем моковый объект базы данных
	//
	//create a mock database object
	updaterMock := mocks.NewUpdater(t)
	updaterMock.On("SetGauge", "testKey", mock.Anything).Return(nil)
	//создаем тестовый объект логгера
	//
	//create a test logger object
	logger := zaptest.NewLogger(t)
	//создаем объект запроса
	//
	//create a request object
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("metric", "gauge")
	rctx.URLParams.Add("key", "testKey")
	rctx.URLParams.Add("value", "10.01")
	request := httptest.NewRequest(http.MethodPost, "/update/gauge/testKey/10.01", nil)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
	//создаем объект ответа
	//
	//create a response object
	response := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	//вызываем обработчик
	//
	//call the handler
	for i := 0; i < b.N; i++ {
		Handler(logger, updaterMock).ServeHTTP(response, request)
	}
}

func BenchmarkHandlerParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		//создаем тестовый объект
		//
		//create a test object
		t := &testing.T{}
		//создаем моковый объект базы данных
		//
		//create a mock database object
		updaterMock := mocks.NewUpdater(t)
		updaterMock.On("SetGauge", "testKey", mock.Anything).Return(nil)
		//создаем тестовый объект логгера
		//
		//create a test logger object
		logger := zaptest.NewLogger(t)
		//создаем объект запроса
		//
		//create a request object
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("metric", "gauge")
		rctx.URLParams.Add("key", "testKey")
		rctx.URLParams.Add("value", "10.01")
		request := httptest.NewRequest(http.MethodPost, "/update/gauge/testKey/10.01", nil)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
		//создаем объект ответа
		//
		//create a response object
		response := httptest.NewRecorder()

		b.ReportAllocs()
		b.ResetTimer()

		for pb.Next() {
			//вызываем обработчик
			//
			//call the handler
			Handler(logger, updaterMock).ServeHTTP(response, request)
		}
	})
}
