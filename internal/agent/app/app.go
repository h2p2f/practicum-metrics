// Package app запускает основную логику агента - инициализирует логгер,
// читает конфиг, запускает мониторинг метрик и отправляет их на сервер.
//
// Package app launches the main agent logic - initializes the logger,
// reads the config, launches the metrics monitoring and sends them to the server.
package app

import (
	"errors"
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
)

// SendOneMetric отправляет метрики на сервер в пуле воркеров по одной метрике за раз
//
// SendOneMetric sends metrics to the server in a worker pool one metric at a time
func SendOneMetric(logger *zap.Logger, config *config.AgentConfig, mCh <-chan []byte, done chan<- bool) {
	// ждем метрику
	// wait for metric
	for m := range mCh {
		// проверяем, установлен ли ключ
		// check if key is presented
		var hash [32]byte
		// получаем хеш данных запроса
		// get hash of request data
		if config.Key != "" {
			// get hash of request data
			hash = hash2.GetHash(config.Key, m)
		}
		// отправляем метрику
		// send metric
		err := httpclient.SendMetric(logger, m, hash, config.ServerAddress)
		if err != nil {
			logger.Error("Error sending metrics: %s",
				zap.Error(err))
			// если broken pipe - переподключаемся
			// if broken pipe - reconnect
			if errors.Is(err, syscall.EPIPE) {
				logger.Error("Broken pipe, reconnecting...")
				time.Sleep(1 * time.Second)
			}
		}
		// отправляем сигнал о завершении
		// send done signal
		done <- true
	}
}

// getRuntimeMetrics запускает мониторинг метрик памяти
//
// getRuntimeMetrics launches memory metrics monitoring
func getRuntimeMetrics(m *storage.MetricStorage, pool time.Duration) {
	for {
		m.RuntimeMetricsMonitor()
		time.Sleep(pool)
	}
}

// getGopsUtilMetrics запускает мониторинг метрик gops
//
// getGopsUtilMetrics launches gops metrics monitoring
func getGopsUtilMetrics(m *storage.MetricStorage, pool time.Duration) {
	for {
		m.GopsUtilizationMonitor()
		time.Sleep(pool)
	}
}

// Run запускает агент
//
// Run launches the agent
func Run(logger *zap.Logger) {
	// читаем конфиг
	// read config
	conf := config.GetConfig()

	fields := []zapcore.Field{
		zap.Int("rate limit", conf.RateLimit),
		zap.String("poll interval", conf.PollInterval.String()),
		zap.String("report interval", conf.ReportInterval.String()),
		zap.String("server address", conf.ServerAddress),
	}
	// если ключ не пустой - добавляем сообщение в лог
	// if the key is not empty - add a message to the log
	if conf.Key != "" {
		msg := "key is presented"
		fields = append(fields, zap.String("key", msg))
	}
	logger.Info("Config loaded", fields...)
	// инициализируем хранилище
	// initialize storage
	memDB := storage.NewAgentStorage()
	// запускаем мониторинг метрик
	// start metrics monitoring
	var wg sync.WaitGroup
	wg.Add(2)
	go getRuntimeMetrics(memDB, conf.PollInterval)
	go getGopsUtilMetrics(memDB, conf.PollInterval)
	// запускаем отправку метрик на сервер в зависимости от наличия лимита
	// start sending metrics to the server depending on the limit
	if conf.RateLimit > 0 {
		for {
			// ждем интервал
			// wait for interval
			time.Sleep(conf.ReportInterval)
			// получаем метрики из хранилища
			// get metrics from storage
			data := memDB.JSONMetrics()
			// создаем каналы для воркеров
			// create channels for workers
			jobs := make(chan []byte, len(data))
			done := make(chan bool, len(data))
			// запускаем воркеров
			// start workers
			for w := 1; w <= conf.RateLimit; w++ {
				wg.Add(1)
				go SendOneMetric(logger, conf, jobs, done)
			}
			// отправляем метрики в канал
			// send metrics to channel
			for _, metric := range data {
				jobs <- metric
			}
			// закрываем каналы
			// close channels
			close(jobs)
			// ждем завершения воркеров
			// wait for workers to finish
			for range data {
				<-done
			}
			// закрываем воркеров
			// close workers
			wg.Done()
		}
	}
	// если лимит не установлен - отправляем метрики пачкой в json на сервер
	// if the limit is not set - send metrics in batches in json to the server
	if conf.RateLimit <= 0 {
		for {
			// ждем интервал
			// wait for interval
			time.Sleep(conf.ReportInterval)
			// получаем метрики из хранилища
			// get metrics from storage
			data := memDB.BatchJSONMetrics()
			// считаем хеш
			// calculate hash
			var hash [32]byte
			if conf.Key != "" {
				hash = hash2.GetHash(conf.Key, data)
			}
			// отправляем метрики
			// send metrics
			err := httpclient.SendBatchJSONMetrics(logger, conf, data, hash)
			if err != nil {
				logger.Error("Error sending metrics: ", zap.Error(err))
			}
		}
	}
	// ждем завершения воркеров
	// wait for workers to finish
	wg.Wait() //nolint:govet
}
