package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// function check if string is numeric
func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// GetFlagAndEnvClient is a function that returns flag and env variables
func GetFlagAndEnvClient() (string, string, time.Duration, time.Duration) {
	var flagRunPort string
	var reportInterval time.Duration
	var poolInterval time.Duration
	var key string
	var r, p int

	//------------------flags and env variables------------------

	flag.StringVar(&flagRunPort, "a", "localhost:8080", "port to run server on")
	flag.IntVar(&r, "r", 10, "report to server interval in seconds")
	flag.IntVar(&p, "p", 2, "pool interval in seconds")
	//flag.DurationVar(&reportInterval, "r", 10*time.Second, "report to server interval in seconds")
	//flag.DurationVar(&poolInterval, "p", 2*time.Second, "pool interval in seconds")
	flag.StringVar(&key, "k", "", "key to calculate data's hash if presented")
	flag.Parse()
	//convert int to duration
	reportInterval = time.Duration(r)
	//set poolInterval
	poolInterval = time.Duration(p)
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
	//------------------start agent------------------
	//set host

	host := "http://"
	//check if port is numeric - some people can try to run agent on :8080 - but it will be localhost:8080
	if isNumeric(flagRunPort) {
		host = host + "localhost:" + flagRunPort
	} else if !strings.Contains(flagRunPort, host) {
		host += flagRunPort
	}
	return host, key, reportInterval, poolInterval
}
