// Package app launches the main agent logic - initializes the logger,
// reads the config, launches the metrics monitoring and sends them to the server.
package app

import (
	"context"
	"errors"
	"log"
	_ "net/http/pprof"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/h2p2f/practicum-metrics/internal/agent/config"
	hash2 "github.com/h2p2f/practicum-metrics/internal/agent/hash"
	"github.com/h2p2f/practicum-metrics/internal/agent/httpclient"
	"github.com/h2p2f/practicum-metrics/internal/agent/storage"
)

// SendOneMetric sends metrics to the server in a worker pool one metric at a time
func SendOneMetric(logger *zap.Logger, config *config.AgentConfig, mCh <-chan []byte, done chan<- bool) {

	// wait for metric
	for m := range mCh {

		// check if key is presented
		var hash [32]byte

		// get hash of request data
		if config.Key != "" {
			// get hash of request data
			hash = hash2.GetHash(m)
		}

		// send metric
		err := httpclient.SendMetric(logger, m, hash, config)
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

// getRuntimeMetrics launches memory metrics monitoring
func getRuntimeMetrics(ctx context.Context, m *storage.MetricStorage, poolTime time.Duration) {
	t := time.NewTicker(poolTime)
	select {
	case <-ctx.Done():
		t.Stop()
		return
	case <-t.C:
		m.RuntimeMetricsMonitor()
	}
}

// getGopsUtilMetrics launches gops metrics monitoring
func getGopsUtilMetrics(ctx context.Context, m *storage.MetricStorage, poolTime time.Duration) {
	t := time.NewTicker(poolTime)
	select {
	case <-ctx.Done():
		t.Stop()
		return
	case <-t.C:
		m.GopsUtilizationMonitor()
	}
}

type App struct {
	db     *storage.MetricStorage
	config *config.AgentConfig
	logger *zap.Logger
}

// Run launches the agent
func Run(sigint chan os.Signal, connectionsClosed chan<- struct{}) {
	// read config
	conf, logger, err := config.GetConfig()
	if err != nil {
		log.Fatal("Failed to read config")
	}
	fields := []zapcore.Field{
		zap.Int("rate limit", conf.RateLimit),
		zap.String("poll interval", conf.PollInterval.String()),
		zap.String("report interval", conf.ReportInterval.String()),
		zap.String("server address", conf.ServerAddress),
		zap.String("log level", conf.LogLevel),
		zap.String("key file", conf.KeyFile),
		zap.String("ip address", conf.IPaddr.String()),
	}

	// if the key is not empty - add a message to the log
	if conf.Key != "" {
		msg := "key is presented"
		fields = append(fields, zap.String("key", msg))
	}
	logger.Info("Config loaded", fields...)

	// initialize storage
	memDB := storage.NewAgentStorage()

	app := App{
		db:     memDB,
		config: conf,
		logger: logger,
	}

	// create context
	ctx, cancel := context.WithCancel(context.Background()) //nolint:govet
	defer cancel()

	// start metrics monitoring
	go getRuntimeMetrics(ctx, memDB, conf.PollInterval)
	go getGopsUtilMetrics(ctx, memDB, conf.PollInterval)

	// start sending metrics to the server depending on the limit
	if conf.RateLimit > 0 {
		go app.sendWithRateLimit(ctx)
	} else {
		// if the limit is not set - send metrics in batches in json to the server
		go app.sendWithoutRateLimit(ctx)
	}

	// wait for done signal

	for signal := range sigint {
		logger.Info("Received signal", zap.String("signal", signal.String()))
		logger.Info("Shutting down agent...")
		// stop sending metrics
		cancel()
		close(sigint)
		logger.Info("Agent shutdown gracefully")
		close(connectionsClosed)
	}
}

// sendWithRateLimit sends metrics to the server one at a time in a goroutine pool,
// the number of which is limited by the RateLimit parameter
func (app *App) sendWithRateLimit(ctx context.Context) {
	app.logger.Info("Sending metrics with rate limit")
	t := time.NewTicker(app.config.ReportInterval)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			app.logger.Info("tick")
			data := app.db.JSONMetrics()
			app.logger.Info("Sending metrics to the server in json one metric at a time")
			// create channels for workers
			jobs := make(chan []byte, len(data))
			done := make(chan bool, len(data))
			// start workers
			for w := 1; w <= app.config.RateLimit; w++ {
				go SendOneMetric(app.logger, app.config, jobs, done)
			}
			// send metrics to channel
			for _, metric := range data {
				jobs <- metric
			}
			// close channels
			close(jobs)
			// wait for workers to finish
			for range data {
				<-done
			}
			close(done)
		}
	}
}

// sendWithoutRateLimit sends metrics to the server in batches in json
func (app *App) sendWithoutRateLimit(ctx context.Context) {
	app.logger.Info("Sending metrics to the server in batches in json")
	t := time.NewTicker(app.config.ReportInterval)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			data := app.db.BatchJSONMetrics()
			// calculate hash
			var hash [32]byte
			if app.config.Key != "" {
				hash = hash2.GetHash(data)
			}
			// send metrics
			err := httpclient.SendBatchJSONMetrics(app.logger, app.config, data, hash)
			if err != nil {
				app.logger.Error("Error sending metrics: ", zap.Error(err))
			}
		}
	}

}
