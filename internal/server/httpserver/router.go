// Package httpserver реализует обработчики запросов к серверу.
//
// Package httpserver implements handlers for server requests.
package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver/middlewares/decryptormiddleware"
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

// DataBaser - интерфейс для работы с хранилищем данных.
//
// DataBaser is an interface for working with a data store.
type DataBaser interface {
	SetCounter(key string, value int64)
	SetGauge(key string, value float64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
	Ping() error
}

// DataBase - структура для работы с хранилищем данных.
//
// DataBase is a structure for working with a data store.
type DataBase struct {
	DataBaser
}

// NewDataBase - конструктор для DataBase.
//
// NewDataBase is a constructor for DataBase.
func NewDataBase(db DataBaser) *DataBase {
	return &DataBase{db}
}

// MetricRouter - конструктор для роутера.
//
// MetricRouter is a constructor for the router.
func MetricRouter(logger *zap.Logger, m DataBaser, config *config.ServerConfig) *chi.Mux {
	db := NewDataBase(m)
	r := chi.NewRouter()
	//регистрация middleware
	//
	// middleware registration
	r.Use(decryptormiddleware.DecryptMiddleware(config.PrivateKey))
	r.Use(loggermiddleware.LogMiddleware(logger))
	r.Use(compressormiddleware.ZipMiddleware)

	if config.Key != "" {
		r.Use(hashmiddleware.HashMiddleware(logger, config.Key))
	}
	//регистрация профайлера
	//
	// profiler registration
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
