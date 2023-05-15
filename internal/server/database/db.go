package database

import "context"

// DataBaser is an interface for database, file or mock
type DataBaser interface {
	Read(ctx context.Context) ([][]byte, error)
	Write(ctx context.Context, metrics [][]byte) error
	Create(ctx context.Context) error
}

// metrics is a struct for json
type metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// DB is a struct for database
type DB struct {
	DataBase DataBaser
}

// NewDB is a function that returns a new DB
func NewDB(db DataBaser) *DB {
	return &DB{DataBase: db}
}
