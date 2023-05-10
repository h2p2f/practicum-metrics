package database

import (
	"database/sql"
	"github.com/h2p2f/practicum-metrics/internal/logger"
)

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

func init() {
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}
}

type PGDB struct {
	db *sql.DB
}

// NewPostgresDB is a function that returns a new PostgresDB
func NewPostgresDB(param string) *PGDB {
	db, err := sql.Open("pgx", param)
	if err != nil {
		logger.Log.Sugar().Errorf("Error opening database connection: %v", err)
	}
	return &PGDB{db: db}
}

// Close is a function that closes db
func (pgdb *PGDB) Close() {
	err := pgdb.db.Close()
	if err != nil {
		logger.Log.Sugar().Errorf("Error closing database connection: %v", err)
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
	logger.Log.Sugar().Info("Table metrics created successfully")
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

func (pgdb *PGDB) GetAllID(ctx context.Context) (ids []string, err error) {
	query := `SELECT id FROM metrics;`
	rows, err := pgdb.db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (pgdb *PGDB) ReadFromDB(ctx context.Context) ([][]byte, error) {
	var result [][]byte
	rows, err := pgdb.db.QueryContext(ctx, "SELECT * FROM metrics;")
	if err != nil {
		logger.Log.Sugar().Errorf("Error reading from database: %v", err)
		return nil, err
	}
	if rows.Err() != nil {
		logger.Log.Sugar().Errorf("Error reading from database: %v", err)
		//log.Println(err)
		return nil, err
	}
	var met metrics
	for rows.Next() {
		err = rows.Scan(&met.ID, &met.MType, &met.Delta, &met.Value)
		if err != nil {
			logger.Log.Sugar().Errorf("Error scan data rows from database: %v", err)
			return nil, err
		}
		metJSON, err := json.Marshal(met)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		result = append(result, metJSON)
	}
	logger.Log.Sugar().Info("Read from DB successfully")
	return result, nil
}

func (pgdb *PGDB) SaveToDB(ctx context.Context, met [][]byte) error {
	logger.Log.Sugar().Info("Saving to DB...")
	truncQuery := `TRUNCATE TABLE metrics;`
	_, err := pgdb.db.ExecContext(ctx, truncQuery)
	if err != nil {
		logger.Log.Sugar().Errorf("Error truncating table: %v", err)
		return err
	}

	for _, metric := range met {
		var met metrics
		err = json.Unmarshal(metric, &met)
		if err != nil {
			logger.Log.Sugar().Errorf("Error unmarshal data: %v", err)
			return err
		}
		tx, err := pgdb.db.BeginTx(ctx, nil)
		if err != nil {
			logger.Log.Sugar().Errorf("Error begin transaction: %v", err)
		}
		_, err = tx.ExecContext(ctx,
			"INSERT INTO metrics (id, mtype, delta, value) VALUES ($1, $2, $3, $4);",
			met.ID, met.MType, met.Delta, met.Value)
		if err != nil {
			logger.Log.Sugar().Errorf("Error inserting data: %v, rollback transaction", err)
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			logger.Log.Sugar().Errorf("Error commit transaction: %v", err)
		}
	}
	logger.Log.Sugar().Info("Save to DB successfully")
	return nil
}

func (pgdb *PGDB) SaveToDBWithoutTruncate(ctx context.Context, met [][]byte) error {
	logger.Log.Sugar().Info("Saving to DB without truncate...")
	for _, metric := range met {
		var met metrics
		err := json.Unmarshal(metric, &met)
		if err != nil {
			logger.Log.Sugar().Errorf("Error unmarshal data: %v", err)
			return err
		}
		ids, err := pgdb.GetAllID(ctx)
		tx, err := pgdb.db.BeginTx(ctx, nil)

		if contains(ids, met.ID) {

			err = pgdb.UpdateMetric(ctx, met.ID, met.MType, met.Delta, met.Value)
			if err != nil {
				logger.Log.Sugar().Errorf("Error updating data: %v, rollback transaction", err)
				err2 := tx.Rollback()
				if err2 != nil {
					logger.Log.Sugar().Errorf("Error rollback transaction: %v", err2)
				}
				return err
			}

		} else {

			err = pgdb.InsertMetric(ctx, met.ID, met.MType, met.Delta, met.Value)
			if err != nil {
				logger.Log.Sugar().Errorf("Error inserting data: %v, rollback transaction", err)
				err2 := tx.Rollback()
				if err2 != nil {
					logger.Log.Sugar().Errorf("Error rollback transaction: %v", err2)
				}
				return err
			}
		}

		err = tx.Commit()
		if err != nil {
			logger.Log.Sugar().Errorf("Error commit transaction: %v", err)
		}

	}
	logger.Log.Sugar().Infof("Commited all transactions, data saved to DB successfully")
	return nil
}

func contains(ids []string, id string) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}
