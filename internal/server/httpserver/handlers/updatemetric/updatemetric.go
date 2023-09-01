// Package updatemetric contains an http.Handler that updates the metric and returns http.StatusOK if successful.
package updatemetric

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Updater is an interface that updates the metric.
//
//go:generate mockery --name Updater --output ./mocks --filename mocks_updatemetric.go
type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
}

// Handler returns a http.HandlerFunc that handles POST requests and updates the metric.
// It returns http.StatusOK if successful.
// Otherwise, it returns an internal server error.
// data to update receive in URI
func Handler(log *zap.Logger, db Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the method is POST
		if r.Method != http.MethodPost {
			log.Sugar().Infow("method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		wrappedIFace := NewUpdaterWithZap(db, log)

		// Get the metric type, key and value from the URI
		metric := chi.URLParam(r, "metric")
		key := chi.URLParam(r, "key")
		value := chi.URLParam(r, "value")
		// Check for required fields is valid
		if metric == "" {
			log.Info("bad request")
			http.Error(w, "Bad request", http.StatusNotFound)
			return
		}
		if key == "" || value == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Processing and validation of the received data
		err := updaterMetric(&wrappedIFace, log, metric, key, value)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}

// updaterMetric - function to update the metric
func updaterMetric(updater *UpdaterWithZap, log *zap.Logger, metric, key, value string) error {
	var (
		i   int64
		f   float64
		err error
	)
	switch metric {
	case "gauge":
		// Parse the value to float64
		f, err = strconv.ParseFloat(value, 64)
		if err != nil {
			log.Error("could not parse float", zap.Error(err))
			return err
		}
		// Check if the value is negative
		if f < 0 {
			log.Error("value must be positive")
			return errors.New("value must be positive")
		}
		// Update the metric
		updater.SetGauge(key, f)
	case "counter":
		// Parse the value to int64
		i, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Error("could not parse int", zap.Error(err))
			return err
		}
		// Check if the value is negative
		if i < 0 {
			log.Error("value must be positive")
			return errors.New("value must be positive")
		}
		// Update the metric
		updater.SetCounter(key, i)
	default:
		log.Error("invalid metric type")
		return errors.New("invalid metric type")
	}
	return nil

}
