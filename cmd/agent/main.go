package main

import (
	"bytes"
	"compress/gzip"
	//"errors"
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/h2p2f/practicum-metrics/internal/client/metrics"
	"log"
	"os"
	"strconv"
	"strings"
	//"syscall"
	"time"
	"unicode"
)

// these variables for start up flags
var flagRunPort string
var reportInterval time.Duration
var poolInterval time.Duration

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// function to monitor metrics
func getMetrics(m *metrics.RuntimeMetrics, pool time.Duration) {
	for {
		m.Monitor()
		time.Sleep(pool * time.Second)
	}
}
func main() {
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
	//print info
	//TODO: add normal logging
	fmt.Println("Running agent for server:", host)
	fmt.Println("Report to server interval:", reportInterval)
	fmt.Println("Pool interval:", poolInterval)
	//create new metrics
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()
	//start monitoring (made with goroutine, because interval is not constant)
	go getMetrics(m, poolInterval)
	//start reporting in main goroutine
	for {
		//we sleep here, because we need to report metrics after poolInterval
		time.Sleep(reportInterval * time.Second)
		//get metrics in json format
		jsonMetrics := m.JSONMetrics()
		//prepare metrics to send
		for _, data := range jsonMetrics {
			//compress data, this comment wrote captain obvious
			buf, err := Compress(data)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			//send data to server
			client := resty.New()
			//some autotests can be faster than server starts, so we need to retry three times, not more :)
			client.SetRetryCount(3).SetRetryWaitTime(200 * time.Millisecond)
			//upgrading request's headers
			resp, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Content-Encoding", "gzip").
				SetBody(buf).
				Post(host + "/update/")
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			//TODO: add normal logging
			//print response status code
			fmt.Println("received response from server: ", resp.StatusCode())

		}

	}

}

// function Compress to compress data
func Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
