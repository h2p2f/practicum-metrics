package database

import (
	"bufio"
	"context"
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}
}

//TODO: put in order this code

// FileDB is a struct that contains file path and interval to store metrics, mutex, file
type FileDB struct {
	File     *os.File
	FilePath string
	Interval time.Duration
	mut      sync.RWMutex
}

// NewFileDB is a function that returns a new fileDB
func NewFileDB(filePath string, interval time.Duration) *FileDB {
	return &FileDB{
		FilePath: filePath,
		Interval: interval,
	}
}
func (f *FileDB) Create(ctx context.Context) error {
	return nil
}

func (f *FileDB) Write(ctx context.Context, metrics [][]byte) error {
	var err error
	f.File, err = os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			logger.Log.Sugar().Fatalf("error while closing file on write: %s", err)
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
	logger.Log.Sugar().Infof("saved to file - success")
	return nil
}

func (f *FileDB) Read(ctx context.Context) ([][]byte, error) {
	var err error
	_, err = os.Stat(f.FilePath)
	if os.IsNotExist(err) {
		return nil, err
	}
	f.File, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			log.Fatalf("error while closing file on read: %s", err)
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
	logger.Log.Sugar().Infof("read from file - success")
	return metrics, nil
}
