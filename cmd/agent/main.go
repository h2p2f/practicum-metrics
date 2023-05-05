package main

import (
	"github.com/h2p2f/practicum-metrics/internal/client/config"
	"github.com/h2p2f/practicum-metrics/internal/client/httpclient"
	"github.com/h2p2f/practicum-metrics/internal/client/metrics"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"log"
	"time"
)

// function to monitor metrics
func getMetrics(m *metrics.RuntimeMetrics, pool time.Duration) {
	for {
		m.Monitor()
		time.Sleep(pool * time.Second)
	}
}
func main() {
	//init logger
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}

	//setup new config
	conf := config.NewConfig()
	conf.SetConfig(GetFlagAndEnvClient())

	//print config
	logger.Log.Sugar().Info("Starting agent")
	logger.Log.Sugar().Info("Running agent for server: %s ", conf.Address)
	logger.Log.Sugar().Info("Report to server interval: %d ", conf.ReportInterval)
	logger.Log.Sugar().Info("Pool interval: %d ", conf.PoolInterval)

	//create new metrics
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()

	//start monitoring (made with goroutine, because interval is not constant)
	go getMetrics(m, conf.PoolInterval)

	//start reporting in main goroutine
	for {
		//we sleep here, because we need to report metrics after poolInterval
		time.Sleep(conf.ReportInterval * time.Second)
		//get metrics in json format
		jsonMetrics := m.JSONMetrics()
		//prepare metrics to send
		err := httpclient.SendMetrics(jsonMetrics, conf.Address)
		if err != nil {
			logger.Log.Sugar().Errorf("Error sending metrics: %s", err)
		}
	}
}
