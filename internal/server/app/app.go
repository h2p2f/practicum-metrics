package app

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/filestorage"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/inmemorystorage"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/postgrestorage"
)

type DataBaser interface {
	SetCounter(key string, value int64)
	SetGauge(key string, value float64)
	GetCounter(name string) (value int64, err error)
	GetGauge(name string) (value float64, err error)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
	Ping() error
}

func Run() {

	conf := config.GetConfig()

	if err := logger.InitLogger(conf.LogLevel); err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	memDB := inmemorystorage.NewMemStorage(logger.Log)

	var db DataBaser = memDB

	file := filestorage.NewFileDB(conf.File.Path, conf.File.StoreInterval, logger.Log)

	pgDB := postgrestorage.NewPostgresDB(conf.DB.Dsn, logger.Log)
	defer pgDB.Close()

	if conf.DB.UsePG {
		err := pgDB.Create()
		if err != nil {
			logger.Log.Sugar().Errorf("Error creating DB: %s", err)
		}
		db = pgDB
	}

	if !conf.DB.UsePG && conf.File.UseFile && conf.File.Restore {
		metrics, err := file.Read(ctx)
		if err != nil {
			logger.Log.Sugar().Errorf("Error reading metrics from file: %s", err)

		}
		err = memDB.RestoreFromSerialized(metrics)
		if err != nil {
			logger.Log.Sugar().Errorf("Error restoring metrics from file: %s", err)
		}
	}

	if !conf.DB.UsePG && conf.File.UseFile {
		t := time.NewTicker(conf.File.StoreInterval)
		defer t.Stop()

		go func() {
			for range t.C {
				metrics := memDB.GetAllSerialized()
				err := file.Write(ctx, metrics)
				if err != nil {
					logger.Log.Sugar().Errorf("Error writing metrics to file: %s", err)
				}
			}
		}()
	}
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
	logger.Log.Info("Started server", fields...)

	err := http.ListenAndServe(conf.Address, httpserver.MetricRouter(logger.Log, db, conf.Key))
	if err != nil {
		panic(err)
	}

}
