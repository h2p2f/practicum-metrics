package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
	"time"
)

type dbmetrics struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}

func convertMetricsToBytes(metrics []dbmetrics) [][]byte {
	var data [][]byte
	for _, met := range metrics {
		buf, err := json.Marshal(met)
		if err != nil {
			log.Println(err)
		}
		data = append(data, buf)
	}
	return data
}

func convertBytesToMetrics(data []byte) dbmetrics {
	var metrics dbmetrics
	//
	err := json.Unmarshal(data, &metrics)
	if err != nil {
		fmt.Println(err)
	}
	return metrics
}
func TestPGDBWriteAndReadFromDB(t *testing.T) {
	//type fields struct {
	//	db *sql.DB
	//}
	tests := []struct {
		dbAddr  string
		name    string
		fields  []dbmetrics
		wantErr bool
	}{
		{
			dbAddr: "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			name:   "positive test1",
			fields: []dbmetrics{
				{
					ID:    "test1",
					MType: "counter",
					Delta: 10,
				}},
			wantErr: false,
		},
		{
			dbAddr: "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			name:   "positive test2",
			fields: []dbmetrics{
				{
					ID:    "test2",
					MType: "gauge",
					Value: 10.00000001,
				},
			},
			wantErr: false,
		},
		{
			dbAddr: "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			name:   "negative test1",
			fields: []dbmetrics{
				{
					ID:    "test3",
					MType: "gauge",
					Value: 10.00002,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if os.Getenv("GITHUB_JOB") == "metricstest" {
			t.Skip("Skipping DB tests")
		}
		t.Run(tt.name, func(t *testing.T) {
			db, err := sql.Open("pgx", tt.dbAddr)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer func() {
				if err := db.Close(); err != nil {
					t.Errorf("can't close db: %v", err)
				}
			}()
			logger := zap.NewExample()
			pg := &PGDB{
				db:     db,
				logger: logger,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err := pg.Write(ctx, convertMetricsToBytes(tt.fields)); err != nil {
				t.Errorf("can't save to db: %v", err)
			}

			res, err := pg.GetValueByID(ctx, convertMetricsToBytes(tt.fields)[0])
			if err != nil {
				t.Errorf("can't read from db: %v", err)
			}
			resdbmetrics := convertBytesToMetrics(res)
			assert.Equal(t, tt.fields[0].Value, resdbmetrics.Value)

		})

	}

}
