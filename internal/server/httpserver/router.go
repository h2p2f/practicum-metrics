package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/dbping"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/getallmetrics"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/getmetric"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatejson"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatemetric"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/handlers/updatesmetrics"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/middlewares/compressormiddleware"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/middlewares/hashmiddleware"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/middlewares/loggermiddleware"
)

func MetricRouter(logger *zap.Logger, m DataBaser, key string) *chi.Mux {
	db := NewDataBase(m)
	r := chi.NewRouter()

	r.Use(loggermiddleware.LogMiddleware(logger))
	r.Use(compressormiddleware.ZipMiddleware)

	if key != "" {
		r.Use(hashmiddleware.HashMiddleware(logger, key))
	}

	r.Mount("/debug", middleware.Profiler())

	r.Post("/update/{metric}/{key}/{value}", updatemetric.Handler(logger, db))
	r.Post("/update/", updatejson.Handler(logger, db))
	r.Post("/value/", updatejson.Handler(logger, db))
	r.Post("/updates/", updatesmetrics.Handler(logger, db))

	r.Get("/value/{metric}/{key}", getmetric.Handler(logger, db))
	r.Get("/", getallmetrics.Handler(logger, db))
	r.Get("/ping", dbping.Handler(logger, db))

	return r
}
