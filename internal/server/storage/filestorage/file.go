// Package filestorage описывает хранилище метрик в файле.
//
// package filestorage describes a metric store in a file.
package filestorage

import (
	"bufio"
	"context"
	"os"
	"time"

	"go.uber.org/zap"
)

// FileDB - структура, описывающая хранилище метрик в файле.
//
// FileDB - a structure that describes a metric store in a file.
type FileDB struct {
	logger   *zap.Logger
	File     *os.File
	FilePath string
	Interval time.Duration
	//mut      sync.RWMutex

}

// NewFileDB - конструктор для FileDB.
//
// NewFileDB is a function that returns a new fileDB
func NewFileDB(filePath string, interval time.Duration, logger *zap.Logger) *FileDB {
	return &FileDB{
		FilePath: filePath,
		Interval: interval,
		logger:   logger,
	}
}

// Write - функция, записывающая метрики в файл.
//
// Write is a function that writes metrics to file
func (f *FileDB) Write(ctx context.Context, metrics [][]byte) error {
	var err error
	f.File, err = os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer func() {
		if err2 := f.File.Close(); err2 != nil {
			f.logger.Sugar().Fatalf("error while closing file on write: %s", err2)
		}
	}()
	if err != nil {
		return err
	}
	for _, metric := range metrics {
		_, err = f.File.Write(append(metric, '\n'))
		if err != nil {
			return err
		}
	}
	f.logger.Sugar().Infof("saved to file - success")
	return nil
}

// Read - функция, считывающая метрики из файла.
//
// Read is a function that reads metrics from file
func (f *FileDB) Read(ctx context.Context) ([][]byte, error) {
	var err error
	_, err = os.Stat(f.FilePath)
	if os.IsNotExist(err) {
		return nil, err
	}
	f.File, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0755)
	defer func() {
		if err2 := f.File.Close(); err2 != nil {
			f.logger.Sugar().Errorf("error while closing file on read: %s", err2)
		}
	}()
	if err != nil {
		return nil, err
	}
	var metrics [][]byte
	scan := bufio.NewScanner(f.File)
	for {
		if !scan.Scan() {
			break
		}
		metrics = append(metrics, scan.Bytes())
	}
	f.logger.Sugar().Infof("read from file - success")
	return metrics, nil
}
