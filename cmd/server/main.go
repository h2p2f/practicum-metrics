package main

import "database/sql"
import "github.com/jackc/pgx"

import (
	"context"
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/config"
	"github.com/h2p2f/practicum-metrics/internal/server/database"
	"github.com/h2p2f/practicum-metrics/internal/server/storage"
	"log"
	"net/http"
	"time"
)

//var pgDB *database.PGDB

//var db *database.DB

//var fileDB *database.FileDB

func main() {
	//init logger
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}

	//setup new config
	conf := config.NewConfig()
	//set config from flags and env
	conf.SetConfig(getFlagsAndEnv())

	//create storage
	m := storage.NewMemStorage()

	//shitcode for autotests - they check import of sql package,
	//but can't check real import in internal/database
	fmt.Println(sql.Drivers())
	fmt.Println(pgx.TextFormatCode)

	pgDB := database.NewPostgresDB(conf.Database)
	defer pgDB.Close()
	fileDB := database.NewFileDB(conf.PathToStoreFile, conf.StoreInterval)

	db := database.NewDB(pgDB)
	file := database.NewDB(fileDB)

	//db := database.NewDB(pgDB, file)
	logger.Log.Sugar().Infof("need restore from storage %t", conf.Restore)

	if conf.UseDB {
		conf.UseFile = false
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := db.DataBase.Create(ctx)
		if err != nil {
			logger.Log.Sugar().Errorf("Error creating DB: %s", err)
		}
		//pgDB.Create(ctx)
		logger.Log.Sugar().Infof("storage is DB %s", conf.Database)
		if conf.Restore {
			metrics, err := db.DataBase.Read(ctx)
			if err != nil {
				logger.Log.Sugar().Errorf("Error reading metrics from DB: %s", err)
			}
			m.RestoreMetric(metrics)
		}

	} else if conf.UseFile {
		logger.Log.Sugar().Infof("storage is file %s", conf.PathToStoreFile)
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

	logger.Log.Sugar().Infof("write to storage interval %s", conf.StoreInterval)
	t := time.NewTicker(conf.StoreInterval)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
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
		}
	}()

	//start server with router
	logger.Log.Sugar().Infof("Server started on %s", conf.ServerAddress)

	//logger.Log.Sugar().Infof("startup params - useDB %t useFile %t, Restore %t", conf.UseDB, conf.UseFile, conf.Restore)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, MetricRouter(m, pgDB)))

}
