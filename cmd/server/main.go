package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"practicum-metrics/internal/server/handlers"
	"practicum-metrics/internal/storage"
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
	flag.StringVar(&flagRunPort, "a", ":8080", "port to run server on")
	flag.Parse()
	//start server
	fmt.Println("Running server on", flagRunPort)
	log.Fatal(http.ListenAndServe(flagRunPort, MetricRouter()))
}
