package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/h2p2f/practicum-metrics/internal/client/metrics"
	"os"
	"strconv"
	"time"
)

var flagRunPort string
var reportInterval time.Duration
var poolInterval time.Duration

//function to monitor metrics
func getMetrics(m *metrics.RuntimeMetrics, pool time.Duration) {
	for {
		m.Monitor()
		time.Sleep(pool * time.Second)
	}
}
func main() {
	//------------------flags and env variables------------------
	//variables for flags
	var r, p int
	//parse flags
	flag.StringVar(&flagRunPort, "a", ":8080", "port to run server on")
	//DurationVar is not working, so I use IntVar with conversion to Duration. TODO: fix it
	flag.IntVar(&r, "r", 10, "report to server interval in seconds")
	flag.IntVar(&p, "p", 2, "pool interval in seconds")
	flag.Parse()
	//convert int to duration
	reportInterval = time.Duration(r)
	//set poolInterval
	poolInterval = time.Duration(p)
	//get env variables
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		envReportInterval, err := strconv.Atoi(envReportInterval)
		if err != nil {
			panic(err)
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
	host := "http://localhost" + flagRunPort
	fmt.Println("Running agent for server:", host)
	fmt.Println("Report to server interval:", reportInterval)
	fmt.Println("Pool interval:", poolInterval)
	//create new metrics
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()
	//start monitoring (made with goroutine, because use time.Sleep with poolInterval)
	go getMetrics(m, poolInterval)
	//start reporting in main goroutine
	for {
		//sleep for reportInterval
		time.Sleep(reportInterval * time.Second)
		//get slice urls
		urls := m.URLMetrics(host)
		for _, url := range urls {
			//send metrics to server with resty
			//create new http client
			client := resty.New()
			resp, err := client.R().
				SetHeader("Content-Type", "text/plain").
				Post(url)
			if err != nil {
				panic(err)
			}
			//I don't fully understand this code, but it works
			fmt.Print(resp)
		}
	}
}
