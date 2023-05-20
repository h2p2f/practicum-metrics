package main

import (
	"context"
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/server/httpserver"
	"net/http"
	"time"

	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/database"
	"github.com/h2p2f/practicum-metrics/internal/server/model"
)

func main() {
	//init logger
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}

	//setup new config
	conf := config.NewConfig()
	//set config from flags and env
	conf.SetConfig(config.GetFlagsAndEnv())

	//create model
	m := model.NewMemStorage()
	//create database and file model
	pgDB := database.NewPostgresDB(conf.Database, logger.Log)
	defer pgDB.Close()
	fileDB := database.NewFileDB(conf.PathToStoreFile, conf.StoreInterval, logger.Log)

	db := database.NewDB(pgDB)
	file := database.NewDB(fileDB)

	//db := database.NewDB(pgDB, file)
	logger.Log.Sugar().Infof("need restore from model %t", conf.Restore)
	//Create DB if not exist, restore metrics if it needs
	if conf.UseDB {
		conf.UseFile = false
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := db.DataBase.Create(ctx)
		if err != nil {
			logger.Log.Sugar().Errorf("Error creating DB: %s", err)
		}
		logger.Log.Sugar().Infof("model is DB %s", conf.Database)
		if conf.Restore {
			metrics, err := db.DataBase.Read(ctx)
			if err != nil {
				logger.Log.Sugar().Errorf("Error reading metrics from DB: %s", err)
			}
			m.RestoreMetric(metrics)
		}
		//if it needs use and restore from file
	} else if conf.UseFile {
		logger.Log.Sugar().Infof("model is file %s", conf.PathToStoreFile)
		if conf.Restore {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			//metrics, err := db.File.ReadFromFile(ctx)
			metrics, err := file.DataBase.Read(ctx)
			if err != nil {
				logger.Log.Sugar().Errorf("Error reading metrics from file: %s", err)
			}
			m.RestoreMetric(metrics)
		}
	}
	//periodically write to model
	logger.Log.Sugar().Infof("write to model interval %s", conf.StoreInterval)
	t := time.NewTicker(conf.StoreInterval)
	defer t.Stop()
	go func() {
		for range t.C {
			logger.Log.Sugar().Info("try to save data")
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			met := m.GetAllInBytesSliced()
			if conf.UseDB {
				err := db.DataBase.Write(ctx, met)
				if err != nil {
					logger.Log.Sugar().Errorf("Error writing metrics to DB: %s", err)
				}
			}
			if conf.UseFile {
				err := file.DataBase.Write(ctx, met)
				if err != nil {
					logger.Log.Sugar().Errorf("Error writing metrics to file: %s", err)
				}
			}
			cancel()
		}
	}()

	//start server with router
	logger.Log.Sugar().Infof("Server started on %s", conf.ServerAddress)
	logger.Log.Sugar().Fatalf("Server stopped with error: %s", http.ListenAndServe(conf.ServerAddress, httpserver.MetricRouter(m, pgDB, conf.Key)))

}
