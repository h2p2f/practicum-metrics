package updatesmetrics

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/server/models"
)

//go:generate mockery --name Updater --output ./mocks --filename mocks_updatesmetrics.go
type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
}

// Handler returns a http.HandlerFunc that handles POST requests and updates the batch metric in JSON.
// It writes "ok" to the response body if the update is successful.
// Otherwise, it returns an internal server error.
func Handler(log *zap.Logger, db Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the method is POST
		if r.Method != http.MethodPost {
			log.Sugar().Infow("method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Create a new metric struct
		var (
			buf     bytes.Buffer
			metrics []models.Metric
		)
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			log.Error("could not read from body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Decode the JSON request body into the slice of metrics struct
		err = json.Unmarshal(buf.Bytes(), &metrics)
		if err != nil {
			log.Error("could not unmarshal body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Iterate over the slice of metrics
		for _, metric := range metrics {
			switch metric.MType {
			case "gauge":
				// Check if value is negative
				if *metric.Value < 0 {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				// Update the metric
				db.SetGauge(metric.ID, *metric.Value)
			case "counter":
				// Check if delta is negative
				if *metric.Delta < 0 {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				// Update the metric
				db.SetCounter(metric.ID, *metric.Delta)
			}
		}
		// Write "ok" to the response body
		w.WriteHeader(http.StatusOK)
	}
}
