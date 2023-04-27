package storage

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type FileDB struct {
	File     *os.File
	FilePath string
	Interval time.Duration
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewFileDB(filePath string, interval time.Duration) *FileDB {
	return &FileDB{
		FilePath: filePath,
		Interval: interval,
	}
}
func (f *FileDB) SaveToFile(metrics []Metrics) (err error) {

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
		_, err = f.File.Write(writeData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileDB) ReadFromFile() (metrics []Metrics, err error) {
	f.File, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0755)
	defer func() {
		if err := f.File.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return nil, err
	}
	var metric Metrics
	for {
		err = json.NewDecoder(f.File).Decode(&metric)
		if err != nil {
			break
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
