package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode"
)

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// MetricRouter function to create router
func MetricRouter(m *storage.MemStorage) chi.Router {
	//create storage

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

func getFlagsAndEnv() (string, time.Duration, string, bool) {
	var (
		flagRunAddr       string
		flagStoreInterval time.Duration
		flagStorePath     string
		flagRestore       bool
		interval          int
	)

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "port to run server on")
	flag.IntVar(&interval, "i", 300, "interval to store metrics in seconds")
	flag.StringVar(&flagStorePath, "f", "/tmp/devops-metrics-db.json", "path to store metrics")
	flag.BoolVar(&flagRestore, "r", true, "restore metrics from file")
	flag.Parse()

	flagStoreInterval = time.Duration(interval)

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		envStoreInterval, err := strconv.Atoi(envStoreInterval)
		if err != nil {
			log.Fatal(err)
		}
		flagStoreInterval = time.Duration(envStoreInterval)
	}
	if envStorePath := os.Getenv("STORE_FILE"); envStorePath != "" {
		flagStorePath = envStorePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		envRestore, err := strconv.ParseBool(envRestore)
		if err != nil {
			log.Fatal(err)
		}
		flagRestore = envRestore
	}
	//hardcode for autotests
	host := "localhost:"
	if isNumeric(flagRunAddr) {
		flagRunAddr = host + flagRunAddr
		fmt.Println("Running server on", flagRunAddr)
	}
	fmt.Println(flagRunAddr, flagStoreInterval, flagStorePath, flagRestore)
	return flagRunAddr, flagStoreInterval, flagStorePath, flagRestore
}
func main() {
	conf := config.NewConfig()
	conf.SetConfig(getFlagsAndEnv())

	m := storage.NewMemStorage()

	fileDB := storage.NewFileDB(conf.PathToStoreFile, conf.StoreInterval)

	if conf.Restore {
		metrics, err := fileDB.ReadFromFile()
		if err != nil {
			fmt.Println(err)
		}
		m.RestoreMetrics(metrics)
	}

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

	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	logger.Log.Sugar().Infof("Server started on %s", conf.ServerAddress)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, MetricRouter(m)))
}
