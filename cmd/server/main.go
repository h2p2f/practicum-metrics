package main

import (
	"net/http"
	"practicum-metrics/internal/server/handlers"
	"practicum-metrics/internal/storage"
)

func main() {

	//Storage starts here
	m := storage.NewMemStorage()
	handler := handlers.NewMetricHandler(m)

	//Handler starts here
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.MainPage)

	//Listener starts here
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
