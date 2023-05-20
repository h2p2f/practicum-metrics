package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
)

// Configurer interface for config
type Configurer interface {
	NewConfig() *serverConfig
	GetConfig() *serverConfig
	SetConfig(string, time.Duration, string, bool) *serverConfig
}

// serverConfig is a struct that contains server config
type serverConfig struct {
	// address to listen on
	ServerAddress string
	// interval to store metrics
	StoreInterval time.Duration
	// path to store file
	PathToStoreFile string
	// restore metrics from file
	Restore  bool
	Database string
	UseDB    bool
	UseFile  bool
	Key      string
}

// NewConfig is a function that returns a new config
func NewConfig() *serverConfig {
	return &serverConfig{
		ServerAddress:   "localhost:8080",
		StoreInterval:   300 * time.Second,
		PathToStoreFile: "/tmp/devops-metrics-db.json",
		Restore:         true,
		Database:        "host=localhost user=practicum password=yandex dbname=postgres sslmode=disable",
		UseDB:           false,
		UseFile:         false,
		Key:             "",
	}
}

// SetServerAddress is a function that sets server address

// GetConfig is a function that returns config
func (c *serverConfig) GetConfig() *serverConfig {
	return c
}

// SetConfig is a function that sets config
func (c *serverConfig) SetConfig(address string, interval time.Duration, path string, restore bool, db string, udb bool, uf bool, key string) *serverConfig {
	c.ServerAddress = address
	c.StoreInterval = interval
	c.PathToStoreFile = path
	c.Restore = restore
	c.Database = db
	c.UseDB = udb
	c.UseFile = uf
	c.Key = key
	return c
}

func GetFlagsAndEnv() (string, time.Duration, string, bool, string, bool, bool, string) {
	var (
		flagRunAddr       string
		flagStoreInterval time.Duration
		flagStorePath     string
		flagRestore       bool
		IntervalDuration  time.Duration
		databaseVar       string
		useDatabase       bool
		useFile           bool
		key               string
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
	flag.StringVar(&key, "k", "", "key to calculate data's hash")

	flag.Parse()

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
	if envKey := os.Getenv("KEY"); envKey != "" {
		key = envKey
	}
	//check if port is numeric - some people can try to run agent on :8080 - but it will be localhost:8080
	host := "localhost:"
	if isNumeric(flagRunAddr) {
		flagRunAddr = host + flagRunAddr
		fmt.Println("Running server on", flagRunAddr)
	}

	return flagRunAddr, flagStoreInterval, flagStorePath, flagRestore, databaseVar, useDatabase, useFile, key
}

func (c *serverConfig) GetKey() string {
	return c.Key
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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
