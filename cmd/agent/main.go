package main

import (
	"errors"
	"github.com/h2p2f/practicum-metrics/internal/client/hash"
	"github.com/h2p2f/practicum-metrics/internal/client/httpclient"
	"log"
	"sync"
	"syscall"
	"time"

	"github.com/h2p2f/practicum-metrics/internal/client/config"
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

func SendBatchMetrics(m *metrics.RuntimeMetrics, reportInterval time.Duration, address string, key string) {
	for {
		time.Sleep(reportInterval)
		jsonMetrics := m.JSONMetrics()
		var checkSum [32]byte
		if key != "" {
			var err error
			checkSum, err = hash.GetHash(key, jsonMetrics)
			if err != nil {
				logger.Log.Sugar().Errorf("key not presented: %s", err)
			}
		}
		err := httpclient.SendBatchMetrics(jsonMetrics, checkSum, address)
		if err != nil {
			logger.Log.Sugar().Errorf("Error sending metrics: %s", err)
			if errors.Is(err, syscall.EPIPE) {
				logger.Log.Sugar().Errorf("Broken pipe, reconnecting...")
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func SendOneMetric(conf *config.Config, mCh <-chan []byte, done chan<- bool) {
	for m := range mCh {
		var checkSum [32]byte
		if conf.Key != "" {
			var err error
			checkSum, err = hash.GetHash(conf.Key, m)
			if err != nil {
				logger.Log.Sugar().Errorf("key not presented: %s", err)
			}
		}
		err := httpclient.SendMetric(m, checkSum, conf.Address)
		if err != nil {
			logger.Log.Sugar().Errorf("Error sending metrics: %s", err)
			if errors.Is(err, syscall.EPIPE) {
				logger.Log.Sugar().Errorf("Broken pipe, reconnecting...")
				time.Sleep(1 * time.Second)
			}
		}
		done <- true
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
	logger.Log.Sugar().Infof("Send batch metrics: %t ", conf.Batch)
	if conf.Key != "" {
		logger.Log.Sugar().Info("calculate checksum - true")
	}

	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()

	var wg sync.WaitGroup
	wg.Add(2 + conf.RateLimit)

	go getMetrics(m, conf.PoolInterval)

	if conf.Batch {
		SendBatchMetrics(m, conf.ReportInterval, conf.Address, conf.Key)
	}
	if !conf.Batch {
		for {
			time.Sleep(conf.ReportInterval)
			jsonMetrics := m.JSONMetricsForSingleSending()
			jobs := make(chan []byte, len(jsonMetrics))
			done := make(chan bool, len(jsonMetrics))
			for work := 0; work < conf.RateLimit; work++ {
				go SendOneMetric(conf, jobs, done)
			}
			for _, metric := range jsonMetrics {
				jobs <- metric
			}
			close(jobs)
			for range jsonMetrics {
				<-done
			}
			wg.Done()
		}
	}
	wg.Wait()
}
