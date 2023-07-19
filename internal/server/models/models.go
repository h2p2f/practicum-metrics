// Package models описывает модель данных, используемые в проекте.
//
// Package models describes the data model used in the project.
package models

// Metric - модель данных для метрики.
//
// Metric - data model for metric.
type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
