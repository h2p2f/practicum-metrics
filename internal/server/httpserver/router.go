package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
)

func MetricRouter(m Storager, db DataBaseHandler) chi.Router {
	//get httpserver
	handler := NewMetricHandler(m, db)
	//create router
	r := chi.NewRouter()
	//add middlewares
	loggedAndZippedRouter := r.With(logger.WithLogging, GzipHanle)
	loggedRouter := r.With(logger.WithLogging)
	//add routes
	loggedAndZippedRouter.Post("/update/", handler.UpdateJSON)
	loggedAndZippedRouter.Post("/value/", handler.ValueJSON)
	loggedAndZippedRouter.Post("/updates/", handler.UpdatesBatch)
	loggedRouter.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
	loggedRouter.Get("/value/{metric}/{key}", handler.GetMetricValue)
	loggedRouter.Get("/ping", handler.DBPing)
	loggedAndZippedRouter.Get("/", handler.MainPage)
	//
	return r
}
