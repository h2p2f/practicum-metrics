package getmetric

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

//go:generate mockery --name Getter --output ./mocks --filename mocks_getmetric.go
type Getter interface {
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
}

// Handler returns a http.HandlerFunc that handles GET requests and gets the metric.
// It writes the metric value to the response body if the metric is found.
// Otherwise, it returns a not found error.
func Handler(logger *zap.Logger, db Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is not GET
		if r.Method != http.MethodGet {
			logger.Sugar().Infow("method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Get the metric name and key from the URL parameters.
		metric := chi.URLParam(r, "metric")
		key := chi.URLParam(r, "key")
		// Check if the metric name and key are empty.
		if metric == "" || key == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Get the metric value from the database.
		switch metric {
		case "gauge":
			value, err := db.GetGauge(key)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			_, err = w.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))
			if err != nil {
				logger.Error("could not write response", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "counter":
			value, err := db.GetCounter(key)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)

			}
			_, err = w.Write([]byte(strconv.FormatInt(value, 10)))
			if err != nil {
				logger.Error("could not write response", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)

			}
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Set the content type header to text/plain and the status code to 200 OK.
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}

}
