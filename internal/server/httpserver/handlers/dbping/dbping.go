package dbping

import (
	"net/http"

	"go.uber.org/zap"
)

//go:generate mockery --name Pinger --output ./mocks --filename mocks_ping.go
type Pinger interface {
	Ping() error
}

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
		// Write "pong" to the response body.
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("pong"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
