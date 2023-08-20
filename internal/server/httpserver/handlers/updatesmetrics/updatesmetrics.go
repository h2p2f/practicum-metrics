// Package updatesmetrics содержит в себе http.Handler, который обновляет метрики в JSON-наборе.
//
// package updatesmetrics contains an http.Handler that updates metrics in JSON batch.
package updatesmetrics

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/server/models"
)

// Updater это интерфейс, который обновляет метрики.
//
// Updater is an interface that updates metrics.
//
//go:generate mockery --name Updater --output ./mocks --filename mocks_updatesmetrics.go
type Updater interface {
	SetGauge(name string, value float64)
	SetCounter(name string, value int64)
}

// Handler возвращает http.HandlerFunc, который обрабатывает POST запросы и обновляет метрики в JSON.
// В негативном случае возвращает внутреннюю ошибку сервера.
// Handler returns a http.HandlerFunc that handles POST requests and updates the batch metric in JSON.
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
				wrappedIFace.SetGauge(metric.ID, *metric.Value)
			case "counter":
				// Check if delta is negative
				if *metric.Delta < 0 {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				// Update the metric
				wrappedIFace.SetCounter(metric.ID, *metric.Delta)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
