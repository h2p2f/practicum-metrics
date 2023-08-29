// package updatejson contains an http.Handler that updates the metric in JSON and returns "ok" if successful.
// It is also used to get the metric from the store in JSON format.
package updatejson

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/server/models"
)

// Updater is an interface that updates the metric.
// It is also used to get the metric from the store.
//
//go:generate mockery --name Updater --output ./mocks --filename mocks_updatejson.go
type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
}

// Handler returns a http.HandlerFunc that handles POST requests and updates the metric in JSON.
// It writes the updated value to the response body if the update is successful.
// Otherwise, it returns an internal server error.
func Handler(log *zap.Logger, db Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the method is POST
		if r.Method != http.MethodPost {
			log.Sugar().Infow("method not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		wrappedIFace := NewUpdaterWithZap(db, log)
		// Create a new metric struct
		var metric models.Metric
		// Decode the JSON request body into the metric struct
		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			log.Error("could not unmarshal json", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check for required fields in the metric
		if metric.ID == "" || metric.MType == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Check if delta is negative
		if metric.Delta != nil && *metric.Delta < 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Check if value is negative
		if metric.Value != nil && *metric.Value < 0 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		// Get the path from the URL
		path := r.URL.Path
		// if path is /update/, update the metric and return updated value
		switch path {
		case "/update/":
			{
				switch metric.MType {
				case "gauge":
					wrappedIFace.SetGauge(metric.ID, *metric.Value)
				case "counter":
					{
						wrappedIFace.SetCounter(metric.ID, *metric.Delta)
						*metric.Delta, _ = wrappedIFace.GetCounter(metric.ID)
					}
				default:
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
			}
		// if path is /value/, get the metric and return the value
		case "/value/":
			{
				switch metric.MType {
				case "gauge":
					{
						n, _ := wrappedIFace.GetGauge(metric.ID)
						metric.Value = &n
					}
				case "counter":
					{
						n, _ := wrappedIFace.GetCounter(metric.ID)
						metric.Delta = &n
					}
				}
			}
		}
		// Marshal the metric struct into JSON
		resp, err := json.Marshal(metric)
		if err != nil {
			log.Error("could not marshal json", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set response headers and write the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(resp); err != nil {
			log.Error("could not write response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
