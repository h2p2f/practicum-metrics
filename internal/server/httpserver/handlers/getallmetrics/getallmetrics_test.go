package getallmetrics

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/getallmetrics/mocks"
)

func TestGetAllMetrics(t *testing.T) {

	tests := []struct {
		name   string
		method string
		want   int
	}{
		{
			name:   "Test 1",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "Test 2",
			method: http.MethodPost,
			want:   http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			getterMock := mocks.NewGetter(t)
			if tt.method == http.MethodGet {
				getterMock.On("GetCounters").Return(map[string]int64{"testKey": 1})
				getterMock.On("GetGauges").Return(map[string]float64{"test1": 10})
			}
			logger := zaptest.NewLogger(t)
			handler := Handler(logger, getterMock)

			request := httptest.NewRequest(tt.method, "/", nil)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.want {
				t.Errorf("GetAllMetrics() = %v, want %v", response.Code, tt.want)
			}

		})
	}
}

func Example() {
	//создаем тестовый объект
	//
	//create a test object
	t := &testing.T{}

	//создаем тестовый объект логгера
	//
	//create a test logger object
	logger := zaptest.NewLogger(t)

	//создаем моковый объект базы данных
	//
	//create a mock database object
	getterMock := mocks.NewGetter(t)

	//прописываем ожидаемый результат
	//
	//specify the expected result
	getterMock.On("GetCounters").Return(map[string]int64{"testKey": 1})
	getterMock.On("GetGauges").Return(map[string]float64{"test1": 10})

	//создаем объект запроса
	//
	//create a request object
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	//создаем объект записи ответа
	//
	//create a response record object
	rr := httptest.NewRecorder()

	//вызываем обработчик
	//
	//call the handler
	Handler(logger, getterMock).ServeHTTP(rr, req)

	//выводим результат
	//
	//display the result
	fmt.Println(rr.Body.String())

	// Output:
	//counters:<br><p> testKey: 1</p>gauges:<br><p> test1: 10.000000</p>
}
