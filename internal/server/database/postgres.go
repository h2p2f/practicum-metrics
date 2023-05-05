package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type PGDB struct {
	db *sql.DB
}

// NewPostgresDB is a function that returns a new PostgresDB
func NewPostgresDB(param string) *PGDB {
	db, err := sql.Open("pgx", param)
	if err != nil {
		fmt.Println("Error opening database connection: ", err)
		log.Fatal(err)
	}
	return &PGDB{db: db}
}

// Close is a function that closes db
func (pgdb *PGDB) Close() {
	err := pgdb.db.Close()
	if err != nil {
		log.Println(err)
	}
}

func (pgdb *PGDB) PingContext(ctx context.Context) error {
	err := pgdb.db.PingContext(ctx)
	return err
}

func (pgdb *PGDB) CreateTable(ctx context.Context) (err error) {
	query := `CREATE TABLE IF NOT EXISTS metrics (
		    id text not null,
		    mtype text not null,
		    delta bigint,
		    value double precision
		    );`
	_, err = pgdb.db.ExecContext(ctx, query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (pgdb *PGDB) InsertMetric(ctx context.Context, id string, mtype string, delta *int64, value *float64) (err error) {
	query := `INSERT INTO metrics (id, mtype, delta, value) VALUES ($1, $2, $3, $4);`
	_, err = pgdb.db.ExecContext(ctx, query, id, mtype, delta, value)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (pgdb *PGDB) UpdateMetric(ctx context.Context, id string, mtype string, delta *int64, value *float64) (err error) {
	query := `UPDATE metrics SET delta = $1, value = $2 WHERE id = $3 AND mtype = $4;`
	_, err = pgdb.db.ExecContext(ctx, query, delta, value, id, mtype)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (pgdb *PGDB) ReadFromDB(ctx context.Context) ([][]byte, error) {
	var result [][]byte
	rows, err := pgdb.db.QueryContext(ctx, "SELECT * FROM metrics;")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if rows.Err() != nil {
		log.Println(err)
		return nil, err
	}
	var met metrics
	for rows.Next() {
		err = rows.Scan(&met.ID, &met.MType, &met.Delta, &met.Value)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		metJSON, err := json.Marshal(met)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, metJSON)
	}
	return result, nil
}

func (pgdb *PGDB) SaveToDB(ctx context.Context, met [][]byte) error {
	truncQuery := `TRUNCATE TABLE metrics;`
	_, err := pgdb.db.ExecContext(ctx, truncQuery)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, metric := range met {
		var met metrics
		err = json.Unmarshal(metric, &met)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = pgdb.db.ExecContext(ctx,
			"INSERT INTO metrics (id, mtype, delta, value) VALUES ($1, $2, $3, $4);",
			met.ID, met.MType, met.Delta, met.Value)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
