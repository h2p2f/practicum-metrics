package main

import (
	"errors"
	"github.com/h2p2f/practicum-metrics/internal/client/hash"
	"log"
	"syscall"
	"time"

	"github.com/h2p2f/practicum-metrics/internal/client/config"
	"github.com/h2p2f/practicum-metrics/internal/client/httpclient"
	"github.com/h2p2f/practicum-metrics/internal/client/metrics"
	"github.com/h2p2f/practicum-metrics/internal/logger"
)

// function to monitor metrics
func getMetrics(m *metrics.RuntimeMetrics, pool time.Duration) {
	for {
		m.Monitor()
		time.Sleep(pool)
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
	logger.Log.Sugar().Infof("Running agent for server: %s ", conf.Address)
	logger.Log.Sugar().Infof("Report to server interval: %s ", conf.ReportInterval)
	logger.Log.Sugar().Infof("Pool interval: %s ", conf.PoolInterval)
	logger.Log.Sugar().Infof("Key to calculate checksum: %s ", conf.Key)

	//create new metrics
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()

	//start monitoring (made with goroutine, because interval is not constant)
	go getMetrics(m, conf.PoolInterval)

	//start reporting in main goroutine
	for {
		//we sleep here, because we need to report metrics after poolInterval
		time.Sleep(conf.ReportInterval)
		//get metrics in json format
		jsonMetrics := m.JSONMetrics()
		var checkSum [32]byte
		if conf.Key != "" {
			var err error
			checkSum, err = hash.GetHash(conf.Key, jsonMetrics)
			if err != nil {
				logger.Log.Sugar().Errorf("key not presented: %s", err)
			}
		}
		//prepare metrics to send
		//if it needs to send metrics by one - uncomment next line and comment line 56
		//err := httpclient.SendMetrics(jsonMetrics, conf.Address)
		err := httpclient.SendBatchMetrics(jsonMetrics, checkSum, conf.Address)
		if err != nil {
			logger.Log.Sugar().Errorf("Error sending metrics: %s", err)
			//if broken pipe - reconnect
			//this code for increment 13, but setRetryCount
			//handles it well without any additional implementation
			//like "see what I can do" :)
			if errors.Is(err, syscall.EPIPE) {
				logger.Log.Sugar().Errorf("Broken pipe, reconnecting...")
				time.Sleep(1 * time.Second)
			}
		}
	}
}
