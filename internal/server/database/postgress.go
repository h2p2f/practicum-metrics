package database

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type PGDB struct {
	db *sql.DB
}

// NewPostgresDB is a function that returns a new PostgresDB
func NewPostgresDB(param string) *PGDB {
	db, err := sql.Open("pgx", param)
	if err != nil {
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
