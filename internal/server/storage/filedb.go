package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

//TODO: put in order this code

// FileDB is a struct that contains file path and interval to store metrics, mutex, file
type FileDB struct {
	File     *os.File
	FilePath string
	Interval time.Duration
	mut      sync.RWMutex
}

// Metrics is a struct that contains all the metrics that are being stored in file
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewMetricsCounter is a constructor function that returns a new counter struct
func NewMetricsCounter(ID, MType string, delta int64) *Metrics {
	return &Metrics{
		ID:    ID,
		MType: MType,
		Delta: &delta,
	}
}

// NewMetricsGauge is a constructor function that returns a new gauge struct
func NewMetricsGauge(ID, MType string, value float64) *Metrics {
	return &Metrics{
		ID:    ID,
		MType: MType,
		Value: &value,
	}
}

// NewFileDB is a function that returns a new fileDB
func NewFileDB(filePath string, interval time.Duration) *FileDB {
	return &FileDB{
		FilePath: filePath,
		Interval: interval,
	}
}

// SaveToFile is a function that saves metrics to file
func (f *FileDB) SaveToFile(metrics []Metrics) (err error) {
	//fmt.Println("saving to file", metrics)
	f.File, err = os.OpenFile(f.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			log.Fatalf("error while closing file: %s", err)
		}
	}()
	//defer f.File.Close()
	if err != nil {
		return err
	}
	for _, metric := range metrics {
		writeData, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		_, err = f.File.Write(append(writeData, '\n'))
		if err != nil {
			return err
		}
	}
	fmt.Println("saved to file - success")
	return nil
}

// ReadFromFile is a function that reads metrics from file
func (f *FileDB) ReadFromFile() (metrics []Metrics, err error) {
	_, err = os.Stat(f.FilePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	f.File, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			panic(err)
		}
	}()
	scan := bufio.NewScanner(f.File)
	for {
		if !scan.Scan() {
			break
		}
		metric := Metrics{}
		fmt.Println("loaded data from file: ", scan.Text())
		data := scan.Bytes()
		err = json.Unmarshal(data, &metric)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
		//fmt.Println("read from file: ", metric)

	}
	return metrics, nil
}
