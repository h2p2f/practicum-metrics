package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
)

// getFlagsAndEnv is a function that returns flags and env variables
func getFlagsAndEnv() (string, time.Duration, string, bool, string) {
	var (
		flagRunAddr       string
		flagStoreInterval time.Duration
		flagStorePath     string
		flagRestore       bool
		interval          int
		database          string
	)
	// parse flags
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "port to run server on")
	flag.IntVar(&interval, "i", 300, "interval to store metrics in seconds")
	flag.StringVar(&flagStorePath, "f", "/tmp/devops-metrics-db.json", "path to store metrics")
	flag.BoolVar(&flagRestore, "r", true, "restore metrics from file")
	flag.StringVar(&database, "d",
		"host=localhost user=practicum password=yandex dbname=postgres sslmode=disable",
		"database to store metrics")
	flag.Parse()
	// convert int to duration
	flagStoreInterval = time.Duration(interval)
	// get env variables, if they exist drop flags
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
	if envDatabase := os.Getenv("DATABASE_DSN"); envDatabase != "" {
		database = envDatabase
	}
	//check if port is numeric - some people can try to run agent on :8080 - but it will be localhost:8080
	host := "localhost:"
	if isNumeric(flagRunAddr) {
		flagRunAddr = host + flagRunAddr
		fmt.Println("Running server on", flagRunAddr)
	}
	return flagRunAddr, flagStoreInterval, flagStorePath, flagRestore, database
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
