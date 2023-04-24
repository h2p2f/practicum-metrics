package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"log"
	"net/http"
	"os"
	"unicode"
)

// flagRunAddr is a param for run address
var flagRunAddr string

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// MetricRouter function to create router
func MetricRouter() chi.Router {
	//create storage
	m := storage.NewMemStorage()
	//get handlers
	handler := handlers.NewMetricHandler(m)
	//create router
	r := chi.NewRouter()
	//set routes
	loggedAndZippedRouter := r.With(logger.WithLogging, handlers.GzipHanle)
	loggedAndZippedRouter.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
	loggedAndZippedRouter.Get("/value/{metric}/{key}", handler.GetMetricValue)
	loggedAndZippedRouter.Get("/", handler.MainPage)
	loggedAndZippedRouter.Post("/update/", handler.UpdateJSON)
	loggedAndZippedRouter.Post("/value/", handler.ValueJSON)
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
	//hardcode for autotests
	host := "localhost:"
	if isNumeric(flagRunAddr) {
		flagRunAddr = host + flagRunAddr
		fmt.Println("Running server on", flagRunAddr)
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
