// Package models describes the data model used in the project.
package models

// Metric - data model for metric.
type Metric struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}
