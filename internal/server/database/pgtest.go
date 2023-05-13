package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
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
func TestPGDB_WriteAndReadFromDB(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	tests := []struct {
		db_addr string
		name    string
		fields  []dbmetrics
		wantErr bool
	}{
		{
			db_addr: "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			name:    "positive test1",
			fields: []dbmetrics{
				{
					ID:    "test1",
					MType: "counter",
					Delta: 10,
				}},
			wantErr: false,
		},
		{
			db_addr: "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			name:    "positive test2",
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
			db_addr: "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
			name:    "negative test1",
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
		t.Run(tt.name, func(t *testing.T) {
			db, err := sql.Open("pgx", tt.db_addr)
			if err != nil {
				t.Errorf("can't open db: %v", err)

			}
			defer db.Close()

			pg := &PGDB{
				db: db,
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
