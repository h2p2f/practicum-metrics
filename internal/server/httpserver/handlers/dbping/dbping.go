// Package dbping содержит в себе http.Handler, который пингует базу данных и возвращает "pong" в случае успеха.
//
// package dbping contains an http.Handler that pings the database and returns "pong" if successful.
package dbping

import (
	"net/http"

	"go.uber.org/zap"
)

// Pinger это интерфейс, который пингует базу данных.
//
// Pinger is an interface that pings the database.
//
//go:generate mockery --name Pinger --output ./mocks --filename mocks_ping.go
type Pinger interface {
	Ping() error
}

// Handler возвращает http.HandlerFunc, который обрабатывает GET запросы и пингует базу данных.
// Он записывает "pong" в тело ответа, если пинг успешен.
// В противном случае возвращает внутреннюю ошибку сервера.
//
// Handler returns a http.HandlerFunc that handles GET requests and pings the database.
// It writes "pong" to the response body if the ping is successful.
// Otherwise, it returns an internal server error.
func Handler(logger *zap.Logger, db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is GET.
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Ping the database.
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Info("Ping successful")
		// Write "pong" to the response body.
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("pong"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
