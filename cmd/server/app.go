package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/database"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"github.com/jackc/pgx"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
)

// getFlagsAndEnv is a function that returns flags and env variables
func getFlagsAndEnv() (string, time.Duration, string, bool, string, bool, bool) {
	var (
		flagRunAddr       string
		flagStoreInterval time.Duration
		flagStorePath     string
		flagRestore       bool
		//interval          int
		IntervalDuration time.Duration
		databaseVar      string
		useDatabase      bool
		useFile          bool
	)
	useFile = false
	useDatabase = false

	// parse flags
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "port to run server on")
	//flag.IntVar(&interval, "i", 300, "interval to store metrics in seconds")
	flag.StringVar(&flagStorePath, "f", "/tmp/devops-metrics-db.json", "path to store metrics")
	flag.DurationVar(&IntervalDuration, "i", 300*time.Second, "interval to store metrics in seconds")
	flag.BoolVar(&flagRestore, "r", true, "restore metrics from file")
	flag.StringVar(&databaseVar, "d",
		"postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
		"databaseVar to store metrics")

	//postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable
	//host=localhost user=practicum password=yandex dbname=postgres sslmode=disable
	flag.Parse()
	// convert int to duration
	//flagStoreInterval = time.Duration(interval) * time.Second
	// get env variables, if they exist drop flags
	flagStoreInterval = IntervalDuration
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	if isFlagPassed("f") {
		useFile = true
	}
	if isFlagPassed("d") {
		useDatabase = true
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		var err error
		if isNumeric(envStoreInterval) {
			envStoreInterval = envStoreInterval + "s"

		}
		flagStoreInterval, err = time.ParseDuration(envStoreInterval)
		if err != nil {
			fmt.Println(err)
		}
	}
	if envStorePath := os.Getenv("FILE_STORAGE_PATH"); envStorePath != "" {
		flagStorePath = envStorePath
		useFile = true
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		envRestore, err := strconv.ParseBool(envRestore)
		if err != nil {
			log.Println(err)
		}
		flagRestore = envRestore
	}
	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		databaseVar = envDatabase
		useDatabase = true
	}
	//check if port is numeric - some people can try to run agent on :8080 - but it will be localhost:8080
	host := "localhost:"
	if isNumeric(flagRunAddr) {
		flagRunAddr = host + flagRunAddr
		fmt.Println("Running server on", flagRunAddr)
	}

	fmt.Println(sql.Drivers())
	fmt.Println(pgx.TextFormatCode)
	return flagRunAddr, flagStoreInterval, flagStorePath, flagRestore, databaseVar, useDatabase, useFile
}

// isNumeric is a function that checks if string contains only digits
func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// MetricRouter function to create chi router
func MetricRouter(m *storage.MemStorage, db *database.PGDB) chi.Router {
	//get handlers
	handler := handlers.NewMetricHandler(m, db)
	//create router
	r := chi.NewRouter()
	//add middlewares
	loggedAndZippedRouter := r.With(logger.WithLogging, handlers.GzipHanle)
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
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
