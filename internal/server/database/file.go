package database

import (
	"bufio"
	"context"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

// FileDB is a struct that contains file path and interval to store metrics, mutex, file
type FileDB struct {
	File     *os.File
	FilePath string
	Interval time.Duration
	mut      sync.RWMutex
	logger   *zap.Logger
}

// NewFileDB is a function that returns a new fileDB
func NewFileDB(filePath string, interval time.Duration, logger *zap.Logger) *FileDB {
	return &FileDB{
		FilePath: filePath,
		Interval: interval,
		logger:   logger,
	}
}

// Create it is a stub function to implement the interface
func (f *FileDB) Create(ctx context.Context) error {
	return nil
}

// Write is a function that writes metrics to file
func (f *FileDB) Write(ctx context.Context, metrics [][]byte) error {
	var err error
	f.File, err = os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			f.logger.Sugar().Fatalf("error while closing file on write: %s", err)
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

// Read is a function that reads metrics from file
func (f *FileDB) Read(ctx context.Context) ([][]byte, error) {
	var err error
	_, err = os.Stat(f.FilePath)
	if os.IsNotExist(err) {
		return nil, err
	}
	f.File, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			f.logger.Sugar().Errorf("error while closing file on read: %s", err)
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
