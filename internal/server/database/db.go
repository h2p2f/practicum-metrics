package database

import "context"

//import _ "github.com/golang/mock/mockgen/model"

type DataBaser interface {
	Read(ctx context.Context) ([][]byte, error)
	Write(ctx context.Context, metrics [][]byte) error
	Create(ctx context.Context) error
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
	return &DB{DataBase: db}
}
