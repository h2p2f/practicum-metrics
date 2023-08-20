// Package updatemetric содержит в себе http.Handler, который обновляет метрику и возвращает http.StatusOK в случае успеха.
//
// package updatemetric contains an http.Handler that updates the metric and returns http.StatusOK if successful.
package updatemetric

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Updater это интерфейс, который обновляет метрику.
//
// Updater is an interface that updates the metric.
//
//go:generate mockery --name Updater --output ./mocks --filename mocks_updatemetric.go
type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
}

// Handler возвращает http.HandlerFunc, который обрабатывает POST запросы и обновляет метрику.
// Он возвращает http.StatusOK в случае успеха.
// В противном случае возвращает внутреннюю ошибку сервера.
// данные для обновления получает в URI
//
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
			log.Sugar().Infow("bad request")

			http.Error(w, "Bad request", http.StatusNotFound)
			return
		}
		if key == "" || value == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Processing and validation of the received data
		switch metric {
		case "gauge":
			// Parse the value to float64
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Error("could not parse float", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// Check if the value is negative
			if f < 0 {
				log.Error("value must be positive")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// Update the metric
			wrappedIFace.SetGauge(key, f)
		case "counter":
			// Parse the value to int64
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Error("could not parse int", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// Check if the value is negative
			if i < 0 {
				log.Error("value must be positive")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// Update the metric
			wrappedIFace.SetCounter(key, i)
		// If the metric type is unknown, return a bad request
		default:
			log.Error("invalid metric type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
