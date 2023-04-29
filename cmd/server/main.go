package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"log"
	"net/http"
	"os"
	"time"
)

// MetricRouter function to create chi router
func MetricRouter(m *storage.MemStorage) chi.Router {
	//get handlers
	handler := handlers.NewMetricHandler(m)
	//create router
	r := chi.NewRouter()
	//add middlewares
	loggedAndZippedRouter := r.With(logger.WithLogging, handlers.GzipHanle)
	//add routes
	loggedAndZippedRouter.Post("/update/", handler.UpdateJSON)
	loggedAndZippedRouter.Post("/value/", handler.ValueJSON)
	loggedAndZippedRouter.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
	loggedAndZippedRouter.Get("/value/{metric}/{key}", handler.GetMetricValue)
	loggedAndZippedRouter.Get("/", handler.MainPage)

	return r
}

func main() {
	//setup new config
	conf := config.NewConfig()
	//set config from flags and env
	conf.SetConfig(getFlagsAndEnv())

	//create storage
	m := storage.NewMemStorage()

	//create fileDB with path and interval from config
	fileDB := storage.NewFileDB(conf.PathToStoreFile, conf.StoreInterval)
	//hardcode to turn off restore from file -  fix bug iter8 autotest
	_, err := os.Stat(conf.PathToStoreFile)
	if err != nil {
		conf.Restore = false
	}
	//restore metrics from file if flag is set
	if conf.Restore {
		metrics, err := fileDB.ReadFromFile()
		if err != nil {
			fmt.Println(err)
		}
		m.RestoreMetrics(metrics)
	}
	//save metrics to file with interval from config
	//made with anonymous function and goroutine
	go func() {
		for {
			time.Sleep(conf.StoreInterval * time.Second)
			metrics := m.GetAllMetricsSliced()
			err := fileDB.SaveToFile(metrics)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	//init logger
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	//start server with router
	logger.Log.Sugar().Infof("Server started on %s", conf.ServerAddress)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, MetricRouter(m)))
}
