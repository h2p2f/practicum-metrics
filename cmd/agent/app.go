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

func GetFlagAndEnvClient() (string, time.Duration, time.Duration) {
	var flagRunPort string
	var reportInterval time.Duration
	var poolInterval time.Duration

	//------------------flags and env variables------------------
	//temporary local variables for flags
	//this code has no grace, but it works
	var r, p int
	//parse flags
	flag.StringVar(&flagRunPort, "a", "localhost:8080", "port to run server on")
	//TODO: fix this shitcode
	flag.IntVar(&r, "r", 10, "report to server interval in seconds")
	flag.IntVar(&p, "p", 2, "pool interval in seconds")
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
	//------------------start agent------------------
	//set host

	host := "http://"
	//check if port is numeric - some people can try to run agent on :8080 - but it will be localhost:8080
	if isNumeric(flagRunPort) {
		host = host + "localhost:" + flagRunPort
	} else if !strings.Contains(flagRunPort, host) {
		host += flagRunPort
	}
	return host, reportInterval, poolInterval
}
