package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"log"
	"net/http"
	"os"
)

var flagRunPort string

func MetricRouter() chi.Router {
	//create storage
	m := storage.NewMemStorage()
	//get handlers
	handler := handlers.NewMetricHandler(m)
	//create router
	r := chi.NewRouter()
	//set routes
	r.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
	r.Get("/value/{metric}/{key}", handler.GetMetricValue)
	r.Get("/", handler.MainPage)
	return r
}
func main() {
	//-----------------parse flags and env variables-----------------
	flag.StringVar(&flagRunPort, "a", ":8080", "port to run server on")
	flag.Parse()
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunPort = envAddress
	}
	//-----------------start server-----------------
	fmt.Println("Running server on", flagRunPort)
	log.Fatal(http.ListenAndServe(flagRunPort, MetricRouter()))
}
