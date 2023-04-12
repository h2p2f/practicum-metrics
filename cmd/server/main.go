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

//flagRunAddr is a param for run address
var flagRunAddr string

//function to create router
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
	// this code for normal server users
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "port to run server on")
	flag.Parse()

	// this code below for crazy people who want to use random flags
	//sliceFlags := flag.NewFlagSet("slice", flag.ContinueOnError)
	//sliceFlags.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	//sliceFlags.Parse(os.Args[1:])
	//code below is feature for handle errors in sliceFlags.Parse
	//if err != nil {
	//	panic(err)
	//}

	//this code for a healthy user
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	//-----------------start server-----------------
	fmt.Println("Running server on", flagRunAddr)
	log.Fatal(http.ListenAndServe(flagRunAddr, MetricRouter()))
}
