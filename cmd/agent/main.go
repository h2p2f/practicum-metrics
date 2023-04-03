package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	"practicum-metrics/internal/client/metrics"
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
	var r, p int
	//parse flags
	flag.StringVar(&flagRunPort, "a", ":8080", "port to run server on")
	//DurationVar is not working, so I use IntVar with conversion to Duration
	flag.IntVar(&r, "r", 10, "report to server interval in seconds")
	flag.IntVar(&p, "p", 2, "pool interval in seconds")
	flag.Parse()
	//set reportInterval

	reportInterval = time.Duration(r)
	//set poolInterval
	poolInterval = time.Duration(p)
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
