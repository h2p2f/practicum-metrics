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

var pgDB *database.PGDB

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

	pgDB = database.NewPostgresDB(conf.Database)
	defer pgDB.Close()

	//create fileDB with path and interval from config
	fileDB := storage.NewFileDB(conf.PathToStoreFile, conf.StoreInterval)
	if conf.UseDB {
		conf.UseFile = false
	}

	//restore metrics from file if flag is set
	if conf.UseFile && conf.Restore {
		logger.Log.Sugar().Info("trying to use file")
		logger.Log.Sugar().Infof("with file %s", conf.PathToStoreFile)
		logger.Log.Sugar().Infof("with param: store interval %s", conf.StoreInterval)
		logger.Log.Sugar().Infof("need restore from file %t", conf.Restore)
		metrics, err := fileDB.ReadFromFile()
		if err != nil {
			fmt.Println(err)
		}
		m.RestoreMetrics(metrics)
		//go func() {
		//	for {
		//		time.Sleep(conf.StoreInterval * time.Second)
		//		met := m.GetAllMetricsSliced()
		//		fmt.Println(met)
		//		err := fileDB.SaveToFile(met)
		//		if err != nil {
		//			fmt.Println(err)
		//		}
		//	}
		//}()
	}

	if conf.UseDB {
		logger.Log.Sugar().Info("trying to use DB")
		logger.Log.Sugar().Infof("with DB %s", conf.Database)
		logger.Log.Sugar().Infof("with param: store interval %s", conf.StoreInterval)
		logger.Log.Sugar().Infof("need restore from DB %t", conf.Restore)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
	}
	//if conf.UseDB {
	//	go func() {
	//		for {
	//			time.Sleep(conf.StoreInterval * time.Second)
	//			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	//			met := m.GetAllInBytesSliced()
	//			//fmt.Println(met)
	//			err := pgDB.SaveToDB(ctx, met)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//			cancel()
	//		}
	//	}()
	//}

	t := time.NewTicker(conf.StoreInterval)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				fmt.Println("tick")
				if conf.UseFile {
					fmt.Println("use file")
					met := m.GetAllMetricsSliced()
					err := fileDB.SaveToFile(met)
					if err != nil {
						fmt.Println(err)
					}
				}
				if conf.UseDB {
					fmt.Println("use DB")
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					met := m.GetAllInBytesSliced()
					err := pgDB.SaveToDB(ctx, met)
					if err != nil {
						fmt.Println(err)
					}
					cancel()
				}
			}
		}
	}()

	//start server with router
	logger.Log.Sugar().Infof("Server started on %s", conf.ServerAddress)

	//logger.Log.Sugar().Infof("startup params - useDB %t useFile %t, Restore %t", conf.UseDB, conf.UseFile, conf.Restore)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, MetricRouter(m, pgDB)))

}
