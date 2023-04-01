package main

import (
	"net/http"
	"practicum-metrics/internal/client/metrics"
	"time"
)

func main() {
	host := "http://localhost:8080"
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()
	timeCounter := 1

	for {
		m.Monitor()

		time.Sleep(2 * time.Second)
		timeCounter++
		if timeCounter%5 == 0 {
			urls := m.UrlMetrics(host)
			for _, url := range urls {
				req, err := http.NewRequest("POST", url, nil)
				if err != nil {
					panic(err)
				}
				client := &http.Client{}

				req.Header.Set("Content-Type", "text/plain")
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
			}
		}

	}

}
