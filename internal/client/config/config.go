package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

// Config struct
type Config struct {
	Address        string
	Key            string
	ReportInterval time.Duration
	PoolInterval   time.Duration
	Batch          bool
	RateLimit      int
}

// NewConfig returns new Config struct with default values
func NewConfig() *Config {
	return &Config{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PoolInterval:   2 * time.Second,
		Key:            "",
		Batch:          true,
		RateLimit:      6,
	}
}

// SetConfig sets config values from flags or env variables
func (c *Config) SetConfig() *Config {

	var flagRunPort string
	var reportInterval time.Duration
	var poolInterval time.Duration
	var key string
	var r, p, rateLimit int
	var batch bool

	flag.StringVar(&flagRunPort, "a", "localhost:8080", "port to run server on")
	flag.IntVar(&r, "r", 10, "report to server interval in seconds")
	flag.IntVar(&p, "p", 2, "pool interval in seconds")
	//flag.DurationVar(&reportInterval, "r", 10*time.Second, "report to server interval in seconds")
	//flag.DurationVar(&poolInterval, "p", 2*time.Second, "pool interval in seconds")
	flag.StringVar(&key, "k", "", "key to calculate data's hash if presented")
	flag.BoolVar(&batch, "b", true, "batch mode")
	flag.IntVar(&rateLimit, "l", 6, "requests rate limit")

	flag.Parse()
	//convert int to duration
	reportInterval = time.Duration(r) * time.Second
	//set poolInterval
	poolInterval = time.Duration(p) * time.Second
	//get env variables, if they exist drop flags
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportInterval, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatal(err)
		}
		reportInterval = time.Duration(envReportInterval)
	}
	if envPoolInterval := os.Getenv("POOL_INTERVAL"); envPoolInterval != "" {
		envPoolInterval, err := strconv.Atoi(envPoolInterval)
		if err != nil {
			panic(err)
		}
		poolInterval = time.Duration(envPoolInterval)
	}
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		flagRunPort = envAddress
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		key = envKey
	}
	if envBatch := os.Getenv("BATCH"); envBatch != "" {
		batch = true
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		envRateLimit, err := strconv.Atoi(envRateLimit)
		if err != nil {
			panic(err)
		}
		rateLimit = envRateLimit
	}
	//set config values
	c.Address = flagRunPort
	c.ReportInterval = reportInterval
	c.PoolInterval = poolInterval
	c.Key = key
	c.Batch = batch
	c.RateLimit = rateLimit
	return c
}
