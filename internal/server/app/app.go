// Package app реализует приложение, в котором создается конфигурация приложения из yaml файла,
// обрабатываются флаги и переменные окружения если присутствуют,
// хранение метрик происходит в памяти, в файле или в базе данных в зависимости от конфигурации
// в процессе запуска создается логгер и стартует http сервер с выбранным хранилищем
//
// package app implements an application in which the application configuration is created from a yaml file,
// flags and environment variables are processed if present,
// metrics are stored in memory, in a file or in a database depending on the configuration
// during startup, a logger is created and an http server with the selected storage starts
package app

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/filestorage"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/inmemorystorage"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/postgrestorage"
)

// DataBaser интерфейс для работы с хранилищем
// интерфейс описывает методы inmemory хранилища и postgreSQL хранилища
//
// DataBaser interface for working with storage
// the interface describes the methods of the inmemory storage and the postgreSQL storage
type DataBaser interface {
	SetCounter(key string, value int64)
	SetGauge(key string, value float64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
	Ping() error
}

// Run запускает приложение
//
// Run starts the application
func Run(sigint <-chan os.Signal, connectionsClosed chan<- struct{}) {
	// читаем конфигурацию
	// read configuration
	conf := config.GetConfig()
	//достаем логгер из структуры конфига

	logger := conf.Logger

	ctx := context.Background()
	// создаем хранилище в памяти
	// create inmemory storage
	memDB := inmemorystorage.NewMemStorage(logger)
	// присваиваем переменной db хранилище в памяти
	// assign the db variable to the inmemory storage
	var db DataBaser = memDB
	// инициализируем хранилище в файле
	// initialize file storage
	file := filestorage.NewFileDB(conf.File.Path, conf.File.StoreInterval, logger)
	// создаем хранилище в postgreSQL
	// create postgreSQL storage
	pgDB := postgrestorage.NewPostgresDB(conf.DB.Dsn, logger)
	defer pgDB.Close()
	// если в конфиге указано использовать postgreSQL - создаем таблицы
	// if the config specifies to use postgreSQL - create tables
	if conf.DB.UsePG {
		err := pgDB.Create()
		if err != nil {
			logger.Sugar().Errorf("Error creating DB: %s", err)
		}
		// присваиваем переменной db хранилище в postgreSQL
		// assign the db variable to the postgreSQL storage
		db = pgDB
	}
	// если в конфиге указано не использовать postgreSQL, а файл - восстанавливаем метрики из файла
	// if the config specifies not to use postgreSQL, but a file - restore metrics from file
	if !conf.DB.UsePG && conf.File.UseFile && conf.File.Restore {
		metrics, err := file.Read(ctx)
		if err != nil {
			logger.Sugar().Errorf("Error reading metrics from file: %s", err)

		}
		err = memDB.RestoreFromSerialized(metrics)
		if err != nil {
			logger.Sugar().Errorf("Error restoring metrics from file: %s", err)
		}
	}
	// если в конфиге указано не использовать postgreSQL, а файл
	//без восстановления сохраненных данных - запускаем запись метрик в файл
	// if the config specifies not to use postgreSQL, but a file
	// without restoring saved data - start writing metrics to file
	if !conf.DB.UsePG && conf.File.UseFile {
		t := time.NewTicker(conf.File.StoreInterval)
		defer t.Stop()

		go func() {
			for range t.C {
				metrics := memDB.GetAllSerialized()
				err := file.Write(ctx, metrics)
				if err != nil {
					logger.Sugar().Errorf("Error writing metrics to file: %s", err)
				}
			}
		}()
	}
	// набираем поля для логгера
	// collect fields for logger
	fields := []zapcore.Field{
		zap.String("address", conf.Address),
		zap.String("log_level", conf.LogLevel),
	}
	if conf.DB.UsePG {
		fields = append(fields, zap.String("postgreSQL database_dsn", conf.DB.Dsn))
	} else if conf.File.UseFile {
		fields = append(fields, zap.String("file path", conf.File.Path))
		fields = append(fields, zap.String("file store interval", conf.File.StoreInterval.String()))
		fields = append(fields, zap.Bool("restore from file", conf.File.Restore))
	}
	logger.Info("Started server", fields...)
	// создаем http сервер
	// create http server
	srv := &http.Server{
		Addr:    conf.Address,
		Handler: httpserver.MetricRouter(logger, db, conf),
	}
	// запускаем http сервер
	// start http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen", zap.Error(err))
		}
	}()
	// ожидаем сигнал о завершении
	// wait for done signal
	for signal := range sigint {
		logger.Info("Received signal", zap.String("signal", signal.String()))
		logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("server shutdown error", zap.Error(err))
		}
		// если в конфиге указано использовать postgreSQL - закрываем соединение
		// if the config specifies to use postgreSQL - close the connection
		if conf.DB.UsePG {
			pgDB.Close()
		}
		cancel()
		logger.Info("Server shutdown gracefully")
		close(connectionsClosed)
		return
	}

}
