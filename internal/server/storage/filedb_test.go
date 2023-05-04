package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestFileDB_ReadFromFile(t *testing.T) {
	type fields struct {
		File     *os.File
		FilePath string
		Interval time.Duration
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []string
	}{
		{
			name: "positive test1 (counters only)",
			fields: fields{
				FilePath: "/tmp/test1.json",
				Interval: 30,
			},
			wantData: []string{"{\"id\":\"TestGet100\",\"type\":\"counter\",\"delta\":13065}",
				"{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":30097}"},
		},
		{
			name: "positive test2 (gauges only)",
			fields: fields{
				FilePath: "/tmp/test2.json",
				Interval: 30,
			},
			wantData: []string{"{\"id\":\"TestGet100\",\"type\":\"gauge\",\"value\":13065}",
				"{\"id\":\"PollCount\",\"type\":\"gauge\",\"value\":30097}"},
		},
		{
			name: "positive test3 (mixed metrics))",
			fields: fields{
				FilePath: "/tmp/test3.json",
				Interval: 30,
			},
			wantData: []string{"{\"id\":\"TestGet100\",\"type\":\"counter\",\"delta\":13065}",
				"{\"id\":\"PollCount\",\"type\":\"gauge\",\"value\":30097}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				File:     tt.fields.File,
				FilePath: tt.fields.FilePath,
				Interval: tt.fields.Interval,
			}
			var err error
			f.File, err = os.OpenFile(tt.fields.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			var want []Metrics
			for _, line := range tt.wantData {
				data := []byte(line)
				metric := Metrics{}
				err = json.Unmarshal(data, &metric)
				want = append(want, metric)
			}

			if err != nil {
				log.Fatalf("error while opening file: %s", err)
			}
			for _, metric := range tt.wantData {
				_, err = f.File.Write([]byte(metric + "\n"))
				if err != nil {
					log.Fatalf("error while writing to file: %s", err)
				}
			}
			if err := f.File.Close(); err != nil {
				panic(err)
			}

			gotMetrics, err := f.ReadFromFile()
			if err != nil {
				t.Errorf("ReadFromFile() error = %v", err)
			}

			assert.Equalf(t, want, gotMetrics, "ReadFromFile()")
		})
	}
}

func TestFileDB_SaveToFile(t *testing.T) {
	type fields struct {
		File     *os.File
		FilePath string
		Interval time.Duration
	}
	type args struct {
		metrics []Metrics
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []string
	}{
		{
			name: "positive test 1 (counters only)",
			fields: fields{
				FilePath: "/tmp/test1.json",
				Interval: 30,
			},
			wantData: []string{"{\"id\":\"TestGet100\",\"type\":\"counter\",\"delta\":13065}",
				"{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":30097}"},
		},
		{
			name: "positive test 2 (gauges only)",
			fields: fields{
				FilePath: "/tmp/test2.json",
				Interval: 30,
			},
			wantData: []string{"{\"id\":\"TestGet100\",\"type\":\"gauge\",\"value\":13065}",
				"{\"id\":\"PollCount\",\"type\":\"gauge\",\"value\":30097}"},
		},
		{
			name: "positive test 3 (mixed metrics)",
			fields: fields{
				FilePath: "/tmp/test3.json",
				Interval: 30,
			},
			wantData: []string{"{\"id\":\"TestGet100\",\"type\":\"counter\",\"delta\":13065}",
				"{\"id\":\"PollCount\",\"type\":\"gauge\",\"value\":30097}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileDB{
				File:     tt.fields.File,
				FilePath: tt.fields.FilePath,
				Interval: tt.fields.Interval,
			}
			var err error
			var want []Metrics
			for _, line := range tt.wantData {
				data := []byte(line)
				metric := Metrics{}
				err = json.Unmarshal(data, &metric)
				want = append(want, metric)
			}
			err = f.SaveToFile(want)
			if err != nil {
				log.Fatalf("error while opening file: %s", err)
			}
			gotMetrics, err := f.ReadFromFile()
			assert.Equalf(t, want, gotMetrics, "SaveToFile()")
		})
	}
}
