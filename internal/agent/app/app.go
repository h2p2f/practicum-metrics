package app

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/h2p2f/practicum-metrics/internal/agent/config"
	hash2 "github.com/h2p2f/practicum-metrics/internal/agent/hash"
	"github.com/h2p2f/practicum-metrics/internal/agent/httpclient"
	"github.com/h2p2f/practicum-metrics/internal/agent/storage"
	"github.com/h2p2f/practicum-metrics/internal/logger"
)

func SendOneMetric(logger *zap.Logger, config *config.AgentConfig, mCh <-chan []byte, done chan<- bool) {
	// wait for metric
	for m := range mCh {
		// check if key is presented
		var hash [32]byte
		if config.Key != "" {
			// get hash of request data
			hash = hash2.GetHash(config.Key, m)
		}
		// send metric
		err := httpclient.SendMetric(logger, m, hash, config.ServerAddress)
		if err != nil {
			logger.Error("Error sending metrics: %s",
				zap.Error(err))
			// if broken pipe - reconnect
			if errors.Is(err, syscall.EPIPE) {
				logger.Error("Broken pipe, reconnecting...")
				time.Sleep(1 * time.Second)
			}
		}
		// send done signal
		done <- true
	}
}

func getRuntimeMetrics(m *storage.MetricStorage, pool time.Duration) {
	for {
		m.RuntimeMetricsMonitor()
		time.Sleep(pool)
	}
}

func getGopsUtilMetrics(m *storage.MetricStorage, pool time.Duration) {
	for {
		m.GopsUtilizationMonitor()
		time.Sleep(pool)
	}
}

func Run() {
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
		return
	}
	logger.Log.Info("Starting agent...")
	logger.Log.Info("Reading config...")

	conf := config.GetConfig()

	fields := []zapcore.Field{
		zap.Int("rate limit", conf.RateLimit),
		zap.String("poll interval", conf.PollInterval.String()),
		zap.String("report interval", conf.ReportInterval.String()),
		zap.String("server address", conf.ServerAddress),
	}
	if conf.Key != "" {
		msg := "key is presented"
		fields = append(fields, zap.String("key", msg))
	}
	logger.Log.Info("Config loaded", fields...)

	memDB := storage.NewAgentStorage()

	var wg sync.WaitGroup
	wg.Add(4 + conf.RateLimit)

	go getRuntimeMetrics(memDB, conf.PollInterval)
	go getGopsUtilMetrics(memDB, conf.PollInterval)

	if conf.RateLimit > 0 {
		for {
			time.Sleep(conf.ReportInterval)
			data := memDB.JSONMetrics()
			jobs := make(chan []byte, len(data))
			done := make(chan bool, len(data))
			for w := 1; w <= conf.RateLimit; w++ {
				go SendOneMetric(logger.Log, conf, jobs, done)
			}
			for _, metric := range data {
				jobs <- metric
			}
			close(jobs)
			for range data {
				<-done
			}
			wg.Done()
		}
	}
	if conf.RateLimit <= 0 {
		for {
			time.Sleep(conf.ReportInterval)
			data := memDB.BatchJSONMetrics()
			var hash [32]byte
			if conf.Key != "" {
				hash = hash2.GetHash(conf.Key, data)
			}
			err := httpclient.SendBatchJSONMetrics(logger.Log, conf, data, hash)
			if err != nil {
				logger.Log.Sugar().Errorf("Error sending metrics: %s", err)
			}
		}
	}
	wg.Wait() //nolint:govet
}
