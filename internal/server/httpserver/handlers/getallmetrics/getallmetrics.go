// package getallmetrics contains an http.Handler that returns all the metrics.
package getallmetrics

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// Getter is an interface that gets all the metrics.
//
//go:generate mockery --name Getter --output ./mocks --filename mocks_getmetrics.go
type Getter interface {
	GetCounters() map[string]int64
	GetGauges() map[string]float64
}

// Handler returns a http.HandlerFunc that handles GET requests and returns all the metrics.
// It writes the counters and gauges to the response body.
// Otherwise, it returns a method not allowed error.
func Handler(logger *zap.Logger, db Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is not GET
		if r.Method != http.MethodGet {
			logger.Sugar().Infow("method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		wrappedIFace := NewGetterWithZap(db, logger)
		// Get the counters from the database
		counters := wrappedIFace.GetCounters()
		//counters := db.GetCounters()

		// Get the gauges from the database
		gauges := wrappedIFace.GetGauges()
		//gauges := db.GetGauges()

		// Set the Content-Type header to text/html
		w.Header().Add("Content-Type", "text/html")

		// Write the "counters:" text to the response writer
		_, err := w.Write([]byte("counters:<br>"))
		if err != nil {
			logger.Error("could not write response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Write the counters to the response writer
		for k, v := range counters {
			_, err2 := w.Write([]byte(fmt.Sprintf("<p> %s: %d</p>", k, v)))
			if err2 != nil {
				logger.Error("could not write response", zap.Error(err2))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		// Write the "gauges:" text to the response writer
		_, err = w.Write([]byte("gauges:<br>"))
		if err != nil {
			logger.Error("could not write response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Write the gauges to the response writer
		for k, v := range gauges {
			_, err := w.Write([]byte(fmt.Sprintf("<p> %s: %f</p>", k, v)))
			if err != nil {
				logger.Error("could not write response", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
}
