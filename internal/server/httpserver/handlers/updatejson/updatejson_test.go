package updatejson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatejson/mocks"
	"github.com/h2p2f/practicum-metrics/internal/server/models"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		method string
		name   string
		metric string
		key    string
		value  float64
		delta  int64
		want   int
	}{
		{
			method: "POST",
			name:   "Test 1",
			metric: "gauge",
			key:    "testKey",
			value:  10.01,
			want:   http.StatusOK,
		},
		{
			method: "POST",
			name:   "Test 2",
			metric: "counter",
			key:    "testKey",
			value:  -10,
			want:   http.StatusBadRequest,
		},
		{
			method: "POST",
			name:   "Test 3",
			metric: "",
			key:    "testKey",
			value:  10,
			want:   http.StatusBadRequest,
		},
		{
			method: "POST",
			name:   "Test 4",
			metric: "gauge",
			key:    "testKey",
			value:  -10.12,
			want:   http.StatusBadRequest,
		},
		{
			method: "GET",
			name:   "Test 5",
			metric: "gauge",
			key:    "testKey",
			value:  10.01,
			want:   http.StatusMethodNotAllowed,
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
				Delta: &tt.delta,
				Value: &tt.value,
				ID:    tt.key,
				MType: tt.metric,
			}

			body, _ := json.Marshal(metric)

			request := httptest.NewRequest(tt.method, "/update/", bytes.NewBuffer(body))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("Handler() got = %v, want %v", response.Code, tt.want)
			}
		})
	}
}

func Example() {
	//создаем тестовый объект
	//
	//create a test object
	gauge := float64(10)
	metric := models.Metric{
		Value: &gauge,
		MType: "gauge",
		ID:    "testKey",
	}
	//маршаллизируем его в json
	//
	//marshal it to json
	body, _ := json.Marshal(metric)
	//создаем тестовую структуру
	//
	//create a test structure
	t := &testing.T{}
	//создаем моковый объект базы данных
	//
	//create a mock database object
	updaterMock := mocks.NewUpdater(t)
	updaterMock.On("SetGauge", metric.ID, gauge).Return(nil)
	//создаем логгер
	//
	//create a logger
	logger := zaptest.NewLogger(t)
	//создаем запрос и стуктуру обработки ответа
	//
	//create a request and response handling structure
	request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	//вызываем обработчик
	//
	//call the handler
	Handler(logger, updaterMock).ServeHTTP(response, request)
	//выводим результаты
	//
	//output the results
	fmt.Println(response.Code)
	fmt.Println(response.Body.String())

	// Output:
	// 200
	// {"value":10,"id":"testKey","type":"gauge"}

}

func BenchmarkHandler(b *testing.B) {
	//создаем тестовый объект
	//
	//create a test object
	gauge := float64(10)
	metric := models.Metric{
		MType: "gauge",
		ID:    "testKey",
		Value: &gauge,
	}
	//маршаллизируем его в json
	//
	//marshal it to json
	body, _ := json.Marshal(metric)
	//создаем тестовую структуру
	//
	//create a test structure
	t := &testing.T{}
	//создаем моковый объект базы данных
	//
	//create a mock database object
	updaterMock := mocks.NewUpdater(t)
	updaterMock.On("SetGauge", metric.ID, gauge).Return(nil)
	//создаем логгер
	//
	//create a logger
	logger := zaptest.NewLogger(t)
	//создаем запрос и стуктуру обработки ответа
	//
	//create a request and response handling structure
	request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(body))
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
		gauge := float64(10)
		metric := models.Metric{
			MType: "gauge",
			ID:    "testKey",
			Value: &gauge,
		}
		//маршаллизируем его в json
		//
		//marshal it to json
		body, _ := json.Marshal(metric)
		//создаем тестовую структуру
		//
		//create a test structure
		t := &testing.T{}
		//создаем моковый объект базы данных
		//
		//create a mock database object
		updaterMock := mocks.NewUpdater(t)
		updaterMock.On("SetGauge", metric.ID, gauge).Return(nil)
		//создаем логгер
		//
		//create a logger
		logger := zaptest.NewLogger(t)
		//создаем запрос и стуктуру обработки ответа
		//
		//create a request and response handling structure
		request := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(body))
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
