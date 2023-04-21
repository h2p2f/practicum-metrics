package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"log"
	"net/http"
	"os"
)

// flagRunAddr is a param for run address
var flagRunAddr string

// MetricRouter function to create router
func MetricRouter() chi.Router {
	//create storage
	m := storage.NewMemStorage()
	//get handlers
	handler := handlers.NewMetricHandler(m)
	//create router
	r := chi.NewRouter()
	//set routes
	loggedRouter := r.With(logger.WithLogging)
	loggedRouter.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
	loggedRouter.Get("/value/{metric}/{key}", handler.GetMetricValue)
	loggedRouter.Get("/", handler.MainPage)
	loggedRouter.Post("/update", handler.UpdateJSON)
	loggedRouter.Post("/value", handler.ValueJSON)
	return r
}
func main() {
	//-----------------parse flags and env variables-----------------
	// this code for normal server users
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "port to run server on")
	flag.Parse()

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	//logger.Log.Info("Server started", zap.String("address", flagRunAddr))
	logger.Log.Sugar().Infof("Server started on %s", flagRunAddr)
	//-----------------start server-----------------
	//fmt.Println("Running server on", flagRunAddr)
	log.Fatal(http.ListenAndServe(flagRunAddr, MetricRouter()))
}
