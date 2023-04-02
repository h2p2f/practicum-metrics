package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"practicum-metrics/internal/server/handlers"
	"practicum-metrics/internal/storage"
)

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
	//start server
	log.Fatal(http.ListenAndServe(":8080", MetricRouter()))
}
