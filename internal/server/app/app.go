// Package app implements an application in which the application configuration is created from a yaml file,
// flags and environment variables are processed if present,
// metrics are stored in memory, in a file or in a database depending on the configuration
// during startup, a logger is created and an http server with the selected storage starts
package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/grpcserver"
	"github.com/h2p2f/practicum-metrics/internal/server/grpcserver/middlewares"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/filestorage"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/inmemorystorage"
	"github.com/h2p2f/practicum-metrics/internal/server/storage/postgrestorage"
	pb "github.com/h2p2f/practicum-metrics/proto"
)

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

// Run starts the application
func Run(sigint chan os.Signal, connectionsClosed chan<- struct{}) {
	// read configuration
	conf, logger, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create inmemory storage
	memDB := inmemorystorage.NewMemStorage(logger)
	// assign the db variable to the inmemory storage
	var db DataBaser = memDB
	// initialize file storage
	file := filestorage.NewFileDB(conf.File.Path, conf.File.StoreInterval, logger)
	// create postgreSQL storage
	pgDB := postgrestorage.NewPostgresDB(conf.DB.Dsn, logger)
	defer pgDB.Close()
	// if the config specifies to use postgreSQL - create tables
	if conf.DB.UsePG {
		err := pgDB.Create()
		if err != nil {
			logger.Sugar().Errorf("Error creating DB: %s", err)
		}
		// assign the db variable to the postgreSQL storage
		db = pgDB
	}
	// if the config specifies not to use postgreSQL, but a file - restore metrics from file
	if !conf.DB.UsePG && conf.File.UseFile && conf.File.Restore {
		restoreFromFile(ctx, logger, file, memDB)
	}
	// if the config specifies not to use postgreSQL, but a file
	// without restoring saved data - start writing metrics to file
	if !conf.DB.UsePG && conf.File.UseFile {
		go saveToFile(ctx, conf.File.StoreInterval, file, logger, memDB)
	}
	// collect fields for logger
	fields := []zapcore.Field{
		zap.String("address", conf.HTTP.Address),
		zap.String("log_level", conf.LogLevel),
	}
	if conf.DB.UsePG {
		fields = append(fields, zap.String("postgreSQL", conf.DB.Dsn))
	} else if conf.File.UseFile {
		fields = append(fields, zap.String("file_path", conf.File.Path))
		fields = append(fields, zap.String("store_interval", conf.File.StoreInterval.String()))
		fields = append(fields, zap.Bool("restore_file", conf.File.Restore))
	}
	logger.Info("Started http server", fields...)
	// create http server
	srv := &http.Server{
		Addr:    conf.HTTP.Address,
		Handler: httpserver.MetricRouter(logger, db, conf),
	}
	// start http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen", zap.Error(err))
		}
	}()
	// create grpc server
	logger.Info("Started grpc server", zap.String("address", conf.GRPC.Address))
	listen, err := net.Listen("tcp", conf.GRPC.Address)
	if err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
	// create grpc server with middlewares
	var opts []grpc.ServerOption
	//opts = middlewares.WithLogging(logger, opts)
	//opts = middlewares.WithCheckingIP(conf.HTTP.TrustSubnet, opts)
	opts = middlewares.WithChekingIPAndLogging(logger, conf.HTTP.TrustSubnet, opts)
	grpcServer := grpc.NewServer(opts...)

	grpcMetrics := grpcserver.NewServer(db, logger)
	pb.RegisterMetricsServiceServer(grpcServer, grpcMetrics)
	// start grpc server
	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			logger.Fatal("listen", zap.Error(err))
		}
	}()

	// wait for done signal
	<-sigint
	logger.Info("Shutting down server...")
	ctx2, cancel2 := context.WithTimeout(ctx, 5*time.Second)
	if err := srv.Shutdown(ctx2); err != nil {
		logger.Fatal("server shutdown error", zap.Error(err))
	}
	if conf.DB.UsePG {
		pgDB.Close()
	}
	cancel2()
	close(sigint)
	logger.Info("Server shutdown gracefully")
	close(connectionsClosed)
	cancel()
}

// saveToFile - function for writing metrics to a file
func saveToFile(
	ctx context.Context,
	interval time.Duration,
	file *filestorage.FileDB,
	logger *zap.Logger,
	memDB *inmemorystorage.MemStorage) {

	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		metrics := memDB.GetAllSerialized()
		err := file.Write(ctx, metrics)
		if err != nil {
			logger.Error("could not write metrics to file", zap.Error(err))
		}
		if ctx.Err() != nil && errors.Is(context.Canceled, ctx.Err()) {
			return
		}
	}
}

// restoreFromFile - function for restoring metrics from a file
func restoreFromFile(
	ctx context.Context,
	logger *zap.Logger,
	file *filestorage.FileDB,
	memDB *inmemorystorage.MemStorage) {
	metrics, err := file.Read(ctx)
	if err != nil {
		logger.Error("could not read metrics from file", zap.Error(err))
	}
	err = memDB.RestoreFromSerialized(metrics)
	if err != nil {
		logger.Error("could not restore metrics from file", zap.Error(err))
	}
}
