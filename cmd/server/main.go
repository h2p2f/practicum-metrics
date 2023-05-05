package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/database"
	"github.com/h2p2f/practicum-metrics/internal/server/handlers"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
)

// MetricRouter function to create chi router
func MetricRouter(m *storage.MemStorage, db *database.PGDB) chi.Router {
	//get handlers
	handler := handlers.NewMetricHandler(m, db)
	//create router
	r := chi.NewRouter()
	//add middlewares
	loggedAndZippedRouter := r.With(logger.WithLogging, handlers.GzipHanle)
	loggedRouter := r.With(logger.WithLogging)
	//add routes
	loggedAndZippedRouter.Post("/update/", handler.UpdateJSON)
	loggedAndZippedRouter.Post("/value/", handler.ValueJSON)
	loggedRouter.Post("/update/{metric}/{key}/{value}", handler.UpdatePage)
	loggedRouter.Get("/value/{metric}/{key}", handler.GetMetricValue)
	loggedRouter.Get("/ping", handler.DBPing)
	loggedAndZippedRouter.Get("/", handler.MainPage)
	//
	return r
}

func main() {

	//init logger
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	//setup new config
	conf := config.NewConfig()
	//set config from flags and env
	conf.SetConfig(getFlagsAndEnv())

	//create storage
	m := storage.NewMemStorage()
	pgDB := database.NewPostgresDB(conf.Database)
	defer pgDB.Close()

	//create fileDB with path and interval from config
	fileDB := storage.NewFileDB(conf.PathToStoreFile, conf.StoreInterval)
	if conf.UseDB {
		conf.UseFile = false
	}
	//restore metrics from file if flag is set
	if conf.UseFile && conf.Restore {
		metrics, err := fileDB.ReadFromFile()
		if err != nil {
			fmt.Println(err)
		}
		m.RestoreMetrics(metrics)
	}
	if conf.UseDB {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := pgDB.CreateTable(ctx)
		if err != nil {
			logger.Log.Sugar().Errorf("Error creating DB table: %s", err)
		}
		if conf.Restore {
			metrics, err := pgDB.ReadFromDB(ctx)
			if err != nil {
				logger.Log.Sugar().Errorf("Error reading metrics from DB: %s", err)
			}
			m.RestoreMetric(metrics)
		}
		go func() {
			for {
				time.Sleep(conf.StoreInterval * time.Second)
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				met := m.GetAllInBytesSliced()
				//fmt.Println(met)
				err := pgDB.SaveToDB(ctx, met)
				if err != nil {
					fmt.Println(err)
				}
				cancel()
			}
		}()
	}

	//save metrics to file with interval from config
	//made with anonymous function and goroutine
	if conf.UseFile {
		go func() {
			for {
				time.Sleep(conf.StoreInterval * time.Second)
				met := m.GetAllMetricsSliced()
				fmt.Println(met)
				err := fileDB.SaveToFile(met)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}

	//start server with router
	logger.Log.Sugar().Infof("Server started on %s", conf.ServerAddress)
	if conf.UseDB {
		logger.Log.Sugar().Infof("with DB %s", conf.Database)
		logger.Log.Sugar().Infof("with param: store interval %s", conf.StoreInterval)
		logger.Log.Sugar().Infof("restore from DB %t", conf.Restore)
	}
	if conf.UseFile {
		logger.Log.Sugar().Infof("with file %s", conf.PathToStoreFile)
		logger.Log.Sugar().Infof("with param: store interval %s", conf.StoreInterval)
		logger.Log.Sugar().Infof("restore from file %t", conf.Restore)
	}
	//logger.Log.Sugar().Infof("startup params", conf.UseDB, conf.UseFile, conf.Restore, conf.StoreInterval, conf.PathToStoreFile, conf.Database)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, MetricRouter(m, pgDB)))

}
