package app

import (
	"context"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/database"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver"
	"github.com/h2p2f/practicum-metrics/internal/server/model"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Run - function to run server
func Run(logger *zap.Logger) {
	//init config
	conf := config.NewConfig()
	conf.SetConfig()

	//create model
	m := model.NewMemStorage()
	//create database and file model
	pgDB := database.NewPostgresDB(conf.Database, logger)
	defer pgDB.Close()
	fileDB := database.NewFileDB(conf.PathToStoreFile, conf.StoreInterval, logger)
	//create db and file models
	db := database.NewDB(pgDB)
	file := database.NewDB(fileDB)

	//db := database.NewDB(pgDB, file)
	logger.Sugar().Infof("need restore from model %t", conf.Restore)
	//Create DB if not exist, restore metrics if it needs
	if conf.UseDB {
		conf.UseFile = false
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := db.DataBase.Create(ctx)
		if err != nil {
			logger.Sugar().Errorf("Error creating DB: %s", err)
		}
		logger.Sugar().Infof("storage is DB %s", conf.Database)
		//if restore need - restore from DB
		if conf.Restore {
			metrics, err := db.DataBase.Read(ctx)
			if err != nil {
				logger.Sugar().Errorf("Error reading metrics from DB: %s", err)
			}
			m.RestoreMetric(metrics)
		}
		//if it needs use and restore from file
	} else if conf.UseFile {
		logger.Sugar().Infof("storage is file %s", conf.PathToStoreFile)
		//if restore need - restore from file
		if conf.Restore {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			//metrics, err := db.File.ReadFromFile(ctx)
			metrics, err := file.DataBase.Read(ctx)
			if err != nil {
				logger.Sugar().Errorf("Error reading metrics from file: %s", err)
			}
			m.RestoreMetric(metrics)
		}
	}
	//periodically write to model
	logger.Sugar().Infof("write to model interval %s", conf.StoreInterval)
	t := time.NewTicker(conf.StoreInterval)
	defer t.Stop()
	go func() {
		for range t.C {
			logger.Sugar().Info("try to save data")
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			met := m.GetAllInBytesSliced()
			if conf.UseDB {
				err := db.DataBase.Write(ctx, met)
				if err != nil {
					logger.Sugar().Errorf("Error writing metrics to DB: %s", err)
				}
			}
			if conf.UseFile {
				err := file.DataBase.Write(ctx, met)
				if err != nil {
					logger.Sugar().Errorf("Error writing metrics to file: %s", err)
				}
			}
			cancel()
		}
	}()

	//start server with router
	logger.Sugar().Infof("Server started on %s", conf.ServerAddress)
	logger.Sugar().Fatalf("Server stopped with error: %s",
		http.ListenAndServe(conf.ServerAddress, httpserver.MetricRouter(m, pgDB, conf.Key)))

}
