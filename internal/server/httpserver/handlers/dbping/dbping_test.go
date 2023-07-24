package dbping

import (
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/dbping/mocks"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	db := mocks.NewPinger(t)

	//прописываем ожидаемый результат
	//
	//specify the expected result
	db.On("Ping").Return(nil)

	//создаем объект запроса
	//
	//create a request object
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	//создаем объект записи ответа
	//
	//create a response record object
	rr := httptest.NewRecorder()

	//вызываем обработчик
	//
	//call the handler
	Handler(logger, db).ServeHTTP(rr, req)

	//выводим результат
	//
	//display the result
	fmt.Println(rr.Body.String())

	// Output:
	// pong
}

func BenchmarkHandler(b *testing.B) {
	//создаем тестовый объект
	//
	//create a test object
	t := &testing.B{}

	//создаем тестовый объект логгера
	//
	//create a test logger object
	logger := zaptest.NewLogger(t)

	//создаем моковый объект базы данных
	//
	//create a mock database object
	db := mocks.NewPinger(t)

	//прописываем ожидаемый результат
	//
	//specify the expected result
	db.On("Ping").Return(nil)

	//создаем объект запроса
	//
	//create a request object
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

	//создаем объект записи ответа
	//
	//create a response record object
	rr := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	//вызываем обработчик
	//
	//call the handler
	for i := 0; i < b.N; i++ {
		Handler(logger, db).ServeHTTP(rr, req)
	}

}

func BenchmarkHandlerParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		//создаем тестовый объект
		//
		//create a test object
		t := &testing.B{}

		//создаем тестовый объект логгера
		//
		//create a test logger object
		logger := zaptest.NewLogger(t)

		//создаем моковый объект базы данных
		//
		//create a mock database object
		db := mocks.NewPinger(t)

		//прописываем ожидаемый результат
		//
		//specify the expected result
		db.On("Ping").Return(nil)

		//создаем объект запроса
		//
		//create a request object
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)

		//создаем объект записи ответа
		//
		//create a response record object
		rr := httptest.NewRecorder()

		b.ReportAllocs()
		b.ResetTimer()

		//вызываем обработчик
		//
		//call the handler
		for pb.Next() {
			Handler(logger, db).ServeHTTP(rr, req)
		}
	})

}
