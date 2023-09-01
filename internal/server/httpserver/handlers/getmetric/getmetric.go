// Package getmetric contains an http.Handler that gets the metric and returns its value.
package getmetric

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Getter is an interface that gets the metric.
//
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

		wrappedIFace := NewGetterWithZap(db, logger)
		// Get the metric name and key from the URL parameters.
		metric := chi.URLParam(r, "metric")
		key := chi.URLParam(r, "key")
		// Check if the metric name and key are empty.
		if metric == "" || key == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Get the metric value from the database.
		value, err := getterMetric(&wrappedIFace, logger, metric, key)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		_, err = w.Write([]byte(value))
		if err != nil {
			logger.Error("could not write response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Set the content type header to text/plain and the status code to 200 OK.
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}

}

// getterMetric - function to get the metric
func getterMetric(getter *GetterWithZap, logger *zap.Logger, metric, key string) (string, error) {
	var (
		i   int64
		f   float64
		err error
	)
	switch metric {
	case "gauge":
		// Get the gauge value from the database.
		f, err = getter.GetGauge(key)
		if err != nil {
			logger.Error("could not get gauge", zap.Error(err))
			return "", err
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case "counter":
		// Get the counter value from the database.
		i, err = getter.GetCounter(key)
		if err != nil {
			logger.Error("could not get counter", zap.Error(err))
			return "", err
		}
		return strconv.FormatInt(i, 10), nil
	default:
		return "", err
	}
}
