package app

import (
	"errors"
	"github.com/h2p2f/practicum-metrics/internal/client/config"
	"github.com/h2p2f/practicum-metrics/internal/client/hash"
	"github.com/h2p2f/practicum-metrics/internal/client/httpclient"
	"github.com/h2p2f/practicum-metrics/internal/client/metrics"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"go.uber.org/zap"
	"sync"
	"syscall"
	"time"
)

// getRuntimeMetrics - function to monitor runtime metrics
func getRuntimeMetrics(m *metrics.RuntimeMetrics, pool time.Duration) {
	for {
		m.RuntimeMonitor()
		time.Sleep(pool)
	}
}

// getGopsUtilMetrics - function to monitor gops metrics
func getGopsUtilMetrics(m *metrics.RuntimeMetrics, pool time.Duration) {
	for {
		m.GopsUtilMonitor()
		time.Sleep(pool)
	}
}

// SendBatchMetrics - function to send all metrics with one request
func SendBatchMetrics(m *metrics.RuntimeMetrics, reportInterval time.Duration, address string, key string) {
	for {
		// wait for report interval
		time.Sleep(reportInterval)
		// get all metrics
		jsonMetrics := m.JSONMetrics()
		// check if key is presented
		var checkSum [32]byte
		if key != "" {
			var err error
			// get hash of request data
			checkSum, err = hash.GetHash(key, jsonMetrics)
			if err != nil {
				logger.Log.Sugar().Errorf("key not presented: %s", err)
			}
		}
		// send metrics
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

// SendOneMetric - function to send one metric with one request
func SendOneMetric(conf *config.Config, mCh <-chan []byte, done chan<- bool) {
	// wait for metric
	for m := range mCh {
		// check if key is presented
		var checkSum [32]byte
		if conf.Key != "" {
			var err error
			// get hash of request data
			checkSum, err = hash.GetHash(conf.Key, m)
			if err != nil {
				logger.Log.Sugar().Errorf("key not presented: %s", err)
			}
		}
		// send metric
		err := httpclient.SendMetric(m, checkSum, conf.Address)
		if err != nil {
			logger.Log.Sugar().Errorf("Error sending metrics: %s", err)
			// if broken pipe - reconnect
			if errors.Is(err, syscall.EPIPE) {
				logger.Log.Sugar().Errorf("Broken pipe, reconnecting...")
				time.Sleep(1 * time.Second)
			}
		}
		// send done signal
		done <- true
	}
}

// Run - function to run agent
func Run(logger *zap.Logger) {
	// get config
	conf := config.NewConfig()
	conf.SetConfig()
	// log start up info
	logger.Sugar().Info("Starting agent")
	logger.Sugar().Infof("Running agent for server: %s ", conf.Address)
	logger.Sugar().Infof("Report to server interval: %s ", conf.ReportInterval)
	logger.Sugar().Infof("Pool interval: %s ", conf.PoolInterval)
	logger.Sugar().Infof("Send batch metrics: %t ", conf.Batch)
	if conf.Key != "" {
		logger.Sugar().Info("calculate checksum - true")
	}
	// init metrics
	m := new(metrics.RuntimeMetrics)
	m.NewMetrics()
	// add wait group
	var wg sync.WaitGroup
	wg.Add(3 + conf.RateLimit)
	// run metrics collectors
	go getRuntimeMetrics(m, conf.PoolInterval)
	go getGopsUtilMetrics(m, conf.PoolInterval)
	// send metrics
	//if batch is true - send all metrics with one request
	if conf.Batch {
		SendBatchMetrics(m, conf.ReportInterval, conf.Address, conf.Key)
	}
	//if batch is false - send one metric with one request
	if !conf.Batch {
		for {
			// wait for report interval
			time.Sleep(conf.ReportInterval)
			// get all metrics
			jsonMetrics := m.JSONMetricsForSingleSending()
			// make jobs buffered channel
			jobs := make(chan []byte, len(jsonMetrics))
			// make channel for done signal
			done := make(chan bool, len(jsonMetrics))
			// send metrics with goroutines
			// number of goroutines = rate limit
			for work := 0; work < conf.RateLimit; work++ {
				go SendOneMetric(conf, jobs, done)
			}
			// send metrics to jobs channel
			for _, metric := range jsonMetrics {
				jobs <- metric
			}
			// close jobs channel
			close(jobs)
			// wait for done signal
			for range jsonMetrics {
				<-done
			}
			// complete wait group
			wg.Done()
		}
	}
	// wait for all goroutines
	wg.Wait()
}
