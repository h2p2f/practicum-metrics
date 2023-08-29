// Package postgrestorage implements a metric store in PostgreSQL.
package postgrestorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type pg struct {
	db     *sql.DB
	logger *zap.Logger
}

// SetCounter sets the counter value by name.
func (pg *pg) SetCounter(key string, value int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mType := "counter"
	query := `INSERT INTO metrics (id, mtype, delta) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + excluded.delta;`
	_, err := pg.db.ExecContext(ctx, query, key, mType, value)
	if err != nil {
		pg.logger.Sugar().Errorf("Error inserting counter: %v", err)
	}
}

// SetGauge sets the gauge value by name.
func (pg *pg) SetGauge(key string, value float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mType := "gauge"
	query := `INSERT INTO metrics (id, mtype, value) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET value = $3;`
	_, err := pg.db.ExecContext(ctx, query, key, mType, value)
	if err != nil {
		pg.logger.Sugar().Errorf("Error inserting gauge: %v", err)
	}
}

// GetCounter returns the counter value by name.
func (pg *pg) GetCounter(name string) (value int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `SELECT delta FROM metrics WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, query, name)
	err = row.Scan(&value)
	if err != nil {
		pg.logger.Sugar().Errorf("Error scanning row: %v", err)
		return 0, nil
	}
	return value, nil
}

// GetGauge returns the gauge value by name.
func (pg *pg) GetGauge(name string) (value float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `SELECT value FROM metrics WHERE id = $1;`
	row := pg.db.QueryRowContext(ctx, query, name)
	err = row.Scan(&value)
	if err != nil {
		pg.logger.Sugar().Errorf("Error scanning row: %v", err)
		return 0, err
	}
	return value, nil
}

// GetCounters returns all counter values.
func (pg *pg) GetCounters() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mType := "counter"
	query := `SELECT id, delta FROM metrics WHERE mtype = $1;`
	rows, err := pg.db.QueryContext(ctx, query, mType)
	if rows.Err() != nil {
		pg.logger.Sugar().Errorf("Error reading from database: %v", err)
		return nil
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.logger.Sugar().Errorf("no row in result: %v", err)
			return nil
		} else {
			pg.logger.Sugar().Errorf("Error querying rows: %v", err)
			return nil
		}
	}
	defer func() {
		err2 := rows.Close()
		if err2 != nil {
			pg.logger.Sugar().Errorf("Error closing rows: %v", err2)
		}
	}()
	counters := make(map[string]int64)
	for rows.Next() {
		var key string
		var value int64
		err = rows.Scan(&key, &value)
		if err != nil {
			pg.logger.Sugar().Errorf("Error scanning row: %v", err)
			return nil
		}
		counters[key] = value
	}
	return counters
}

// GetGauges returns all gauge values.
func (pg *pg) GetGauges() map[string]float64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mType := "gauge"
	query := `SELECT id, value FROM metrics WHERE mtype = $1;`
	rows, err := pg.db.QueryContext(ctx, query, mType)
	if rows.Err() != nil {
		pg.logger.Sugar().Errorf("Error reading from database: %v", err)
		return nil
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.logger.Sugar().Errorf("no row in result: %v", err)
			return nil
		} else {
			pg.logger.Sugar().Errorf("Error querying rows: %v", err)
			return nil
		}
	}
	defer func() {
		err2 := rows.Close()
		if err2 != nil {
			pg.logger.Sugar().Errorf("Error closing rows: %v", err2)
		}
	}()
	gauges := make(map[string]float64)
	for rows.Next() {
		var key string
		var value float64
		err = rows.Scan(&key, &value)
		if err != nil {
			pg.logger.Sugar().Errorf("Error scanning row: %v", err)
			return nil
		}
		gauges[key] = value
	}
	return gauges
}

// NewPostgresDB creates a new instance of PostgresDB.
func NewPostgresDB(param string, logger *zap.Logger) *pg {

	db, err := sql.Open("pgx", param)
	if err != nil {
		// If the error is a connection exception, try to reconnect
		//this code wrote for increment #13
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			logger.Sugar().Errorf("Error opening database connection: %v, trying reconnect", err)
			//time delta for reconnect is 1, 3, 5 seconds
			waitTime := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
			//try to reconnect
			for _, v := range waitTime {
				time.Sleep(v)
				db, err = sql.Open("pgx", param)
				if err == nil {
					//if connection is successful, break the loop
					break
				}
			}
		}
	}
	return &pg{db: db, logger: logger}
}

// Close closes the database connection.
func (pg *pg) Close() {
	err := pg.db.Close()
	if err != nil {
		pg.logger.Sugar().Errorf("Error closing database connection: %v", err)
	}
}

// Create creates the metrics table.
func (pg *pg) Create() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `CREATE TABLE IF NOT EXISTS metrics (
		    id text not null PRIMARY KEY,
		    mtype text not null,
		    delta bigint,
		    value double precision
		    );`
	_, err = pg.db.ExecContext(ctx, query)
	if err != nil {
		pg.logger.Sugar().Errorf("Error creating table: %v", err)
		return err
	}
	pg.logger.Sugar().Info("Table metrics created successfully")
	return nil
}

// Ping checks the database connection.
func (pg *pg) Ping() error {
	ctx := context.Background()
	err := pg.db.PingContext(ctx)
	pg.logger.Sugar().Info("PingContext successfully")
	return err
}
