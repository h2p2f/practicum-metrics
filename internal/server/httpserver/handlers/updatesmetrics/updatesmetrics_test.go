package updatesmetrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatesmetrics/mocks"
	"github.com/h2p2f/practicum-metrics/internal/server/models"
)

func TestHandler(t *testing.T) {
	var gauge = 0.001
	var counter int64 = 10
	var wrongCounter int64 = -10

	tests := []struct {
		name    string
		metrics []models.Metric
		want    int
	}{
		{
			name: "Test 1",
			metrics: []models.Metric{
				{
					ID:    "testKey",
					MType: "gauge",
					Value: &gauge,
				},
				{
					ID:    "testKey2",
					MType: "counter",
					Delta: &counter,
				},
			},
			want: http.StatusOK,
		},
		{
			name: "Test 2",
			metrics: []models.Metric{
				{
					ID:    "testKey2",
					MType: "counter",
					Delta: &wrongCounter,
				},
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatersMock := mocks.NewUpdater(t)
			if tt.want == http.StatusOK {
				updatersMock.On("SetGauge", tt.metrics[0].ID, *tt.metrics[0].Value).Return(nil)
				updatersMock.On("SetCounter", tt.metrics[1].ID, *tt.metrics[1].Delta).Return(nil)

			}
			logger := zaptest.NewLogger(t)
			handler := Handler(logger, updatersMock)

			body, _ := json.Marshal(tt.metrics)

			request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(body))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("Handler() got = %v, want %v", response.Code, tt.want)
			}
		})
	}
}

func Example() {
	//объявляем тестовые переменные
	//
	//declare test variables
	var gauge = 0.001
	var counter int64 = 10
	//создаем тестовый объект
	//
	//create test object
	t := &testing.T{}
	//создаем моки
	//

	updatersMock := mocks.NewUpdater(t)
	updatersMock.On("SetGauge", "testKey", gauge).Return(nil)
	updatersMock.On("SetCounter", "testKey", counter).Return(nil)
	//создаем логгер
	//
	//create logger
	logger := zaptest.NewLogger(t)
	//создаем слайс метрик и маршалим его в json
	//
	//create metrics slice and marshal it to json
	metrics := []models.Metric{
		{
			ID:    "testKey",
			MType: "gauge",
			Value: &gauge,
		}, {
			ID:    "testKey",
			MType: "counter",
			Delta: &counter,
		},
	}
	body, _ := json.Marshal(metrics)
	//создаем запрос, записываем в него тело и обявляем объект ответа
	//
	//create request, write body to it and declare response object
	request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	//вызываем обработчик
	//
	//call handler
	Handler(logger, updatersMock).ServeHTTP(response, request)
	//выводим код ответа
	//
	//print response code
	fmt.Println(response.Code)
	// Output:
	// 200
}

func BenchmarkHandler(b *testing.B) {
	//объявляем тестовые переменные
	//
	//declare test variables
	var gauge = 0.001
	var counter int64 = 10
	//создаем тестовый объект
	//
	//create test object
	t := &testing.T{}
	//создаем моки
	//

	updatersMock := mocks.NewUpdater(t)
	updatersMock.On("SetGauge", "testKey", gauge).Return(nil)
	updatersMock.On("SetCounter", "testKey", counter).Return(nil)
	//создаем логгер
	//
	//create logger
	logger := zaptest.NewLogger(t)
	//создаем слайс метрик и маршалим его в json
	//
	//create metrics slice and marshal it to json
	metrics := []models.Metric{
		{
			ID:    "testKey",
			MType: "gauge",
			Value: &gauge,
		}, {
			ID:    "testKey",
			MType: "counter",
			Delta: &counter,
		},
	}
	body, _ := json.Marshal(metrics)
	//создаем запрос, записываем в него тело и обявляем объект ответа
	//
	//create request, write body to it and declare response object
	request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(body))
	response := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	//вызываем обработчик
	//
	//call handler
	for i := 0; i < b.N; i++ {
		Handler(logger, updatersMock).ServeHTTP(response, request)

	}
}

func BenchmarkHandlerParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		//объявляем тестовые переменные
		//
		//declare test variables
		var gauge = 0.001
		var counter int64 = 10
		//создаем тестовый объект
		//
		//create test object
		t := &testing.T{}
		//создаем моки
		//

		updatersMock := mocks.NewUpdater(t)
		updatersMock.On("SetGauge", "testKey", gauge).Return(nil)
		updatersMock.On("SetCounter", "testKey", counter).Return(nil)
		//создаем логгер
		//
		//create logger
		logger := zaptest.NewLogger(t)
		//создаем слайс метрик и маршалим его в json
		//
		//create metrics slice and marshal it to json
		metrics := []models.Metric{
			{
				ID:    "testKey",
				MType: "gauge",
				Value: &gauge,
			}, {
				ID:    "testKey",
				MType: "counter",
				Delta: &counter,
			},
		}
		body, _ := json.Marshal(metrics)
		//создаем запрос, записываем в него тело и обявляем объект ответа
		//
		//create request, write body to it and declare response object
		request := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		b.ReportAllocs()
		b.ResetTimer()
		//вызываем обработчик
		//
		//call handler
		for pb.Next() {
			Handler(logger, updatersMock).ServeHTTP(response, request)
		}
	})
}
