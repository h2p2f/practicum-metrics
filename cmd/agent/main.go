package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"practicum-metrics/internal/client/metrics"
	"time"
)

func main() {
	//set host
	host := "http://localhost:8080"
	//create new metrics
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()
	//start monitoring
	timeCounter := 1
	for {
		m.Monitor()
		time.Sleep(2 * time.Second)
		timeCounter++
		//send metrics to server every 10 seconds
		if timeCounter%5 == 0 {
			urls := m.UrlMetrics(host)
			for _, url := range urls {
				//send metrics to server with resty
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
}
