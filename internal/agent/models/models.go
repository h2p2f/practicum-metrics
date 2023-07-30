//Package models реализует модели данных.
//
// Package models implements data models.

package models

type Metric struct {
	Value *float64 `json:"value,omitempty"`
	Delta *int64   `json:"delta,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}
