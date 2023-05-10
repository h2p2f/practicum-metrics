package database

import "context"

type DataBaser interface {
	CreateTable(ctx context.Context) error
	ReadFromDB(ctx context.Context) ([][]byte, error)
	SaveToDB(ctx context.Context, metrics [][]byte) error
	SaveToDBWithoutTruncate(ctx context.Context, metrics [][]byte) error
}

type metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type DB struct {
	DataBase DataBaser
}

func NewDB(db DataBaser) *DB {
	return &DB{
		DataBase: db,
	}
}
